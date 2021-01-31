//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"
	"github.com/traPtitech/anke-to/traq"
)

var (
	administratorRepository = wire.Bind(new(model.AdministratorRepository), new(*model.Administrator))
	optionRepository        = wire.Bind(new(model.OptionRepository), new(*model.Option))
	questionnaireRepository = wire.Bind(new(model.QuestionnaireRepository), new(*model.Questionnaire))
	questionRepository      = wire.Bind(new(model.QuestionRepository), new(*model.Question))
	respondentRepository    = wire.Bind(new(model.RespondentRepository), new(*model.Respondent))
	responseRepository      = wire.Bind(new(model.ResponseRepository), new(*model.Response))
	scaleLabelRepository    = wire.Bind(new(model.ScaleLabelRepository), new(*model.ScaleLabel))
	targetRepository        = wire.Bind(new(model.TargetRepository), new(*model.Target))
	validationRepository    = wire.Bind(new(model.ValidationRepository), new(*model.Validation))

	webhookBind = wire.Bind(new(traq.IWebhook), new(*traq.Webhook))
)

func InjectAPIServer() *router.API {
	wire.Build(
		router.NewAPI,
		router.NewMiddleware,
		router.NewQuestionnaire,
		router.NewQuestion,
		router.NewResponse,
		router.NewResult,
		router.NewUser,
		model.NewAdministrator,
		model.NewOption,
		model.NewQuestionnaire,
		model.NewQuestion,
		model.NewRespondent,
		model.NewResponse,
		model.NewScaleLabel,
		model.NewTarget,
		model.NewValidation,
		traq.NewWebhook,
		administratorRepository,
		optionRepository,
		questionnaireRepository,
		questionRepository,
		respondentRepository,
		responseRepository,
		scaleLabelRepository,
		targetRepository,
		validationRepository,
		webhookBind,
	)

	return nil
}
