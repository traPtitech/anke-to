//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

// ITarget TargetのRepository
type ITarget interface {
	InsertTargets(questionnaireID int, targets []string) error
	DeleteTargets(questionnaireID int) error
	GetTargets(questionnaireIDs []int) ([]Targets, error)
}
