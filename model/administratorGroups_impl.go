package model

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
)

// AdministratorGroup AdministratorGroupRepositoryの実装
type AdministratorGroup struct{}

// NewAdministratorGroup AdministratorGroupRepositoryのコンストラクタ
func NewAdministratorGroup() *AdministratorGroup {
	return &AdministratorGroup{}
}

type AdministratorGroups struct {
	QuestionnaireID int       `gorm:"type:int(11);not null;primaryKey"`
	GroupID         uuid.UUID `gorm:"type:varchar(36);size:36;not null;primaryKey"`
}

// InsertAdministratorGroups アンケートの管理者グループを追加
func (*AdministratorGroup) InsertAdministratorGroups(ctx context.Context, questionnaireID int, groupID []uuid.UUID) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if len(groupID) == 0 {
		return nil
	}

	dbAdministratorGroups := make([]AdministratorGroups, 0, len(groupID))
	for _, administratorGroup := range groupID {
		dbAdministratorGroups = append(dbAdministratorGroups, AdministratorGroups{
			QuestionnaireID: questionnaireID,
			GroupID:         administratorGroup,
		})
	}

	err = db.Create(&dbAdministratorGroups).Error
	if err != nil {
		return fmt.Errorf("failed to insert administrator groups: %w", err)
	}

	return nil
}

// DeleteAdministratorGroups アンケートの管理者グループを削除
func (*AdministratorGroup) DeleteAdministratorGroups(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(AdministratorGroups{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete administrator groups: %w", err)
	}

	return nil
}

// GetAdministratorGroups アンケートの管理者グループを取得
func (*AdministratorGroup) GetAdministratorGroups(ctx context.Context, questionnaireIDs []int) ([]AdministratorGroups, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var administratorGroups []AdministratorGroups
	err = db.
		Where("questionnaire_id IN ?", questionnaireIDs).
		Find(&administratorGroups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get administrator groups: %w", err)
	}

	return administratorGroups, nil
}
