package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

//Response responseテーブルの構造体
type Response struct {
	ResponseID int         `json:"-" gorm:"type:int(11) NOT NULL;"`
	QuestionID int         `json:"-" gorm:"type:int(11) NOT NULL;"`
	Body       null.String `json:"response" gorm:"type:text;default:NULL;"`
	ModifiedAt time.Time   `json:"-" gorm:"type:timestamp NOT NULL;DEFAULT:CURRENT_TIMESTAMP;"`
	DeletedAt  null.Time   `json:"-" gorm:"type:timestamp NULL;default:NULL;"`
}

//TableName テーブル名が単数形なのでその対応
func (*Response) TableName() string {
	return "response"
}

// ResponseBody 質問に対する回答の構造体
type ResponseBody struct {
	QuestionID     int         `json:"questionID" gorm:"column:id"`
	QuestionType   string      `json:"question_type" gorm:"column:type"`
	Body           null.String `json:"response,omitempty"`
	OptionResponse []string    `json:"option_response"`
}

// Responses 質問に対する回答一覧の構造体
type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt null.Time      `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}

// InsertResponse 質問に対する回答の追加
func InsertResponse(c echo.Context, responseID int, questionID int, data string) error {
	err := db.Create(&Response{
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

// DeleteResponse 質問に対する回答の削除
func DeleteResponse(c echo.Context, responseID int) error {
	err := db.
		Where("response_id = ?", responseID).
		Delete(&Response{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}
