package router

import (
	"github.com/traPtitech/anke-to/model"
)

// Result Resultの構造体
type Result struct {
	model.IRespondent
	model.IQuestionnaire
	model.IAdministrator
}

// NewResult Resultのコンストラクタ
func NewResult(respondent model.IRespondent, questionnaire model.IQuestionnaire, administrator model.IAdministrator) *Result {
	return &Result{
		IRespondent:    respondent,
		IQuestionnaire: questionnaire,
		IAdministrator: administrator,
	}
}
