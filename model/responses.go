package model

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

//Response responseテーブルの構造体
type Response struct {
	ResponseID int         `gorm:"type:int(11);NOT NULL;"`
	QuestionID int         `gorm:"type:int(11);NOT NULL;"`
	Body       null.String `gorm:"type:text;"`
	ModifiedAt time.Time   `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;"`
	DeletedAt  null.Time   `gorm:"type:timestamp;"`
}

//TableName テーブル名が単数形なのでその対応
func (*Response) TableName() string {
	return "response"
}

type ResponseBody struct {
	QuestionID     int      `json:"questionID"`
	QuestionType   string   `json:"question_type"`
	Response       string   `json:"response"`
	OptionResponse []string `json:"option_response"`
}

type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt null.Time      `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}

type ResponseInfo struct {
	QuestionnaireID int       `db:"questionnaire_id"`
	ResponseID      int       `db:"response_id"`
	ModifiedAt      time.Time `db:"modified_at"`
	SubmittedAt     null.Time `db:"submitted_at"`
}

type UserResponse struct {
	ResponseID  int       `db:"response_id"`
	UserID      string    `db:"user_traqid"`
	ModifiedAt  time.Time `db:"modified_at"`
	SubmittedAt null.Time `db:"submitted_at"`
}

type ResponseID struct {
	QuestionnaireID int       `db:"questionnaire_id"`
	ModifiedAt      null.Time `db:"modified_at"`
	SubmittedAt     null.Time `db:"submitted_at"`
}

func InsertResponse(c echo.Context, responseID int, req Responses, body ResponseBody, data string) error {
	err := gormDB.Create(Response{
		ResponseID: responseID,
		QuestionID: body.QuestionID,
		Body:       null.NewString(data, true),
	}).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to insert response: %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
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

func RespondedAtBytraQID(c echo.Context, questionnaireID int, traQID string) (string, error) {
	respondedAt := sql.NullString{}
	if err := db.Get(&respondedAt,
		`SELECT MAX(submitted_at) FROM respondents
		WHERE user_traqid = ? AND questionnaire_id = ? AND deleted_at IS NULL`,
		traQID, questionnaireID); err != nil {
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

func UpdateRespondents(c echo.Context, questionnaireID int, responseID int, submittedAt null.Time) error {
	if !submittedAt.Valid {
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
	requestUser := GetUserID(c)

	res, err := db.Exec(
		`UPDATE respondents resp SET deleted_at = ? WHERE response_id = ? AND ( user_traqid = ? OR 
		EXISTS( SELECT * FROM administrators admin WHERE admin.questionnaire_id = resp.questionnaire_id AND admin.user_traqid = ? ))`,
		time.Now(), responseID, requestUser, requestUser)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return echo.NewHTTPError(http.StatusForbidden)
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
