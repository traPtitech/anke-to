//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"context"
	"time"

	"gopkg.in/guregu/null.v4"
)

type MyResponseQuestionnaireInfo struct {
	QuestionnaireID     int
	Title               string
	CreatedAt           time.Time
	ModifiedAt          time.Time
	ResponseDueDateTime null.Time
	IsAnonymous         bool
	IsTargetingMe       bool
}

type MyResponseGroup struct {
	QuestionnaireInfo MyResponseQuestionnaireInfo
	Responses         []RespondentDetail
}

// IRespondent RespondentのRepository
type IRespondent interface {
	InsertRespondent(ctx context.Context, userID string, questionnaireID int, submittedAt null.Time) (int, error)
	UpdateSubmittedAt(ctx context.Context, responseID int) error
	UpdateModifiedAt(ctx context.Context, responseID int) error
	DeleteRespondent(ctx context.Context, responseID int) error
	GetRespondent(ctx context.Context, responseID int) (*Respondents, error)
	GetRespondentInfos(ctx context.Context, userID string, questionnaireIDs ...int) ([]RespondentInfo, error)
	GetRespondentDetail(ctx context.Context, responseID int) (RespondentDetail, error)
	GetRespondentDetails(ctx context.Context, questionnaireID int, sort string, onlyMyResponse bool, userID string, isDraft *bool) ([]RespondentDetail, error)
	GetMyResponseGroups(ctx context.Context, userID string, questionnaireIDs []int, isDraft *bool, pageNum int) ([]MyResponseGroup, int, error)
	GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]Respondents, error)
	GetMyResponseIDs(ctx context.Context, sort string, userID string, questionnaireIDs []int, isDraft *bool) ([]int, error)
	CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error)
}
