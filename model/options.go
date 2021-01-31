//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

// IOption OptionのRepository
type IOption interface {
	InsertOption(lastID int, num int, body string) error
	UpdateOptions(options []string, questionID int) error
	DeleteOptions(questionID int) error
	GetOptions(questionIDs []int) ([]Options, error)
}
