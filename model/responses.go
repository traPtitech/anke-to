package model

import (
	"net/http"
	"strconv"
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

type UserResponse struct {
	ResponseID  int            `db:"response_id"`
	UserID      string         `db:"user_traqid"`
	ModifiedAt  time.Time      `db:"modified_at"`
	SubmittedAt mysql.NullTime `db:"submitted_at"`
}

type ResponseID struct {
	QuestionnaireID int            `db:"questionnaire_id"`
	ModifiedAt      mysql.NullTime `db:"modified_at"`
	SubmittedAt     mysql.NullTime `db:"submitted_at"`
}

func InsertRespondents(c echo.Context, req Responses) (int, error) {
	var result sql.Result
	var err error
	if req.SubmittedAt == "" || req.SubmittedAt == "NULL" {
		req.SubmittedAt = "NULL"
		if result, err = db.Exec(
			`INSERT INTO respondents (questionnaire_id, user_traqid, modified_at) VALUES (?, ?, ?)`,
			req.ID, GetUserID(c), time.Now()); err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if result, err = db.Exec(
			`INSERT INTO respondents
				(questionnaire_id, user_traqid, submitted_at, modified_at) VALUES (?, ?, ?, ?)`,
			req.ID, GetUserID(c), req.SubmittedAt, time.Now()); err != nil {
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
	if _, err := db.Exec(
		`INSERT INTO response (response_id, question_id, body, modified_at) VALUES (?, ?, ?, ?)`,
		responseID, body.QuestionID, data, time.Now()); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func GetRespondents(c echo.Context, questionnaireID int) ([]string, error) {
	respondents := []string{}
	if err := db.Select(&respondents,
		"SELECT user_traqid FROM respondents WHERE questionnaire_id = ? AND deleted_at IS NULL AND submitted_at IS NOT NULL",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return respondents, nil
}

func GetMyResponses(c echo.Context) ([]ResponseInfo, error) {
	responsesinfo := []ResponseInfo{}

	if err := db.Select(&responsesinfo,
		`SELECT questionnaire_id, response_id, modified_at, submitted_at from respondents
		WHERE user_traqid = ? AND deleted_at IS NULL ORDER BY modified_at DESC`,
		GetUserID(c)); err != nil {
		c.Logger().Error(err)
		return []ResponseInfo{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return responsesinfo, nil
}

func GetMyResponsesByID(c echo.Context, questionnaireID int) ([]ResponseInfo, error) {
	responsesinfo := []ResponseInfo{}

	if err := db.Select(&responsesinfo,
		`SELECT questionnaire_id, response_id, modified_at, submitted_at from respondents
		WHERE user_traqid = ? AND deleted_at IS NULL AND questionnaire_id = ? ORDER BY modified_at DESC`,
		GetUserID(c), questionnaireID); err != nil {
		c.Logger().Error(err)
		return []ResponseInfo{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return responsesinfo, nil
}

func GetResponsesInfo(c echo.Context, responsesinfo []ResponseInfo) ([]MyResponse, error) {
	myresponses := []MyResponse{}

	for _, response := range responsesinfo {
		title, resTimeLimit, err := GetTitleAndLimit(c, response.QuestionnaireID)
		if title == "" {
			continue
		}
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
		if err := db.Select(&option,
			`SELECT body from response
			WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
			responseID, body.QuestionID); err != nil {
			c.Logger().Error(err)
			return ResponseBody{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
		body.OptionResponse = option
		// sortで比較するため
		for _, op := range option {
			if body.Response != "" {
				body.Response += ", "
			}
			body.Response += op
		}
	default:
		var response string
		if err := db.Get(&response,
			`SELECT body from response
			WHERE response_id = ? AND question_id = ? AND deleted_at IS NULL`,
			responseID, body.QuestionID); err != nil {
			if err != sql.ErrNoRows {
				c.Logger().Error(err)
				return ResponseBody{}, echo.NewHTTPError(http.StatusInternalServerError)
			}
		}
		body.Response = response
	}
	return body, nil
}

func RespondedAt(c echo.Context, questionnaireID int) (string, error) {
	respondedAt := sql.NullString{}
	if err := db.Get(&respondedAt,
		`SELECT MAX(submitted_at) FROM respondents
		WHERE user_traqid = ? AND questionnaire_id = ? AND deleted_at IS NULL`,
		GetUserID(c), questionnaireID); err != nil {
		c.Logger().Error(err)
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}
	return NullStringConvert(respondedAt), nil
}

func GetRespondentByID(c echo.Context, responseID int) (ResponseID, error) {
	respondentInfo := ResponseID{}
	if err := db.Get(&respondentInfo,
		`SELECT questionnaire_id, modified_at, submitted_at from respondents
		WHERE response_id = ? AND user_traqid = ? AND deleted_at IS NULL`,
		responseID, GetUserID(c)); err != nil {
		if err != sql.ErrNoRows {
			c.Logger().Error(err)
			return ResponseID{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return respondentInfo, nil
}

func UpdateRespondents(c echo.Context, questionnaireID int, responseID int, submittedAt string) error {
	if submittedAt == "" || submittedAt == "NULL" {
		submittedAt = "NULL"
		if _, err := db.Exec(
			`UPDATE respondents
			SET questionnaire_id = ?, submitted_at = NULL, modified_at = ? WHERE response_id = ?`,
			questionnaireID, time.Now(), responseID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := db.Exec(
			`UPDATE respondents
			SET questionnaire_id = ?, submitted_at = ?, modified_at = ? WHERE response_id = ?`,
			questionnaireID, submittedAt, time.Now(), responseID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil

}

func DeleteResponse(c echo.Context, responseID int) error {
	if _, err := db.Exec(
		`UPDATE response SET deleted_at = ? WHERE response_id = ?`,
		time.Now(), responseID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteMyResponse(c echo.Context, responseID int) error {
	if _, err := db.Exec(
		`UPDATE respondents SET deleted_at = ? WHERE response_id = ? AND user_traqid = ?`,
		time.Now(), responseID, GetUserID(c)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if _, err := db.Exec(
		`UPDATE response SET deleted_at = ? WHERE response_id = ?`,
		time.Now(), responseID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// アンケートの回答を確認できるか
func CheckResponseConfirmable(c echo.Context, resSharedTo string, questionnaireID int) error {
	switch resSharedTo {
	case "administrators":
		AmAdmin, err := CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	case "respondents":
		AmAdmin, err := CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		RespondedAt, err := RespondedAt(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin && RespondedAt == "NULL" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	}
	return nil
}

type ResponseAndBody struct {
	ResponseID  int            `db:"response_id"`
	UserID      string         `db:"user_traqid"`
	ModifiedAt  time.Time      `db:"modified_at"`
	SubmittedAt mysql.NullTime `db:"submitted_at"`
	QuestionID  int            `db:"question_id"`
	Body        string         `db:"body"`
}

func GetResponsesByID(questionnaireID int) ([]ResponseAndBody, error) {
	responses := []ResponseAndBody{}
	if err := db.Select(&responses,
		`SELECT respondents.response_id AS response_id,
		user_traqid, 
		respondents.modified_at AS modified_at,
		respondents.submitted_at AS submitted_at,
		response.question_id,
		response.body
		FROM respondents
		RIGHT OUTER JOIN response
		ON respondents.response_id = response.response_id
		WHERE respondents.questionnaire_id = ?
		AND respondents.deleted_at IS NULL
		AND response.deleted_at IS NULL
		AND respondents.submitted_at IS NOT NULL`, questionnaireID); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return responses, nil
}

// sortされた回答者の情報を返す
func GetSortedRespondents(c echo.Context, questionnaireID int, sortQuery string) ([]UserResponse, int, error) {
	sql := `SELECT response_id, user_traqid, modified_at, submitted_at from respondents
			WHERE deleted_at IS NULL AND questionnaire_id = ? AND submitted_at IS NOT NULL`

	sortNum := 0
	switch sortQuery {
	case "traqid":
		sql += ` ORDER BY user_traqid`
	case "-traqid":
		sql += ` ORDER BY user_traqid DESC`
	case "submitted_at":
		sql += ` ORDER BY submitted_at`
	case "-submitted_at":
		sql += ` ORDER BY submitted_at DESC`
	case "":
	default:
		var err error
		sortNum, err = strconv.Atoi(sortQuery)
		if err != nil {
			c.Logger().Error(err)
			return []UserResponse{}, 0, echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	responsesinfo := []UserResponse{}
	if err := db.Select(&responsesinfo, sql,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return []UserResponse{}, 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return responsesinfo, sortNum, nil
}

type QIDandResponse struct {
	QuestionID int
	Response   string
}

func GetResponseBodyList(c echo.Context, questionTypeList []QuestionIDType, responses []QIDandResponse) []ResponseBody {
	bodyList := []ResponseBody{}

	for _, qType := range questionTypeList {
		response := ""
		optionResponse := []string{}
		for _, respInfo := range responses {
			// 質問IDが一致したら追加
			if qType.ID == respInfo.QuestionID {
				switch qType.Type {
				case "MultipleChoice", "Checkbox", "Dropdown":
					if response != "" {
						response += ","
					}
					response += respInfo.Response
					optionResponse = append(optionResponse, respInfo.Response)
				default:
					response += respInfo.Response
				}
			}
		}
		// 回答内容の配列に追加
		bodyList = append(bodyList,
			ResponseBody{
				QuestionID:     qType.ID,
				QuestionType:   qType.Type,
				Response:       response,
				OptionResponse: optionResponse,
			})
	}
	return bodyList
}
