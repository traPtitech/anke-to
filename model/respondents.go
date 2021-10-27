//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"context"

	"gopkg.in/guregu/null.v4"
)

// IRespondent RespondentのRepository
type IRespondent interface {
	InsertRespondent(ctx context.Context, userID string, questionnaireID int, submitedAt null.Time) (int, error)
	UpdateSubmittedAt(ctx context.Context, responseID int) error
	DeleteRespondent(ctx context.Context, responseID int) error
	GetRespondent(ctx context.Context, responseID int) (*Respondents, error)
	GetRespondentInfos(ctx context.Context, userID string, questionnaireIDs ...int) ([]RespondentInfo, error)
	GetRespondentDetail(ctx context.Context, responseID int) (RespondentDetail, error)
	GetRespondentDetails(ctx context.Context, questionnaireID int, sort string) ([]RespondentDetail, error)
	GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]Respondents, error)
	CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error)
}
