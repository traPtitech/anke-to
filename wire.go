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
	administratorBind = wire.Bind(new(model.IAdministrator), new(*model.Administrator))
	optionBind        = wire.Bind(new(model.IOption), new(*model.Option))
	questionnaireBind = wire.Bind(new(model.IQuestionnaire), new(*model.Questionnaire))
	questionBind      = wire.Bind(new(model.IQuestion), new(*model.Question))
	respondentBind    = wire.Bind(new(model.IRespondent), new(*model.Respondent))
	responseBind      = wire.Bind(new(model.IResponse), new(*model.Response))
	scaleLabelBind    = wire.Bind(new(model.IScaleLabel), new(*model.ScaleLabel))
	targetBind        = wire.Bind(new(model.ITarget), new(*model.Target))
	validationBind    = wire.Bind(new(model.IValidation), new(*model.Validation))
	transactionBind   = wire.Bind(new(model.ITransaction), new(*model.Transaction))

	webhookBind = wire.Bind(new(traq.IWebhook), new(*traq.Webhook))
)

func InjectAPIServer() *handler.Handler {
	wire.Build(
		handler.NewHandler,
		controller.NewResponse,
		controller.NewQuestionnaire,
		controller.NewMiddleware,
		model.NewAdministrator,
		model.NewOption,
		model.NewQuestionnaire,
		model.NewQuestion,
		model.NewRespondent,
		model.NewResponse,
		model.NewScaleLabel,
		model.NewTarget,
		model.NewValidation,
		model.NewTransaction,
		traq.NewWebhook,
		administratorBind,
		optionBind,
		questionnaireBind,
		questionBind,
		respondentBind,
		responseBind,
		scaleLabelBind,
		targetBind,
		validationBind,
		transactionBind,
		webhookBind,
	)
	return &handler.Handler{}
}
