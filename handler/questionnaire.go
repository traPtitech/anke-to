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
	validate, err := getValidator(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get validator: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get validator: %w", err))
	}

	err = validate.StructCtx(ctx.Request().Context(), params)
	if err != nil {
		ctx.Logger().Errorf("failed to validate request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to validate request body: %w", err))
	}

	res := openapi.QuestionnaireDetail{}
	q := controller.NewQuestionnaire()
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	res, err = q.PostQuestionnaire(ctx, userID, params)
	if err != nil {
		ctx.Logger().Errorf("failed to post questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to post questionnaire: %w", err))
	}

	return ctx.JSON(201, res)
}

// (GET /questionnaires/{questionnaireID})
func (h Handler) GetQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.QuestionnaireDetail{}
	q := controller.NewQuestionnaire()
	res, err := q.GetQuestionnaire(ctx, questionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to get questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire: %w", err))
	}
	return ctx.JSON(200, res)
}

// (PATCH /questionnaires/{questionnaireID})
func (h Handler) EditQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	params := openapi.EditQuestionnaireJSONRequestBody{}
	if err := ctx.Bind(&params); err != nil {
		ctx.Logger().Errorf("failed to bind request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request body: %w", err))
	}

	q := controller.NewQuestionnaire()
	err := q.EditQuestionnaire(ctx, questionnaireID, params)
	if err != nil {
		ctx.Logger().Errorf("failed to edit questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to edit questionnaire: %w", err))
	}

	return ctx.NoContent(200)
}

// (DELETE /questionnaires/{questionnaireID})
func (h Handler) DeleteQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	q := controller.NewQuestionnaire()
	err := q.DeleteQuestionnaire(ctx, questionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to delete questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete questionnaire: %w", err))
	}

	return ctx.NoContent(200)
}

// (GET /questionnaires/{questionnaireID}/myRemindStatus)
func (h Handler) GetQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.QuestionnaireIsRemindEnabled{}
	q := controller.NewQuestionnaire()
	status, err := q.GetQuestionnaireMyRemindStatus(ctx, questionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to get questionnaire my remind status: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire my remind status: %w", err))
	}
	res.IsRemindEnabled = status

	return ctx.JSON(200, res)
}

// (PATCH /questionnaires/{questionnaireID}/myRemindStatus)
func (h Handler) EditQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	q := controller.NewQuestionnaire()
	err := q.EditQuestionnaireMyRemindStatus(ctx, questionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to edit questionnaire my remind status: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to edit questionnaire my remind status: %w", err))
	}
	return ctx.NoContent(200)
}

// (GET /questionnaires/{questionnaireID}/responses)
func (h Handler) GetQuestionnaireResponses(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath, params openapi.GetQuestionnaireResponsesParams) error {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}
	q := controller.NewQuestionnaire()
	res, err := q.GetQuestionnaireResponses(ctx, questionnaireID, params, userID)

	return ctx.JSON(200, res)
}

// (POST /questionnaires/{questionnaireID}/responses)
func (h Handler) PostQuestionnaireResponse(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.Response{}
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	params := openapi.PostQuestionnaireResponseJSONRequestBody{}
	if err := ctx.Bind(&params); err != nil {
		ctx.Logger().Errorf("failed to bind request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request body: %w", err))
	}
	validate, err := getValidator(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get validator: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get validator: %w", err))
	}

	err = validate.StructCtx(ctx.Request().Context(), params)
	if err != nil {
		ctx.Logger().Errorf("failed to validate request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to validate request body: %w", err))
	}

	q := controller.NewQuestionnaire()
	res, err = q.PostQuestionnaireResponse(ctx, questionnaireID, params, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to post questionnaire response: %+v", err)
		return err
	}

	return ctx.JSON(201, res)
}

// (GET /questionnaires/{questionnaireID}/result)
func (h Handler) GetQuestionnaireResult(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.Result{}
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	q := controller.NewQuestionnaire()
	res, err = q.GetQuestionnaireResult(ctx, questionnaireID, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get questionnaire result: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire result: %w", err))
	}

	return ctx.JSON(200, res)
}
