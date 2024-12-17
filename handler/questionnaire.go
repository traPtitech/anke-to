package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
)

// (GET /questionnaires)
func (h Handler) GetQuestionnaires(ctx echo.Context, params openapi.GetQuestionnairesParams) error {
	res := openapi.QuestionnaireList{}
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	res, err = h.Questionnaire.GetQuestionnaires(ctx, userID, params)
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
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	res, err = h.Questionnaire.PostQuestionnaire(ctx, userID, params)
	if err != nil {
		ctx.Logger().Errorf("failed to post questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to post questionnaire: %w", err))
	}

	return ctx.JSON(201, res)
}

// (GET /questionnaires/{questionnaireID})
func (h Handler) GetQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.QuestionnaireDetail{}
	res, err := h.Questionnaire.GetQuestionnaire(ctx, questionnaireID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("questionnaire not found: %w", err))
		}
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

	err := h.Questionnaire.EditQuestionnaire(ctx, questionnaireID, params)
	if err != nil {
		ctx.Logger().Errorf("failed to edit questionnaire: %+v", err)
		return err
	}

	return ctx.NoContent(200)
}

// (DELETE /questionnaires/{questionnaireID})
func (h Handler) DeleteQuestionnaire(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	err := h.Questionnaire.DeleteQuestionnaire(ctx, questionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to delete questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete questionnaire: %w", err))
	}

	return ctx.NoContent(200)
}

// (GET /questionnaires/{questionnaireID}/myRemindStatus)
func (h Handler) GetQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	res := openapi.QuestionnaireIsRemindEnabled{}
	status, err := h.Questionnaire.GetQuestionnaireMyRemindStatus(ctx, questionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to get questionnaire my remind status: %+v", err)
		return err
	}
	res.IsRemindEnabled = status

	return ctx.JSON(200, res)
}

// (PATCH /questionnaires/{questionnaireID}/myRemindStatus)
func (h Handler) EditQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID openapi.QuestionnaireIDInPath) error {
	params := openapi.EditQuestionnaireMyRemindStatusJSONRequestBody{}
	if err := ctx.Bind(&params); err != nil {
		ctx.Logger().Errorf("failed to bind request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request body: %w", err))
	}

	err := h.Questionnaire.EditQuestionnaireMyRemindStatus(ctx, questionnaireID, params.IsRemindEnabled)
	if err != nil {
		ctx.Logger().Errorf("failed to edit questionnaire my remind status: %+v", err)
		return err
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
	res, err := h.Questionnaire.GetQuestionnaireResponses(ctx, questionnaireID, params, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get questionnaire responses: %+v", err)
		return err
	}

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

	res, err = h.Questionnaire.PostQuestionnaireResponse(ctx, questionnaireID, params, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to post questionnaire response: %+v", err)
		return err
	}

	return ctx.JSON(201, res)
}
