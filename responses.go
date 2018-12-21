package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func InsertResponse(c echo.Context, responseID int, req responses, body responseBody, data string) error {
	if req.SubmittedAt.Valid {
		if _, err := db.Exec(
			`INSERT INTO responses
				(questionnaire_id, question_id, response_id, body, user_traqid, submitted_at)
				VALUES (?, ?, ?, ?, ?, ?)`,
			req.ID, body.QuestionID, responseID, data, getUserID(c), req.SubmittedAt); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := db.Exec(
			`INSERT INTO responses
				(questionnaire_id, question_id, response_id, body, user_traqid)
				VALUES (?, ?, ?, ?, ?)`,
			req.ID, body.QuestionID, responseID, data, getUserID(c)); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

func postResponse(c echo.Context) error {

	req := responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// responseIDを全部取って最大値+1をIDとする(もっといい方法がありそう)
	responseIDlist := []int{}

	if err := db.Select(&responseIDlist, "SELECT response_id from responses"); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	maxResponseID := 0
	for _, responseID := range responseIDlist {
		if maxResponseID < responseID {
			maxResponseID = responseID
		}
	}

	responseID := maxResponseID + 1

	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				if err := InsertResponse(c, responseID, req, body, option); err != nil {
					return err
				}
			}
		default:
			if err := InsertResponse(c, responseID, req, body, body.Response); err != nil {
				return err
			}
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": req.ID,
		"submitted_at":    timeConvert(req.SubmittedAt),
		"body":            req.Body,
	})
}

func getMyResponses(c echo.Context) error {
	responsesinfo := []struct {
		QuestionnaireID int            `db:"questionnaire_id"`
		ResponseID      int            `db:"response_id"`
		ModifiedAt      mysql.NullTime `db:"modified_at"`
		SubmittedAt     mysql.NullTime `db:"submitted_at"`
	}{}

	if err := db.Select(&responsesinfo,
		`SELECT questionnaire_id, response_id, modified_at, submitted_at from responses
		WHERE user_traqid = ? AND deleted_at IS NULL`,
		getUserID(c)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	type MyResponse struct {
		ResponseID      int    `json:"responseID"`
		QuestionnaireID int    `json:"questionnaireID"`
		Title           string `json:"questionnaire_title"`
		ResTimeLimit    string `json:"res_time_limit"`
		SubmittedAt     string `json:"submitted_at"`
		ModifiedAt      string `json:"modified_at"`
	}
	myresponses := []MyResponse{}

	for _, response := range responsesinfo {
		duplication := false
		for _, other := range myresponses {
			if response.ResponseID == other.ResponseID {
				duplication = true
				break
			}
		}
		if !duplication {
			fmt.Println(response.QuestionnaireID)
			title, resTimeLimit, err := GetTitleAndLimit(c, response.QuestionnaireID)
			if err != nil {
				return err
			}
			myresponses = append(myresponses,
				MyResponse{
					ResponseID:      response.ResponseID,
					QuestionnaireID: response.QuestionnaireID,
					Title:           title,
					ResTimeLimit:    resTimeLimit,
					SubmittedAt:     timeConvert(response.SubmittedAt),
					ModifiedAt:      timeConvert(response.ModifiedAt),
				})
		}
	}

	return c.JSON(http.StatusOK, myresponses)
}

func getResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	responsesinfo := []struct {
		QuestionnaireID int            `db:"questionnaire_id"`
		QuestionID      int            `db:"question_id"`
		Body            string         `db:"body"`
		ModifiedAt      mysql.NullTime `db:"modified_at"`
		SubmittedAt     mysql.NullTime `db:"submitted_at"`
	}{}

	if err := db.Select(&responsesinfo,
		`SELECT questionnaire_id, question_id, body, modified_at, submitted_at from responses
		WHERE response_id = ? AND user_traqid = ? AND deleted_at IS NULL`,
		responseID, getUserID(c)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if len(responsesinfo) == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	responses := struct {
		QuestionnaireID int            `json:"questionnaireID"`
		SubmittedAt     string         `json:"submitted_at"`
		ModifiedAt      string         `json:"modified_at"`
		Body            []responseBody `json:"body"`
	}{responsesinfo[0].QuestionnaireID,
		timeConvert(responsesinfo[0].SubmittedAt),
		timeConvert(responsesinfo[0].ModifiedAt),
		[]responseBody{},
	}

	body := []responseBody{}
	for _, response := range responsesinfo {
		questionType := ""
		questionID := response.QuestionID
		if err := db.Get(&questionType,
			"SELECT type FROM questions WHERE id = ?", questionID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		switch questionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			duplication := false
			for i, b := range body {
				if b.QuestionID == questionID {
					body[i].OptionResponse = append(body[i].OptionResponse, response.Body)
					duplication = true
					break
				}
			}
			if !duplication {
				optionResponse := []string{response.Body}
				body = append(body,
					responseBody{
						QuestionID:     questionID,
						QuestionType:   questionType,
						Response:       "",
						OptionResponse: optionResponse,
					})
			}
		default:
			body = append(body,
				responseBody{
					QuestionID:     questionID,
					QuestionType:   questionType,
					Response:       response.Body,
					OptionResponse: []string{},
				})
		}
	}
	responses.Body = body

	return c.JSON(http.StatusOK, responses)
}

func editResponse(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func deleteResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if _, err := db.Exec(
		`UPDATE responses SET deleted_at = CURRENT_TIMESTAMP
		WHERE response_id = ? AND user_traqid = ?`,
		responseID, getUserID(c)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
