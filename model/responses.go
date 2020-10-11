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
	QuestionID     int         `json:"questionID" gorm:"column:id"`
	QuestionType   string      `json:"question_type" gorm:"column:type"`
	Body           null.String `json:"response,omitempty"`
	OptionResponse []string    `json:"option_response"`
}

type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt null.Time      `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}

func InsertResponse(c echo.Context, responseID int, questionID int, data string) error {
	err := gormDB.Create(&Response{
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
