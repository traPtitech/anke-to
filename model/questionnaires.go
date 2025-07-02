//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"context"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// IQuestionnaire QuestionnaireのRepository
type IQuestionnaire interface {
	InsertQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, isPublished bool, isAnonymous bool, isDuplicateAnswerAllowed bool) (int, error)
	UpdateQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int, isPublished bool, isAnonymous bool, isDuplicateAnswerAllowed bool) error
	DeleteQuestionnaire(ctx context.Context, questionnaireID int) error
	GetQuestionnaires(ctx context.Context, userID string, sort string, search string, pageNum int, onlyTargetingMe bool, onlyAdministratedByMe bool, notOverDue bool, isDraft *bool, hasMyResponse *bool, hasMyDraft *bool) ([]QuestionnaireInfo, int, error)
	GetAdminQuestionnaires(ctx context.Context, userID string) ([]Questionnaires, error)
	GetQuestionnaireInfo(ctx context.Context, questionnaireID int) (*Questionnaires, []string, []string, []uuid.UUID, []string, []string, []uuid.UUID, []string, error)
	GetTargettedQuestionnaires(ctx context.Context, userID string, answered string, sort string) ([]TargettedQuestionnaire, error)
	GetQuestionnaireLimit(ctx context.Context, questionnaireID int) (null.Time, error)
	GetQuestionnaireLimitByResponseID(ctx context.Context, responseID int) (null.Time, error)
	GetResponseReadPrivilegeInfoByResponseID(ctx context.Context, userID string, responseID int) (*ResponseReadPrivilegeInfo, error)
	GetResponseReadPrivilegeInfoByQuestionnaireID(ctx context.Context, userID string, questionnaireID int) (*ResponseReadPrivilegeInfo, error)
	GetResponseIsAnonymousByQuestionnaireID(ctx context.Context, questionnaireID int) (bool, error)
	GetQuestionnairesInfoForReminder(ctx context.Context) ([]Questionnaires, error)
}
