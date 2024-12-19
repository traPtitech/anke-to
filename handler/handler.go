package handler

import "github.com/traPtitech/anke-to/controller"

type Handler struct {
	Questionnaire *controller.Questionnaire
	Response      *controller.Response
	Middleware    *controller.Middleware
}

func NewHandler(questionnaire *controller.Questionnaire,
	response *controller.Response,
	middleware *controller.Middleware,
) *Handler {
	return &Handler{
		Questionnaire: questionnaire,
		Response:      response,
		Middleware:    middleware,
	}
}
