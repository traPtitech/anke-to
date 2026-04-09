package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReminderTarget struct{}

func NewReminderTarget() *ReminderTarget {
	return new(ReminderTarget)
}

type ReminderTargets struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
	IsCanceled      bool   `gorm:"type:tinyint(1);not null;default:0"`
}

func (*ReminderTarget) GetReminderTarget(ctx context.Context, questionnaireID int, userID string) (*ReminderTargets, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	reminderTarget := ReminderTargets{}
	err = db.
		Where("questionnaire_id = ? AND user_traqid = ?", questionnaireID, userID).
		First(&reminderTarget).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get reminder target: %w", err)
	}

	return &reminderTarget, nil
}

func (*ReminderTarget) GetReminderTargets(ctx context.Context, questionnaireID int) ([]ReminderTargets, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	reminderTargets := []ReminderTargets{}
	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Find(&reminderTargets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder targets: %w", err)
	}

	return reminderTargets, nil
}

func (*ReminderTarget) UpsertReminderTarget(ctx context.Context, questionnaireID int, userID string, isCanceled bool) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "questionnaire_id"}, {Name: "user_traqid"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_canceled"}),
		}).
		Create(&ReminderTargets{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
			IsCanceled:      isCanceled,
		}).Error
	if err != nil {
		return fmt.Errorf("failed to upsert reminder target: %w", err)
	}

	return nil
}

func (*ReminderTarget) DeleteReminderTargets(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&ReminderTargets{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete reminder targets: %w", err)
	}

	return nil
}
