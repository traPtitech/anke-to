//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"context"

	"gopkg.in/guregu/null.v4"
)

// IQuestionnaire QuestionnaireのRepository
type IQuestionnaire interface {
	InsertQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, isPublished bool) (int, error)
	UpdateQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int, isPublished bool) error
	DeleteQuestionnaire(ctx context.Context, questionnaireID int) error
	GetQuestionnaires(ctx context.Context, userID string, sort string, search string, pageNum int, onlyTargetingMe bool, onlyAdministratedByMe bool) ([]QuestionnaireInfo, int, error)
	GetAdminQuestionnaires(ctx context.Context, userID string) ([]Questionnaires, error)
	GetQuestionnaireInfo(ctx context.Context, questionnaireID int) (*Questionnaires, []string, []string, []string, []string, []string, error)
	GetTargettedQuestionnaires(ctx context.Context, userID string, answered string, sort string) ([]TargettedQuestionnaire, error)
	GetQuestionnaireLimit(ctx context.Context, questionnaireID int) (null.Time, error)
	GetQuestionnaireLimitByResponseID(ctx context.Context, responseID int) (null.Time, error)
	GetResponseReadPrivilegeInfoByResponseID(ctx context.Context, userID string, responseID int) (*ResponseReadPrivilegeInfo, error)
	GetResponseReadPrivilegeInfoByQuestionnaireID(ctx context.Context, userID string, questionnaireID int) (*ResponseReadPrivilegeInfo, error)
}
