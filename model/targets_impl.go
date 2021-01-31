package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

// Target TargetRepositoryの実装
type Target struct{}

// NewTarget Targetのコンストラクター
func NewTarget() *Target {
	return new(Target)
}

//Targets targetsテーブルの構造体
type Targets struct {
	QuestionnaireID int    `sql:"type:int(11);not null;primary_key;"`
	UserTraqid      string `gorm:"type:char(30);not null;primary_key;"`
}

// InsertTargets アンケートの対象を追加
func (*Target) InsertTargets(questionnaireID int, targets []string) error {
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
func (*Target) DeleteTargets(questionnaireID int) error {
	err := db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&Targets{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete targets: %w", err)
	}

	return nil
}

// GetTargets アンケートの対象一覧を取得
func (*Target) GetTargets(questionnaireIDs []int) ([]Targets, error) {
	targets := []Targets{}
	err := db.
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Find(&targets).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	return targets, nil
}
