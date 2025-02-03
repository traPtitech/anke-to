package model

import (
	"context"
	"fmt"
)

// AdministratorUser AdministratorUserRepositoryの実装
type AdministratorUser struct{}

// NewAdministratorUser AdministratorUserRepositoryのコンストラクタ
func NewAdministratorUser() *AdministratorUser {
	return new(AdministratorUser)
}

type AdministratorUsers struct {
	QuestionnaireID int    `gorm:"type:int(11);not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

// InsertAdministratorUsers 選択したアンケート管理者（ユーザー）を追加
func (*AdministratorUser) InsertAdministratorUsers(ctx context.Context, questionnaireID int, UserTraqid []string) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if len(UserTraqid) == 0 {
		return nil
	}

	dbAdministratorUsers := make([]AdministratorUsers, 0, len(UserTraqid))
	for _, administratorUser := range UserTraqid {
		dbAdministratorUsers = append(dbAdministratorUsers, AdministratorUsers{
			QuestionnaireID: questionnaireID,
			UserTraqid:      administratorUser,
		})
	}

	err = db.Create(&dbAdministratorUsers).Error
	if err != nil {
		return fmt.Errorf("failed to insert administrator users: %w", err)
	}

	return nil
}

// DeleteAdministratorUsers 選択したアンケート管理者（ユーザー）を削除
func (*AdministratorUser) DeleteAdministratorUsers(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(AdministratorUsers{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete administrator users: %w", err)
	}

	return nil
}

// GetAdministratorUsers 選択したアンケート管理者（ユーザー）を取得
func (*AdministratorUser) GetAdministratorUsers(ctx context.Context, questionnaireIDs []int) ([]AdministratorUsers, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var administratorUsers []AdministratorUsers
	err = db.
		Where("questionnaire_id IN ?", questionnaireIDs).
		Find(&administratorUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get administrator users: %w", err)
	}

	return administratorUsers, nil
}
