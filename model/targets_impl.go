package model

import (
	"context"
	"fmt"
)

// Target TargetRepositoryの実装
type Target struct{}

// NewTarget Targetのコンストラクター
func NewTarget() *Target {
	return new(Target)
}

// Targets targetsテーブルの構造体
type Targets struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
	IsCanceled      bool   `gorm:"type:tinyint(1);not null;default:0"`
}

// InsertTargets アンケートの対象を追加
func (*Target) InsertTargets(ctx context.Context, questionnaireID int, targets []string) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if len(targets) == 0 {
		return nil
	}

	dbTargets := make([]Targets, 0, len(targets))
	for _, target := range targets {
		dbTargets = append(dbTargets, Targets{
			QuestionnaireID: questionnaireID,
			UserTraqid:      target,
			IsCanceled:      false,
		})
	}

	err = db.Create(&dbTargets).Error
	if err != nil {
		return fmt.Errorf("failed to insert targets: %w", err)
	}

	return nil
}

// DeleteTargets アンケートの対象を削除
func (*Target) DeleteTargets(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&Targets{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete targets: %w", err)
	}

	return nil
}

// GetTargets アンケートの対象一覧を取得
func (*Target) GetTargets(ctx context.Context, questionnaireIDs []int) ([]Targets, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	targets := []Targets{}
	err = db.
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Find(&targets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	return targets, nil
}

func (*Target) IsTargetingMe(ctx context.Context, questionnairID int, userID string) (bool, error) {
	db, err := getTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %w", err)
	}

	var count int64
	err = db.
		Model(&Targets{}).
		Where("questionnaire_id = ? AND user_traqid = ?", questionnairID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to get targets which are targeting me: %w", err)
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (*Target) GetTargetRemindStatus(ctx context.Context, questionnaireID int, target string) (bool, error) {
	db, err := getTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %w", err)
	}

	var cancelStatus bool
	err = db.
		Model(&Targets{}).
		Select("is_canceled").
		Where("questionnaire_id = ? AND user_traqid = ?", questionnaireID, target).
		Take(&cancelStatus).Error
	if err != nil {
		return false, fmt.Errorf("failed to get targets remind status: %w", err)
	}

	return !cancelStatus, nil
}

func (*Target) UpdateTargetsRemindStatus(ctx context.Context, questionnaireID int, targets []string, remindStatus bool) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Model(&Targets{}).
		Where("questionnaire_id = ? AND user_traqid IN (?)", questionnaireID, targets).
		Update("is_canceled", !remindStatus).Error
	if err != nil {
		return fmt.Errorf("failed to cancel targets: %w", err)
	}

	return nil
}
