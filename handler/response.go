package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/openapi"
)

// (GET /responses/myResponses)
func (h Handler) GetMyResponses(ctx echo.Context, params openapi.GetMyResponsesParams) error {
	userID, err := h.Middleware.GetUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	res, err := h.Response.GetMyResponses(ctx, params, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get my responses: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get my responses: %w", err))
	}
	return ctx.JSON(200, res)
}

// (DELETE /responses/{responseID})
func (h Handler) DeleteResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	err := h.Response.DeleteResponse(ctx, responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to delete response: %+v", err)
		return err
	}

	return ctx.NoContent(200)
}

// (GET /responses/{responseID})
func (h Handler) GetResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	res, err := h.Response.GetResponse(ctx, responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to get response: %+v", err)
		return err
	}
	return ctx.JSON(200, res)
}

// (PATCH /responses/{responseID})
func (h Handler) EditResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	req := openapi.EditResponseJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Errorf("failed to bind Responses: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind Responses: %w", err))
	}

	validate, err := h.Middleware.GetValidator(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get validator: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get validator: %w", err))
	}

	err = validate.Struct(req)
	if err != nil {
		ctx.Logger().Errorf("failed to validate request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to validate request body: %w", err))
	}

	err = h.Response.EditResponse(ctx, responseID, req)
	if err != nil {
		ctx.Logger().Errorf("failed to edit response: %+v", err)
		return err
	}

	return ctx.NoContent(200)
}
