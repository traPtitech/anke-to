package handler

import (
	"github.com/traPtitech/anke-to/controller"
	traqAPI "github.com/traPtitech/anke-to/traq"
)

type Handler struct {
	Questionnaire *controller.Questionnaire
	Response      *controller.Response
	Reminder      *controller.Reminder
	Middleware    *controller.Middleware
	TraqClient    *traqAPI.APIClient
}

func NewHandler(questionnaire *controller.Questionnaire,
	response *controller.Response,
	reminder *controller.Reminder,
	middleware *controller.Middleware,
	traqClient *traqAPI.APIClient,
) *Handler {
	reminder.ReminderInit()
	return &Handler{
		Questionnaire: questionnaire,
		Response:      response,
		Reminder:      reminder,
		Middleware:    middleware,
		TraqClient:    traqClient,
	}
}
