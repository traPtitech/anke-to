//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IResponse ResponseのRepository
type IResponse interface {
	InsertResponses(ctx context.Context, responseID int, responseMetas []*ResponseMeta) error
	DeleteResponse(ctx context.Context, responseID int) error
}
