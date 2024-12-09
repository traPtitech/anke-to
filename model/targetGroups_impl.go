package model

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
)

// TargetGroup TargetGroupsRepositoryの実装
type TargetGroup struct{}

// NewTargetGroups TargetGroupsのコンストラクター
func NewTargetGroups() *TargetGroups {
	return new(TargetGroups)
}

// TargetGroups targets_groupsテーブルの構造体
type TargetGroups struct {
	QuestionnaireID int       `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	GroupID         uuid.UUID `gorm:"type:char(36);size:36;not null;primaryKey"`
}

// InsertTargetGroups アンケートの対象としてuser_groupを追加
func (*TargetGroup) InsertTargetGroups(ctx context.Context, questionnaireID int, groupID []uuid.UUID) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if len(groupID) == 0 {
		return nil
	}

	dbTargetGroups := make([]TargetGroups, 0, len(groupID))
	for _, targetGroup := range groupID {
		dbTargetGroups = append(dbTargetGroups, TargetGroups{
			QuestionnaireID: questionnaireID,
			GroupID:         targetGroup,
		})
	}

	err = db.Create(&dbTargetGroups).Error
	if err != nil {
		return fmt.Errorf("failed to insert target groups: %w", err)
	}

	return nil
}

// GetTargetGroups アンケートの対象としてuser_groupを取得
func (*TargetGroup) GetTargetGroups(ctx context.Context, questionnaireIDs []int) ([]TargetGroups, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var targetGroups []TargetGroups
	err = db.
		Where("questionnaire_id IN ?", questionnaireIDs).
		Find(&targetGroups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get target groups: %w", err)
	}

	return targetGroups, nil
}

// DeleteTargetGroups アンケートの対象としてuser_groupを削除
func (*TargetGroup) DeleteTargetGroups(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&TargetGroups{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete target groups: %w", err)
	}

	return nil
}
