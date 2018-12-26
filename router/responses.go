package router

import (
	"net/http"
	"strconv"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

func PostResponse(c echo.Context) error {

	req := model.Responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	responseID, err := model.InsertRespondents(c, req)
	if err != nil {
		return nil
	}

	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				if err := model.InsertResponse(c, responseID, req, body, option); err != nil {
					return err
				}
			}
		default:
			if err := model.InsertResponse(c, responseID, req, body, body.Response); err != nil {
				return err
			}
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": req.ID,
		"submitted_at":    model.TimeConvert(req.SubmittedAt),
		"body":            req.Body,
	})
}

func GetMyResponses(c echo.Context) error {
	responsesinfo := []struct {
		QuestionnaireID int            `db:"questionnaire_id"`
		ResponseID      int            `db:"response_id"`
		ModifiedAt      mysql.NullTime `db:"modified_at"`
		SubmittedAt     mysql.NullTime `db:"submitted_at"`
	}{}

	if err := model.DB.Select(&responsesinfo,
		`SELECT questionnaire_id, response_id, modified_at, submitted_at from respondents
		WHERE user_traqid = ? AND deleted_at IS NULL`,
		model.GetUserID(c)); err != nil {
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
		title, resTimeLimit, err := model.GetTitleAndLimit(c, response.QuestionnaireID)
		if err != nil {
			return err
		}
		myresponses = append(myresponses,
			MyResponse{
				ResponseID:      response.ResponseID,
				QuestionnaireID: response.QuestionnaireID,
				Title:           title,
				ResTimeLimit:    resTimeLimit,
				SubmittedAt:     model.TimeConvert(response.SubmittedAt),
				ModifiedAt:      model.TimeConvert(response.ModifiedAt),
			})
	}

	return c.JSON(http.StatusOK, myresponses)
}

func GetResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	respondentInfo := struct {
		QuestionnaireID int            `db:"questionnaire_id"`
		ModifiedAt      mysql.NullTime `db:"modified_at"`
		SubmittedAt     mysql.NullTime `db:"submitted_at"`
	}{}

	if err := model.DB.Get(&respondentInfo,
		`SELECT questionnaire_id, modified_at, submitted_at from respondents
		WHERE response_id = ? AND user_traqid = ? AND deleted_at IS NULL`,
		responseID, model.GetUserID(c)); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	responses := struct {
		QuestionnaireID int                  `json:"questionnaireID"`
		SubmittedAt     string               `json:"submitted_at"`
		ModifiedAt      string               `json:"modified_at"`
		Body            []model.ResponseBody `json:"body"`
	}{
		respondentInfo.QuestionnaireID,
		model.TimeConvert(respondentInfo.SubmittedAt),
		model.TimeConvert(respondentInfo.ModifiedAt),
		[]model.ResponseBody{},
	}

	questionTypeList, err := model.GetQuestionsType(c, responses.QuestionnaireID)
	if err != nil {
		return err
	}

	bodyList := []model.ResponseBody{}
	for _, questionType := range questionTypeList {
		body := model.ResponseBody{
			QuestionID:   questionType.ID,
			QuestionType: questionType.Type,
		}
		switch questionType.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			option := []string{}
			if err := model.DB.Select(&option,
				`SELECT body from responses
				WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
				responseID, body.QuestionID); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			body.OptionResponse = option
		default:
			var response string
			if err := model.DB.Get(&response,
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

func EditResponse(c echo.Context) error {
	// 後で実装
	return c.NoContent(http.StatusOK)
}

func DeleteResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if _, err := model.DB.Exec(
		`UPDATE respondents SET deleted_at = CURRENT_TIMESTAMP
		WHERE response_id = ? AND user_traqid = ?`,
		responseID, model.GetUserID(c)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if _, err := model.DB.Exec(
		`UPDATE responses SET deleted_at = CURRENT_TIMESTAMP WHERE response_id = ?`,
		responseID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
