// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"
	"github.com/traPtitech/anke-to/traq"
)

import (
	_ "net/http/pprof"
)

// Injectors from wire.go:

func InjectAPIServer() *router.API {
	administrator := model.NewAdministrator()
	respondent := model.NewRespondent()
	question := model.NewQuestion()
	questionnaire := model.NewQuestionnaire()
	middleware := router.NewMiddleware(administrator, respondent, question, questionnaire)
	target := model.NewTarget()
	option := model.NewOption()
	scaleLabel := model.NewScaleLabel()
	validation := model.NewValidation()
	transaction := model.NewTransaction()
	webhook := traq.NewWebhook()
	routerQuestionnaire := router.NewQuestionnaire(questionnaire, target, administrator, question, option, scaleLabel, validation, transaction, webhook)
	routerQuestion := router.NewQuestion(validation, question, option, scaleLabel)
	response := model.NewResponse()
	routerResponse := router.NewResponse(questionnaire, validation, scaleLabel, respondent, response)
	result := router.NewResult(respondent, questionnaire, administrator)
	user := router.NewUser(respondent, questionnaire, target, administrator)
	api := router.NewAPI(middleware, routerQuestionnaire, routerQuestion, routerResponse, result, user)
	return api
}

// wire.go:

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
