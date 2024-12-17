package handler

import "github.com/traPtitech/anke-to/controller"

type Handler struct {
	Questionnaire *controller.Questionnaire
	Response      *controller.Response
}

func NewHandler(questionnaire *controller.Questionnaire,
	response *controller.Response,
) *Handler {
	return &Handler{
		Questionnaire: questionnaire,
		Response:      response,
	}
}
