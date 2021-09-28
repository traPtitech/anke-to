//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IScaleLabel ScaleLabelのRepository
type IScaleLabel interface {
	InsertScaleLabel(ctx context.Context, lastID int, label ScaleLabels) error
	UpdateScaleLabel(ctx context.Context, questionID int, label ScaleLabels) error
	DeleteScaleLabel(ctx context.Context, questionID int) error
	GetScaleLabels(ctx context.Context, questionIDs []int) ([]ScaleLabels, error)
	CheckScaleLabel(label ScaleLabels, response string) error
}
