package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/openapi"
)

// (GET /questionnaires)
func (h Handler) GetQuestionnaires(ctx echo.Context, params openapi.GetQuestionnairesParams) error {
	res := openapi.QuestionnaireList{}

	return ctx.JSON(200, res)
}

// (POST /questionnaires)
func (h Handler) PostQuestionnaire(ctx echo.Context) error {
	res := openapi.QuestionnaireDetail{}

	return ctx.JSON(200, res)
}

// (GET /questionnaires/{questionnaireID})
func (h Handler) GetQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.QuestionnaireDetail{}

	return ctx.JSON(200, res)
}

// (PATCH /questionnaires/{questionnaireID})
func (h Handler) EditQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	return ctx.NoContent(200)
}

// (DELETE /questionnaires/{questionnaireID})
func (h Handler) DeleteQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	return ctx.NoContent(200)
}

// (GET /questionnaires/{questionnaireID}/myRemindStatus)
func (h Handler) GetQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.QuestionnaireIsRemindEnabled{}

	return ctx.JSON(200, res)
}

// (PATCH /questionnaires/{questionnaireID}/myRemindStatus)
func (h Handler) EditQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	return ctx.NoContent(200)
}

// (GET /questionnaires/{questionnaireID}/responses)
func (h Handler) GetQuestionnaireResponses(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath, params openapi.GetQuestionnaireResponsesParams) error {
	res := openapi.Responses{}

	return ctx.JSON(200, res)
}

// (POST /questionnaires/{questionnaireID}/responses)
func (h Handler) PostQuestionnaireResponse(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.Response{}

	return ctx.JSON(201, res)
}

// (GET /questionnaires/{questionnaireID}/result)
func (h Handler) GetQuestionnaireResult(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.Result{}

	return ctx.JSON(200, res)
}
