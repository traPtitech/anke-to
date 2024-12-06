package model

import (
	"context"
)

// ITargetGroup TargetGroupのRepository
type ITargetGroup interface {
	InsertTargetGroups(ctx context.Context, questionnaireID int, groupID []string) error
	GetTargetGroups(ctx context.Context, questionnaireIDs []int) ([]TargetGroups, error)
	DeleteTargetGroups(ctx context.Context, questionnaireIDs int) error
}
