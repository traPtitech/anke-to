package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/controller"
	"github.com/traPtitech/anke-to/openapi"
)

// (GET /questionnaires)
func (h Handler) GetQuestionnaires(ctx echo.Context, params openapi.GetQuestionnairesParams) error {
	res := openapi.QuestionnaireList{}
	q := controller.NewQuestionnaire()
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	res, err = q.GetQuestionnaires(ctx, userID, params)
	if err != nil {
		ctx.Logger().Errorf("failed to get questionnaires: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaires: %w", err))
	}

	return ctx.JSON(200, res)
}

// (POST /questionnaires)
func (h Handler) PostQuestionnaire(ctx echo.Context) error {
	params := openapi.PostQuestionnaireJSONRequestBody{}
	if err := ctx.Bind(&params); err != nil {
		ctx.Logger().Errorf("failed to bind request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request body: %w", err))
	}

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
