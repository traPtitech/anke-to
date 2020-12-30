package model

import (
	"fmt"

	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

//Targets targetsテーブルの構造体
type Targets struct {
	QuestionnaireID int    `sql:"type:int(11);not null;primary_key;"`
	UserTraqid      string `gorm:"type:char(30);not null;primary_key;"`
}

// InsertTargets アンケートの対象を追加
func InsertTargets(questionnaireID int, targets []string) error {
	rowTargets := make([]interface{}, 0, len(targets))
	for _, target := range targets {
		rowTargets = append(rowTargets, Targets{
			QuestionnaireID: questionnaireID,
			UserTraqid:      target,
		})
	}

	err := gormbulk.BulkInsert(db, rowTargets, len(rowTargets))
	if err != nil {
		return fmt.Errorf("failed to insert target: %w", err)
	}

	return nil
}

// DeleteTargets アンケートの対象を削除
func DeleteTargets(questionnaireID int) error {
	err := db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&Targets{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete targets: %w", err)
	}

	return nil
}
