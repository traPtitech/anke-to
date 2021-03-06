package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-sql-driver/mysql"
)

// Question QuestionRepositoryの実装
type Question struct{}

// NewQuestion Questionのコンストラクター
func NewQuestion() *Question {
	return new(Question)
}

//Questions questionテーブルの構造体
type Questions struct {
	ID              int            `json:"id"                  gorm:"type:int(11) AUTO_INCREMENT NOT NULL PRIMARY KEY;"`
	QuestionnaireID int            `json:"questionnaireID"     gorm:"type:int(11);default:NULL;"`
	PageNum         int            `json:"page_num"            gorm:"type:int(11) NOT NULL;"`
	QuestionNum     int            `json:"question_num"        gorm:"type:int(11) NOT NULL;"`
	Type            string         `json:"type"                gorm:"type:char(20) NOT NULL;"`
	Body            string         `json:"body"                gorm:"type:text;default:NULL;"`
	IsRequired      bool           `json:"is_required"         gorm:"type:tinyint(4) NOT NULL;default:0;"`
	DeletedAt       mysql.NullTime `json:"deleted_at"          gorm:"type:timestamp NULL;default:NULL;"`
	CreatedAt       time.Time      `json:"created_at"          gorm:"type:timestamp NOT NULL;default:CURRENT_TIMESTAMP;"`
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
func (*Question) InsertQuestion(questionnaireID int, pageNum int, questionNum int, questionType string,
	body string, isRequired bool) (int, error) {
	question := Questions{
		QuestionnaireID: questionnaireID,
		PageNum:         pageNum,
		QuestionNum:     questionNum,
		Type:            questionType,
		Body:            body,
		IsRequired:      isRequired,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&question).Error
		if err != nil {
			return fmt.Errorf("failed to insert a question record: %w", err)
		}

		err = tx.
			Select("id").
			Last(&question).Error
		if err != nil {
			return fmt.Errorf("failed to get the last question record: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed in transaction: %w", err)
	}

	return question.ID, nil
}

//UpdateQuestion 質問の修正
func (*Question) UpdateQuestion(questionnaireID int, pageNum int, questionNum int, questionType string,
	body string, isRequired bool, questionID int) error {
	question := map[string]interface{}{
		"questionnaire_id": questionnaireID,
		"page_num":         pageNum,
		"question_num":     questionNum,
		"type":             questionType,
		"body":             body,
		"is_required":      isRequired,
	}

	err := db.
		Model(&Questions{}).
		Where("id = ?", questionID).
		Update(question).Error
	if err != nil {
		return fmt.Errorf("failed to update a question record: %w", err)
	}

	return nil
}

//DeleteQuestion 質問の削除
func (*Question) DeleteQuestion(questionID int) error {
	result := db.
		Where("id = ?", questionID).
		Delete(&Questions{})
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to delete a question record: %w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

//GetQuestions 質問一覧の取得
func (*Question) GetQuestions(questionnaireID int) ([]Questions, error) {
	questions := []Questions{}

	err := db.
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
func (*Question) CheckQuestionAdmin(userID string, questionID int) (bool, error) {
	err := db.
		Table("question").
		Joins("INNER JOIN administrators ON question.questionnaire_id = administrators.questionnaire_id").
		Where("question.id = ? AND administrators.user_traqid = ?", questionID, userID).
		Select("question.id").
		Find(&Questions{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get question_id: %w", err)
	}

	return true, nil
}
