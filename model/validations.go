package model

// ValidationRepository ValidationのRepository
type ValidationRepository interface {
	InsertValidation(lastID int, validation Validations) error
	UpdateValidation(questionID int, validation Validations) error
	DeleteValidation(questionID int) error
	GetValidations(qustionIDs []int) ([]Validations, error)
	CheckNumberValidation(validation Validations, Body string) error
	CheckTextValidation(validation Validations, Response string) error
	CheckNumberValid(MinBound, MaxBound string) error
}
