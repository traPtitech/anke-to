//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// ITarget TargetのRepository
type ITarget interface {
	InsertTargets(ctx context.Context, questionnaireID int, targets []string) error
	DeleteTargets(ctx context.Context, questionnaireID int) error
	GetTargets(ctx context.Context, questionnaireIDs []int) ([]Targets, error)
	CancelTargets(ctx context.Context, questionnaireID int, targets []string) error
}
