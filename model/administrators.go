//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

// IAdministrator AdministratorのRepository
type IAdministrator interface {
	InsertAdministrators(questionnaireID int, administrators []string) error
	DeleteAdministrators(questionnaireID int) error
	GetAdministrators(questionnaireIDs []int) ([]Administrators, error)
	CheckQuestionnaireAdmin(userID string, questionnaireID int) (bool, error)
}
