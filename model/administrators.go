//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IAdministrator AdministratorのRepository
type IAdministrator interface {
	InsertAdministrators(ctx context.Context, questionnaireID int, administrators []string) error
	DeleteAdministrators(ctx context.Context, questionnaireID int) error
	GetAdministrators(ctx context.Context, questionnaireIDs []int) ([]Administrators, error)
	CheckQuestionnaireAdmin(ctx context.Context, userID string, questionnaireID int) (bool, error)
}
