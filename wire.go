//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/traPtitech/anke-to/model"
)

var (
	administratorRepository = wire.Bind(new(model.AdministratorRepository), new(model.Administrator))
	optionRepository        = wire.Bind(new(model.OptionRepository), new(model.Option))
	questionnaireRepository = wire.Bind(new(model.QuestionnaireRepository), new(model.Questionnaire))
	questionRepository      = wire.Bind(new(model.QuestionRepository), new(model.Question))
	respondentRepository    = wire.Bind(new(model.RespondentRepository), new(model.Respondent))
	responseRepository      = wire.Bind(new(model.ResponseRepository), new(model.Response))
	scaleLabelRepository    = wire.Bind(new(model.ScaleLabelRepository), new(model.ScaleLabel))
	targetRepository        = wire.Bind(new(model.TargetRepository), new(model.Target))
	validationRepository    = wire.Bind(new(model.ValidationRepository), new(model.Validation))
)
