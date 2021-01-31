package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Administrators administratorsテーブルの構造体
type Administrators struct {
	QuestionnaireID int    `sql:"type:int(11);not null;primary_key;"`
	UserTraqid      string `sql:"type:char(32);not null;primary_key;"`
}

type Administrator struct{}

// InsertAdministrators アンケートの管理者を追加
func (*Administrator) InsertAdministrators(questionnaireID int, administrators []string) error {
	var administrator Administrators
	var err error
	for _, v := range administrators {
		administrator = Administrators{
			QuestionnaireID: questionnaireID,
			UserTraqid:      v,
		}
		err = db.Create(&administrator).Error
		if err != nil {
			return fmt.Errorf("failed to insert administrators: %w", err)
		}
	}
	return nil
}

// DeleteAdministrators アンケートの管理者の削除
func (*Administrator) DeleteAdministrators(questionnaireID int) error {
	err := db.
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
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		Find(&Administrators{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get a administrator: %w", err)
	}

	return true, nil
}
