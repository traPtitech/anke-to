package model

import (
	"fmt"
	"time"

	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
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
	Body           null.String `json:"response"`
	OptionResponse []string    `json:"option_response"`
}

// Responses 質問に対する回答一覧の構造体
type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt null.Time      `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}

// ResponseMeta 質問に対する回答の構造体
type ResponseMeta struct {
	QuestionID int
	Data       string
}

// InsertResponses 質問に対する回答の追加
func InsertResponses(responseID int, responseMetas []*ResponseMeta) error {
	responses := make([]interface{}, 0, len(responseMetas))
	for _, responseMeta := range responseMetas {
		responses = append(responses, Response{
			ResponseID: responseID,
			QuestionID: responseMeta.QuestionID,
			Body:       null.NewString(responseMeta.Data, true),
		})
	}
	err := gormbulk.BulkInsert(db, responses, len(responses), "ModifiedAt", "DeletedAt")
	if err != nil {
		return fmt.Errorf("failed to insert response: %w", err)
	}

	return nil
}

// DeleteResponse 質問に対する回答の削除
func DeleteResponse(responseID int) error {
	err := db.
		Where("response_id = ?", responseID).
		Delete(&Response{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete response: %w", err)
	}

	return nil
}
