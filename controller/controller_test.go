package controller

import (
	"os"
	"testing"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
)

var (
	IQuestionnaire *model.Questionnaire
	IRespondent    *model.Respondent
	IResponse      *model.Response
	ITarget        *model.Target
	IQuestion      *model.Question
	IValidation    *model.Validation
	IScaleLabel    *model.ScaleLabel

	ITargetGroup        *model.TargetGroup
	ITargetUser         *model.TargetUser
	IAdministrator      *model.Administrator
	IAdministratorGroup *model.AdministratorGroup
	IAdministratorUser  *model.AdministratorUser
	IOption             *model.Option
	ITransaction        *model.Transaction
	IWebhook            *traq.Webhook

	re *Reminder
	r  *Response
	q  *Questionnaire
)

func TestMain(m *testing.M) {
	IQuestionnaire = model.NewQuestionnaire()
	IRespondent = model.NewRespondent()
	IResponse = model.NewResponse()
	ITarget = model.NewTarget()
	IQuestion = model.NewQuestion()
	IOption = model.NewOption()
	IValidation = model.NewValidation()
	IScaleLabel = model.NewScaleLabel()

	ITargetGroup = model.NewTargetGroup()
	ITargetUser = model.NewTargetUser()
	IAdministrator = model.NewAdministrator()
	IAdministratorGroup = model.NewAdministratorGroup()
	IAdministratorUser = model.NewAdministratorUser()
	ITransaction = model.NewTransaction()
	IWebhook = traq.NewWebhook()

	re = NewReminder()
	r = NewResponse(IQuestionnaire, IRespondent, IResponse, ITarget, IQuestion, IOption, IValidation, IScaleLabel)
	q = NewQuestionnaire(IQuestionnaire, ITarget, ITargetGroup, ITargetUser, IAdministrator, IAdministratorGroup, IAdministratorUser, IQuestion, IOption, IScaleLabel, IValidation, ITransaction, IRespondent, IWebhook, r, re)

	err := model.EstablishConnection(true)
	if err != nil {
		panic(err)
	}

	_, err = model.Migrate()
	if err != nil {
		panic(err)
	}

	setupSampleQuestionnaire()
	setupSampleResponse()

	code := m.Run()
	os.Exit(code)
}
