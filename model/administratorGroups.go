package model

import "context"

// IAdministratorGroup AdministratorGroupのRepository
type IAdministratorGroup interface {
	InsertAdministratorGroups(ctx context.Context, questionnaireID int, administratorGroups []string) error
	DeleteAdministratorGroups(ctx context.Context, questionnaireID int) error
	GetAdministratorGroups(ctx context.Context, questionnaireIDs []int) ([]AdministratorGroups, error)
}
