//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"context"

	"github.com/google/uuid"
)

// ITargetGroup TargetGroupのRepository
type ITargetGroup interface {
	InsertTargetGroups(ctx context.Context, questionnaireID int, groupID []uuid.UUID) error
	GetTargetGroups(ctx context.Context, questionnaireIDs []int) ([]TargetGroups, error)
	DeleteTargetGroups(ctx context.Context, questionnaireIDs int) error
}
