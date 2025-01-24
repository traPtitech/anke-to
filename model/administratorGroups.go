package model

import (
	"context"

	"github.com/google/uuid"
)

// IAdministratorGroup AdministratorGroupのRepository
type IAdministratorGroup interface {
	InsertAdministratorGroups(ctx context.Context, questionnaireID int, groupID []uuid.UUID) error
	DeleteAdministratorGroups(ctx context.Context, questionnaireID int) error
	GetAdministratorGroups(ctx context.Context, questionnaireIDs []int) ([]AdministratorGroups, error)
}
