package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Target TargetRepositoryの実装
type Target struct{}

// NewTarget Targetのコンストラクター
func NewTarget() *Target {
	return new(Target)
}

//Targets targetsテーブルの構造体
type Targets struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:char(30);size:30;not null;primaryKey"`
}

// InsertTargets アンケートの対象を追加
func (*Target) InsertTargets(questionnaireID int, targets []string) error {
	if len(targets) == 0 {
		return nil
	}

	dbTargets := make([]Targets, 0, len(targets))
	for _, target := range targets {
		dbTargets = append(dbTargets, Targets{
			QuestionnaireID: questionnaireID,
			UserTraqid:      target,
		})
	}

	err := db.
		Session(&gorm.Session{NewDB: true}).
		Create(&dbTargets).Error
	if err != nil {
		return fmt.Errorf("failed to insert targets: %w", err)
	}

	return nil
}

// DeleteTargets アンケートの対象を削除
func (*Target) DeleteTargets(questionnaireID int) error {
	err := db.
		Session(&gorm.Session{NewDB: true}).
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
		Session(&gorm.Session{NewDB: true}).
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Find(&targets).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	return targets, nil
}
