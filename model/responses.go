package model

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

//Response responseテーブルの構造体
type Response struct {
	ResponseID int         `json:"-" gorm:"type:int(11);NOT NULL;"`
	QuestionID int         `json:"-" gorm:"type:int(11);NOT NULL;"`
	Body       null.String `json:"response" gorm:"type:text;"`
	ModifiedAt time.Time   `json:"-" gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;"`
	DeletedAt  null.Time   `json:"-" gorm:"type:timestamp;"`
}

//TableName テーブル名が単数形なのでその対応
func (*Response) TableName() string {
	return "response"
}

type ResponseBody struct {
	Question `gorm:"embedded"`
	Body null.String
	OptionResponse []string `json:"option_response"`
}

type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt null.Time      `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}

type UserResponse struct {
	ResponseID  int       `db:"response_id"`
	UserID      string    `db:"user_traqid"`
	ModifiedAt  time.Time `db:"modified_at"`
	SubmittedAt null.Time `db:"submitted_at"`
}

func InsertResponse(c echo.Context, responseID int, questionID int, data string) error {
	err := gormDB.Create(Response{
		ResponseID: responseID,
		QuestionID: questionID,
		Body:       null.NewString(data, true),
	}).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to insert response: %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func DeleteResponse(c echo.Context, responseID int) error {
	err := gormDB.
		Where("response_id = ?", responseID).
		Delete(&Response{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
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
				Question: Question{
					ID: qType.ID,
					Type: qType.Type,
				},
				Body:       null.NewString(response, true),
				OptionResponse: optionResponse,
			})
	}
	return bodyList
}
