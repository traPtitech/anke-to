//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "gopkg.in/guregu/null.v3"

// IQuestionnaire QuestionnaireのRepository
type IQuestionnaire interface {
	InsertQuestionnaire(title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error)
	UpdateQuestionnaire(title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error
	DeleteQuestionnaire(questionnaireID int) error
	GetQuestionnaires(userID string, sort string, search string, pageNum int, nontargeted bool) ([]QuestionnaireInfo, int, error)
	GetAdminQuestionnaires(userID string) ([]Questionnaires, error)
	GetQuestionnaireInfo(questionnaireID int) (*Questionnaires, []string, []string, []string, error)
	GetTargettedQuestionnaires(userID string, answered string, sort string) ([]TargettedQuestionnaire, error)
	GetQuestionnaireLimit(questionnaireID int) (null.Time, error)
	GetQuestionnaireLimitByResponseID(responseID int) (null.Time, error)
	GetResShared(questionnaireID int) (string, error)
	GetResponseReadPrivilegeInfoByResponseID(userID string, responseID int) (*ResponseReadPrivilegeInfo, error)
}
