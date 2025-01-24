package model

import (
	"context"
)

// ITargetUser TargetUserのRepository
type ITargetUser interface {
	InsertTargetUsers(ctx context.Context, questionnaireID int, traqID []string) error
	GetTargetUsers(ctx context.Context, questionnaireIDs []int) ([]TargetUsers, error)
	DeleteTargetUsers(ctx context.Context, questionnaireIDs int) error
}
