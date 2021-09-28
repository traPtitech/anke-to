//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IOption OptionのRepository
type IOption interface {
	InsertOption(ctx context.Context, lastID int, num int, body string) error
	UpdateOptions(ctx context.Context, options []string, questionID int) error
	DeleteOptions(ctx context.Context, questionID int) error
	GetOptions(ctx context.Context, questionIDs []int) ([]Options, error)
}
