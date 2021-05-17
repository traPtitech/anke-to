//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"
	"github.com/traPtitech/anke-to/router/session"
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
	sessionBind       = wire.Bind(new(model.ISession), new(*model.Session))

	sessionStoreBind = wire.Bind(new(session.ISessionStore), new(*session.SessionStore))

	webhookBind = wire.Bind(new(traq.IWebhook), new(*traq.Webhook))
	userBind    = wire.Bind(new(traq.IUser), new(*traq.User))
)

func InjectAPIServer() (*router.API, error) {
	wire.Build(
		router.NewAPI,
		router.NewMiddleware,
		router.NewQuestionnaire,
		router.NewQuestion,
		router.NewResponse,
		router.NewResult,
		router.NewUser,
		router.NewOAuth2,
		model.NewAdministrator,
		model.NewOption,
		model.NewQuestionnaire,
		model.NewQuestion,
		model.NewRespondent,
		model.NewResponse,
		model.NewScaleLabel,
		model.NewTarget,
		model.NewValidation,
		model.NewSession,
		session.NewSessionStore,
		traq.NewUser,
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
		sessionBind,
		sessionStoreBind,
		userBind,
		webhookBind,
	)

	return nil, nil
}
