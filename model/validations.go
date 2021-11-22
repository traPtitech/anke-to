//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IValidation ValidationのRepository
type IValidation interface {
	InsertValidation(ctx context.Context, lastID int, validation Validations) error
	UpdateValidation(ctx context.Context, questionID int, validation Validations) error
	DeleteValidation(ctx context.Context, questionID int) error
	GetValidations(ctx context.Context, questionIDs []int) ([]Validations, error)
	CheckNumberValidation(validation Validations, Body string) error
	CheckTextValidation(validation Validations, Response string) error
	CheckNumberValid(MinBound, MaxBound string) error
}
