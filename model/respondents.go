//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "gopkg.in/guregu/null.v3"

// RespondentRepository RespondentのRepository
type RespondentRepository interface {
	InsertRespondent(userID string, questionnaireID int, submitedAt null.Time) (int, error)
	UpdateSubmittedAt(responseID int) error
	DeleteRespondent(userID string, responseID int) error
	GetRespondentInfos(userID string, questionnaireIDs ...int) ([]RespondentInfo, error)
	GetRespondentDetail(responseID int) (RespondentDetail, error)
	GetRespondentDetails(questionnaireID int, sort string) ([]RespondentDetail, error)
	GetRespondentsUserIDs(questionnaireIDs []int) ([]Respondents, error)
	CheckRespondent(userID string, questionnaireID int) (bool, error)
	CheckRespondentByResponseID(userID string, responseID int) (bool, error)
}
