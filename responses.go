package main

import (
	"fmt"
	"net/http"
	"strconv"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func InsertRespondents(c echo.Context, req responses) (int, error) {
	var result sql.Result
	var err error
	if req.SubmittedAt.Valid {
		if result, err = db.Exec(
			`INSERT INTO respondents
				(questionnaire_id, user_traqid, submitted_at) VALUES (?, ?, ?)`,
			req.ID, getUserID(c), req.SubmittedAt); err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if result, err = db.Exec(
			`INSERT INTO respondents (questionnaire_id, user_traqid) VALUES (?, ?)`,
			req.ID, getUserID(c)); err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		c.Logger().Error(err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return int(lastID), nil
}

func InsertResponse(c echo.Context, responseID int, req responses, body responseBody, data string) error {
	if _, err := db.Exec(
		`INSERT INTO responses (response_id, question_id, body) VALUES (?, ?, ?)`,
		responseID, body.QuestionID, data); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func postResponse(c echo.Context) error {

	req := responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	responseID, err := InsertRespondents(c, req)
	if err != nil {
		return nil
	}

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

	respondentInfo := struct {
		QuestionnaireID int            `db:"questionnaire_id"`
		ModifiedAt      mysql.NullTime `db:"modified_at"`
		SubmittedAt     mysql.NullTime `db:"submitted_at"`
	}{}

	if err := db.Get(&respondentInfo,
		`SELECT questionnaire_id, modified_at, submitted_at from respondents
		WHERE response_id = ? AND user_traqid = ? AND deleted_at IS NULL`,
		responseID, getUserID(c)); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	responses := struct {
		QuestionnaireID int            `json:"questionnaireID"`
		SubmittedAt     string         `json:"submitted_at"`
		ModifiedAt      string         `json:"modified_at"`
		Body            []responseBody `json:"body"`
	}{
		respondentInfo.QuestionnaireID,
		timeConvert(respondentInfo.SubmittedAt),
		timeConvert(respondentInfo.ModifiedAt),
		[]responseBody{},
	}

	questionTypeList, err := getQuestionsType(c, responses.QuestionnaireID)
	if err != nil {
		return err
	}

	bodyList := []responseBody{}
	for _, questionType := range questionTypeList {
		body := responseBody{
			QuestionID:   questionType.ID,
			QuestionType: questionType.Type,
		}
		switch questionType.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			option := []string{}
			if err := db.Select(&option,
				`SELECT body from responses
				WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
				responseID, body.QuestionID); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			body.OptionResponse = option
		default:
			var response string
			if err := db.Get(&response,
				`SELECT body from responses
				WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
				responseID, body.QuestionID); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			body.Response = response
		}
		bodyList = append(bodyList, body)
	}
	responses.Body = bodyList
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
