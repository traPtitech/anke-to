package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Question QuestionRepositoryの実装
type Question struct{}

// NewQuestion Questionのコンストラクター
func NewQuestion() *Question {
	return new(Question)
}

//Questions questionテーブルの構造体
type Questions struct {
	ID              int            `json:"id"                  gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID"     gorm:"type:int(11);not null"`
	PageNum         int            `json:"page_num"            gorm:"type:int(11);not null"`
	QuestionNum     int            `json:"question_num"        gorm:"type:int(11);not null"`
	Type            string         `json:"type"                gorm:"type:char(20);size:20;not null"`
	Body            string         `json:"body"                gorm:"type:text;default:NULL"`
	IsRequired      bool           `json:"is_required"         gorm:"type:tinyint(4);size:4;not null;default:0"`
	DeletedAt       gorm.DeletedAt `json:"-"          gorm:"type:TIMESTAMP NULL;default:NULL"`
	CreatedAt       time.Time      `json:"created_at"          gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Options         []Options      `json:"-"  gorm:"foreignKey:QuestionID"`
	Responses       []Responses    `json:"-"  gorm:"foreignKey:QuestionID"`
	ScaleLabels     []ScaleLabels  `json:"-"  gorm:"foreignKey:QuestionID"`
	Validations     []Validations  `json:"-"  gorm:"foreignKey:QuestionID"`
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (questionnaire *Questions) BeforeCreate(tx *gorm.DB) error {
	questionnaire.CreatedAt = time.Now()

	return nil
}

//TableName テーブル名が単数形なのでその対応
func (*Questions) TableName() string {
	return "question"
}

//QuestionIDType 質問のIDと種類の構造体
type QuestionIDType struct {
	ID   int
	Type string
}

//InsertQuestion 質問の追加
func (*Question) InsertQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool) (int, error) {
	db, err := getTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction: %w", err)
	}

	question := Questions{
		QuestionnaireID: questionnaireID,
		PageNum:         pageNum,
		QuestionNum:     questionNum,
		Type:            questionType,
		Body:            body,
		IsRequired:      isRequired,
	}

	err = db.
		Create(&question).Error
	if err != nil {
		return 0, fmt.Errorf("failed to insert a question record: %w", err)
	}

	return question.ID, nil
}

//UpdateQuestion 質問の修正
func (*Question) UpdateQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool, questionID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	question := map[string]interface{}{
		"questionnaire_id": questionnaireID,
		"page_num":         pageNum,
		"question_num":     questionNum,
		"type":             questionType,
		"body":             body,
		"is_required":      isRequired,
	}

	err = db.
		Model(&Questions{}).
		Where("id = ?", questionID).
		Updates(question).Error
	if err != nil {
		return fmt.Errorf("failed to update a question record: %w", err)
	}

	return nil
}

//DeleteQuestion 質問の削除
func (*Question) DeleteQuestion(ctx context.Context, questionID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	result := db.
		Where("id = ?", questionID).
		Delete(&Questions{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete a question record: %w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

//GetQuestions 質問一覧の取得
func (*Question) GetQuestions(ctx context.Context, questionnaireID int) ([]Questions, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	questions := []Questions{}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Find(&questions).Error
	// アンケートidの一致する質問を取る
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	return questions, nil
}

// CheckQuestionAdmin Questionの管理者か
func (*Question) CheckQuestionAdmin(ctx context.Context, userID string, questionID int) (bool, error) {
	db, err := getTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Joins("INNER JOIN administrators ON question.questionnaire_id = administrators.questionnaire_id").
		Where("question.id = ? AND administrators.user_traqid = ?", questionID, userID).
		Select("question.id").
		First(&Questions{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get question_id: %w", err)
	}

	return true, nil
}
