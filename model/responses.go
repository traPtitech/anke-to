package model

import (
	"net/http"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"
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

type ResponseInfo struct {
	QuestionnaireID int            `db:"questionnaire_id"`
	ResponseID      int            `db:"response_id"`
	ModifiedAt      time.Time      `db:"modified_at"`
	SubmittedAt     mysql.NullTime `db:"submitted_at"`
}

type MyResponse struct {
	ResponseID      int    `json:"responseID"`
	QuestionnaireID int    `json:"questionnaireID"`
	Title           string `json:"questionnaire_title"`
	ResTimeLimit    string `json:"res_time_limit"`
	SubmittedAt     string `json:"submitted_at"`
	ModifiedAt      string `json:"modified_at"`
}

func InsertRespondents(c echo.Context, req Responses) (int, error) {
	var result sql.Result
	var err error
	if req.SubmittedAt == "" || req.SubmittedAt == "NULL" {
		req.SubmittedAt = "NULL"
		if result, err = DB.Exec(
			`INSERT INTO respondents (questionnaire_id, user_traqid) VALUES (?, ?)`,
			req.ID, GetUserID(c)); err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if result, err = DB.Exec(
			`INSERT INTO respondents
				(questionnaire_id, user_traqid, submitted_at) VALUES (?, ?, ?)`,
			req.ID, GetUserID(c), req.SubmittedAt); err != nil {
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
		`INSERT INTO response (response_id, question_id, body) VALUES (?, ?, ?)`,
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

func GetResponsesInfo(c echo.Context, responsesinfo []ResponseInfo) ([]MyResponse, error) {
	myresponses := []MyResponse{}

	for _, response := range responsesinfo {
		title, resTimeLimit, err := GetTitleAndLimit(c, response.QuestionnaireID)
		if err != nil {
			return nil, err
		}
		myresponses = append(myresponses,
			MyResponse{
				ResponseID:      response.ResponseID,
				QuestionnaireID: response.QuestionnaireID,
				Title:           title,
				ResTimeLimit:    resTimeLimit,
				SubmittedAt:     NullTimeToString(response.SubmittedAt),
				ModifiedAt:      response.ModifiedAt.Format(time.RFC3339),
			})
	}
	return myresponses, nil
}

func GetResponseBody(c echo.Context, responseID int, questionID int, questionType string) (ResponseBody, error) {
	body := ResponseBody{
		QuestionID:   questionID,
		QuestionType: questionType,
	}
	switch questionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		option := []string{}
		if err := DB.Select(&option,
			`SELECT body from response
			WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
			responseID, body.QuestionID); err != nil {
			c.Logger().Error(err)
			return ResponseBody{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		body.OptionResponse = option
	default:
		var response string
		if err := DB.Get(&response,
			`SELECT body from response
			WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
			responseID, body.QuestionID); err != nil {
			c.Logger().Error(err)
			return ResponseBody{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		body.Response = response
	}
	return body, nil
}

func RespondedAt(c echo.Context, questionnaireID int) (string, error) {
	respondedAt := sql.NullString{}
	if err := DB.Get(&respondedAt,
		`SELECT MAX(submitted_at) FROM respondents
		WHERE user_traqid = ? AND questionnaire_id = ? AND deleted_at IS NULL`,
		GetUserID(c), questionnaireID); err != nil {
		c.Logger().Error(err)
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}
	return NullStringConvert(respondedAt), nil
}
