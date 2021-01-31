package model

// ScaleLabelRepository ScaleLabelのRepository
type ScaleLabelRepository interface {
	InsertScaleLabel(lastID int, label ScaleLabels) error
	UpdateScaleLabel(questionID int, label ScaleLabels) error
	DeleteScaleLabel(questionID int) error
	GetScaleLabels(questionIDs []int) ([]ScaleLabels, error)
	CheckScaleLabel(label ScaleLabels, response string) error
}
