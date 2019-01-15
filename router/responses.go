package router

import (
	"net/http"
	"strconv"
	"time"

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
		"responseID":      responseID,
		"questionnaireID": req.ID,
		"submitted_at":    req.SubmittedAt,
		"body":            req.Body,
	})
}

func GetMyResponses(c echo.Context) error {
	responsesinfo := []model.ResponseInfo{}

	if err := model.DB.Select(&responsesinfo,
		`SELECT questionnaire_id, response_id, modified_at, submitted_at from respondents
		WHERE user_traqid = ? AND deleted_at IS NULL`,
		model.GetUserID(c)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	myresponses, err := model.GetResponsesInfo(c, responsesinfo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myresponses)
}

func GetMyResponsesByID(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	responsesinfo := []model.ResponseInfo{}

	if err := model.DB.Select(&responsesinfo,
		`SELECT questionnaire_id, response_id, modified_at, submitted_at from respondents
		WHERE user_traqid = ? AND deleted_at IS NULL AND questionnaire_id = ?`,
		model.GetUserID(c), questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	myresponses, err := model.GetResponsesInfo(c, responsesinfo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myresponses)
}

func GetResponsesByID(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resSharedTo := ""
	if err := model.DB.Get(&resSharedTo,
		`SELECT res_shared_to FROM questionnaires WHERE deleted_at IS NULL AND id = ?`,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	switch resSharedTo {
	case "administrators":
		AmAdmin, err := model.CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	case "respondents":
		AmAdmin, err := model.CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		RespondedAt, err := model.RespondedAt(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin && RespondedAt == "NULL" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	}

	responsesinfo := []struct {
		ResponseID  int            `db:"response_id"`
		UserID      string         `db:"user_traqid"`
		ModifiedAt  time.Time      `db:"modified_at"`
		SubmittedAt mysql.NullTime `db:"submitted_at"`
	}{}

	if err := model.DB.Select(&responsesinfo,
		`SELECT response_id, user_traqid, modified_at, submitted_at from respondents
		WHERE deleted_at IS NULL AND questionnaire_id = ? AND submitted_at IS NOT NULL`,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	type ResponsesInfo struct {
		ResponseID  int                  `json:"responseID"`
		UserID      string               `json:"traqID"`
		SubmittedAt string               `json:"submitted_at"`
		ModifiedAt  string               `json:"modified_at"`
		Body        []model.ResponseBody `json:"response_body"`
	}
	responses := []ResponsesInfo{}

	questionTypeList, err := model.GetQuestionsType(c, questionnaireID)
	if err != nil {
		return err
	}

	for _, response := range responsesinfo {
		bodyList := []model.ResponseBody{}
		for _, questionType := range questionTypeList {
			body, err := model.GetResponseBody(c, response.ResponseID, questionType.ID, questionType.Type)
			if err != nil {
				return err
			}
			bodyList = append(bodyList, body)
		}
		responses = append(responses,
			ResponsesInfo{
				ResponseID:  response.ResponseID,
				UserID:      response.UserID,
				SubmittedAt: model.NullTimeToString(response.SubmittedAt),
				ModifiedAt:  response.ModifiedAt.Format(time.RFC3339),
				Body:        bodyList,
			})
	}

	return c.JSON(http.StatusOK, responses)
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
		model.NullTimeToString(respondentInfo.SubmittedAt),
		model.NullTimeToString(respondentInfo.ModifiedAt),
		[]model.ResponseBody{},
	}

	questionTypeList, err := model.GetQuestionsType(c, responses.QuestionnaireID)
	if err != nil {
		return err
	}

	bodyList := []model.ResponseBody{}
	for _, questionType := range questionTypeList {
		body, err := model.GetResponseBody(c, responseID, questionType.ID, questionType.Type)
		if err != nil {
			return err
		}
		bodyList = append(bodyList, body)
	}
	responses.Body = bodyList
	return c.JSON(http.StatusOK, responses)
}

func EditResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	req := model.Responses{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if req.SubmittedAt == "" || req.SubmittedAt == "NULL" {
		req.SubmittedAt = "NULL"
		if _, err := model.DB.Exec(
			`UPDATE respondents
			SET questionnaire_id = ?, submitted_at = NULL, modified_at = CURRENT_TIMESTAMP
			WHERE response_id = ?`,
			req.ID, responseID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := model.DB.Exec(
			`UPDATE respondents
			SET questionnaire_id = ?, submitted_at = ?, modified_at = CURRENT_TIMESTAMP
			WHERE response_id = ?`,
			req.ID, req.SubmittedAt, responseID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	//全消し&追加(レコード数爆発しそう)
	if _, err := model.DB.Exec(
		`UPDATE responses SET deleted_at = CURRENT_TIMESTAMP WHERE response_id = ?`,
		responseID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
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
		`UPDATE response SET deleted_at = CURRENT_TIMESTAMP WHERE response_id = ?`,
		responseID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
