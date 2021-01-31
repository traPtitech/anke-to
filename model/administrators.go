package model

// AdministratorRepository AdministratorのRepository
type AdministratorRepository interface {
	InsertAdministrators(questionnaireID int, administrators []string) error
	DeleteAdministrators(questionnaireID int) error
	GetAdministrators(questionnaireIDs []int) ([]Administrators, error)
	CheckQuestionnaireAdmin(userID string, questionnaireID int) (bool, error)
}
