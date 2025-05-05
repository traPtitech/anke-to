//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// ITarget TargetのRepository
type ITarget interface {
	InsertTargets(ctx context.Context, questionnaireID int, targets []string) error
	DeleteTargets(ctx context.Context, questionnaireID int) error
	GetTargets(ctx context.Context, questionnaireIDs []int) ([]Targets, error)
	IsTargetingMe(ctx context.Context, quesionnairID int, userID string) (bool, error)
	GetTargetRemindStatus(ctx context.Context, questionnaireID int, target string) (bool, error)
	UpdateTargetsRemindStatus(ctx context.Context, questionnaireID int, targets []string, remindStatus bool) error
}
