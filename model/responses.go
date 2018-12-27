package model

import (
	"net/http"

	"database/sql"
	"github.com/labstack/echo"
)

type ResponseBody struct {
	QuestionID     int      `json:"questionID"`
	QuestionType   string   `json:"question_type"`
	Response       string   `json:"response"`
	OptionResponse []string `json:"option_response"`
}

type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt string         `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}

func InsertRespondents(c echo.Context, req Responses) (int, error) {
	var result sql.Result
	var err error
	if req.SubmittedAt != "NULL" {
		if result, err = DB.Exec(
			`INSERT INTO respondents
				(questionnaire_id, user_traqid, submitted_at) VALUES (?, ?, ?)`,
			req.ID, GetUserID(c), req.SubmittedAt); err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if result, err = DB.Exec(
			`INSERT INTO respondents (questionnaire_id, user_traqid) VALUES (?, ?)`,
			req.ID, GetUserID(c)); err != nil {
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

func InsertResponse(c echo.Context, responseID int, req Responses, body ResponseBody, data string) error {
	if _, err := DB.Exec(
		`INSERT INTO responses (response_id, question_id, body) VALUES (?, ?, ?)`,
		responseID, body.QuestionID, data); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func GetRespondents(c echo.Context, questionnaireID int) ([]string, error) {
	respondents := []string{}
	if err := DB.Select(&respondents,
		"SELECT user_traqid FROM respondents WHERE questionnaire_id = ? AND deleted_at IS NULL",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return respondents, nil
}
