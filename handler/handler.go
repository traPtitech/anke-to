package handler

import "github.com/traPtitech/anke-to/controller"

type Handler struct {
	Questionnaire *controller.Questionnaire
	Response      *controller.Response
	Reminder      *controller.Reminder
	Middleware    *controller.Middleware
}

func NewHandler(questionnaire *controller.Questionnaire,
	response *controller.Response,
	reminder *controller.Reminder,
	middleware *controller.Middleware,
) *Handler {
	reminder.ReminderInit()
	return &Handler{
		Questionnaire: questionnaire,
		Response:      response,
		Reminder:      reminder,
		Middleware:    middleware,
	}
}
