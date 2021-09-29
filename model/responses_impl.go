package model

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

// Response ResponseRepositoryの実装
type Response struct{}

// NewResponse Responseのコンストラクター
func NewResponse() *Response {
	return new(Response)
}

//Responses responseテーブルの構造体
type Responses struct {
	ResponseID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	QuestionID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	Body       null.String    `json:"response" gorm:"type:text;default:NULL"`
	ModifiedAt time.Time      `json:"-" gorm:"type:timestamp;not null;dafault:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"type:TIMESTAMP NULL;default:NULL"`
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (r *Responses) BeforeCreate(tx *gorm.DB) error {
	r.ModifiedAt = time.Now()

	return nil
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (r *Responses) BeforeUpdate(tx *gorm.DB) error {
	r.ModifiedAt = time.Now()

	return nil
}

//TableName テーブル名が単数形なのでその対応
func (*Responses) TableName() string {
	return "response"
}

// ResponseBody 質問に対する回答の構造体
type ResponseBody struct {
	QuestionID     int         `json:"questionID" gorm:"column:id" validate:"min=0"`
	QuestionType   string      `json:"question_type" gorm:"column:type" validate:"required,oneof=Text TextArea Number MultipleChoice Checkbox LinearScale"`
	Body           null.String `json:"response" validate:"required"`
	OptionResponse []string    `json:"option_response" validate:"required_if=QuestionType Checkbox,required_if=QuestionType MultipleChoice,dive,max=50"`
}

// ResponseMeta 質問に対する回答の構造体
type ResponseMeta struct {
	QuestionID int
	Data       string
}

// InsertResponses 質問に対する回答の追加
func (*Response) InsertResponses(ctx context.Context, responseID int, responseMetas []*ResponseMeta) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	responses := make([]Responses, 0, len(responseMetas))
	for _, responseMeta := range responseMetas {
		responses = append(responses, Responses{
			ResponseID: responseID,
			QuestionID: responseMeta.QuestionID,
			Body:       null.NewString(responseMeta.Data, true),
		})
	}
	err = db.Create(&responses).Error
	if err != nil {
		return fmt.Errorf("failed to insert responses: %w", err)
	}

	return nil
}

// DeleteResponse 質問に対する回答の削除
func (*Response) DeleteResponse(ctx context.Context, responseID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	result := db.
		Where("response_id = ?", responseID).
		Delete(&Responses{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete response: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete response: %w", ErrNoRecordDeleted)
	}

	return nil
}
