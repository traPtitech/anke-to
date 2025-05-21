package model

import (
	"context"
	"fmt"
)

// TargetUser TargetUsersRepositoryの実装
type TargetUser struct{}

// NewTargetUsers TargetUsersのコンストラクター
func NewTargetUser() *TargetUser {
	return new(TargetUser)
}

// TargetUsers targets_usersテーブルの構造体
type TargetUsers struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

// InsertTargetUsers 選択したアンケート対象者（ユーザー）を追加
func (*TargetUser) InsertTargetUsers(ctx context.Context, questionnaireID int, traqID []string) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if len(traqID) == 0 {
		return nil
	}

	dbTargetUsers := make([]TargetUsers, 0, len(traqID))
	for _, targetUser := range traqID {
		dbTargetUsers = append(dbTargetUsers, TargetUsers{
			QuestionnaireID: questionnaireID,
			UserTraqid:      targetUser,
		})
	}

	err = db.Create(&dbTargetUsers).Error
	if err != nil {
		return fmt.Errorf("failed to insert target users: %w", err)
	}

	return nil
}

// GetTargetUsers 選択したアンケート対象者（ユーザー）を取得
func (*TargetUser) GetTargetUsers(ctx context.Context, questionnaireIDs []int) ([]TargetUsers, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var TargetUsers []TargetUsers
	err = db.
		Where("questionnaire_id IN ?", questionnaireIDs).
		Find(&TargetUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get target users: %w", err)
	}

	return TargetUsers, nil
}

// DeleteTargetUsers 選択したアンケート対象者（ユーザー）を削除
func (*TargetUser) DeleteTargetUsers(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&TargetUsers{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete target users: %w", err)
	}

	return nil
}
