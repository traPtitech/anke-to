//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IAdministratorUser AdministratorUserのRepository
type IAdministratorUser interface {
	InsertAdministratorUsers(ctx context.Context, questionnaireID int, traqID []string) error
	DeleteAdministratorUsers(ctx context.Context, questionnaireID int) error
	GetAdministratorUsers(ctx context.Context, questionnaireIDs []int) ([]AdministratorUsers, error)
}
