package model

// ResponseRepository ResponseのRepository
type ResponseRepository interface {
	InsertResponses(responseID int, responseMetas []*ResponseMeta) error
	DeleteResponse(responseID int) error
}
