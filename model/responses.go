//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

// ResponseRepository ResponseのRepository
type ResponseRepository interface {
	InsertResponses(responseID int, responseMetas []*ResponseMeta) error
	DeleteResponse(responseID int) error
}
