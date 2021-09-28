package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Administrator AdministratorRepositoryの実装
type Administrator struct{}

// NewAdministrator Administratorのコンストラクター
func NewAdministrator() *Administrator {
	return new(Administrator)
}

// Administrators administratorsテーブルの構造体
type Administrators struct {
	QuestionnaireID int    `gorm:"type:int(11);not null;primaryKey"`
	UserTraqid      string `gorm:"type:char(30);size:30;not null;primaryKey"`
}

// InsertAdministrators アンケートの管理者を追加
func (*Administrator) InsertAdministrators(questionnaireID int, administrators []string) error {
	dbAdministrators := make([]Administrators, 0, len(administrators))

	if len(administrators) == 0 {
		return nil
	}

	var err error
	for _, v := range administrators {
		dbAdministrators = append(dbAdministrators, Administrators{
			QuestionnaireID: questionnaireID,
			UserTraqid:      v,
		})
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Create(&dbAdministrators).Error
	if err != nil {
		return fmt.Errorf("failed to insert administrators: %w", err)
	}

	return nil
}

// DeleteAdministrators アンケートの管理者の削除
func (*Administrator) DeleteAdministrators(questionnaireID int) error {
	err := db.
		Session(&gorm.Session{NewDB: true}).
		Where("questionnaire_id = ?", questionnaireID).
		Delete(Administrators{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete administrators: %w", err)
	}

	return nil
}

// GetAdministrators アンケートの管理者を取得
func (*Administrator) GetAdministrators(questionnaireIDs []int) ([]Administrators, error) {
	administrators := []Administrators{}
	err := db.
		Session(&gorm.Session{NewDB: true}).
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Find(&administrators).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get administrators: %w", err)
	}

	return administrators, nil
}

// CheckQuestionnaireAdmin 自分がアンケートの管理者か判定
func (*Administrator) CheckQuestionnaireAdmin(userID string, questionnaireID int) (bool, error) {
	err := db.
		Session(&gorm.Session{NewDB: true}).
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		First(&Administrators{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get a administrator: %w", err)
	}

	return true, nil
}
