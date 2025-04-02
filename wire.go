//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/traPtitech/anke-to/controller"
	"github.com/traPtitech/anke-to/handler"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
)

var (
	administratorBind      = wire.Bind(new(model.IAdministrator), new(*model.Administrator))
	administratorGroupBind = wire.Bind(new(model.IAdministratorGroup), new(*model.AdministratorGroup))
	administratorUserBind  = wire.Bind(new(model.IAdministratorUser), new(*model.AdministratorUser))
	optionBind             = wire.Bind(new(model.IOption), new(*model.Option))
	questionnaireBind      = wire.Bind(new(model.IQuestionnaire), new(*model.Questionnaire))
	questionBind           = wire.Bind(new(model.IQuestion), new(*model.Question))
	respondentBind         = wire.Bind(new(model.IRespondent), new(*model.Respondent))
	responseBind           = wire.Bind(new(model.IResponse), new(*model.Response))
	scaleLabelBind         = wire.Bind(new(model.IScaleLabel), new(*model.ScaleLabel))
	targetBind             = wire.Bind(new(model.ITarget), new(*model.Target))
	targetGroupBind        = wire.Bind(new(model.ITargetGroup), new(*model.TargetGroup))
	targetUserBind         = wire.Bind(new(model.ITargetUser), new(*model.TargetUser))
	validationBind         = wire.Bind(new(model.IValidation), new(*model.Validation))
	transactionBind        = wire.Bind(new(model.ITransaction), new(*model.Transaction))
	webhookBind            = wire.Bind(new(traq.IWebhook), new(*traq.Webhook))
)

func InjectAPIServer() *handler.Handler {
	wire.Build(
		handler.NewHandler,
		controller.NewResponse,
		controller.NewQuestionnaire,
		controller.NewReminder,
		controller.NewMiddleware,
		model.NewAdministrator,
		model.NewAdministratorGroup,
		model.NewAdministratorUser,
		model.NewOption,
		model.NewQuestionnaire,
		model.NewQuestion,
		model.NewRespondent,
		model.NewResponse,
		model.NewScaleLabel,
		model.NewTarget,
		model.NewTargetGroup,
		model.NewTargetUser,
		model.NewValidation,
		model.NewTransaction,
		traq.NewWebhook,
		administratorBind,
		administratorGroupBind,
		administratorUserBind,
		optionBind,
		questionnaireBind,
		questionBind,
		respondentBind,
		responseBind,
		scaleLabelBind,
		targetBind,
		targetGroupBind,
		targetUserBind,
		validationBind,
		transactionBind,
		webhookBind,
	)
	return &handler.Handler{}
}
