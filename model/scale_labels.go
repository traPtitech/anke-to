//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

// IScaleLabel ScaleLabelのRepository
type IScaleLabel interface {
	InsertScaleLabel(lastID int, label ScaleLabels) error
	UpdateScaleLabel(questionID int, label ScaleLabels) error
	DeleteScaleLabel(questionID int) error
	GetScaleLabels(questionIDs []int) ([]ScaleLabels, error)
	CheckScaleLabel(label ScaleLabels, response string) error
}
