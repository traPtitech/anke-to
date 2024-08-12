package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/controller"
	"github.com/traPtitech/anke-to/openapi"
)

// (GET /responses/myResponses)
func (h Handler) GetMyResponses(ctx echo.Context, params openapi.GetMyResponsesParams) error {
	res := openapi.ResponsesWithQuestionnaireInfo{}
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	r := controller.NewResponse()
	res, err = r.GetMyResponses(ctx, params, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get my responses: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get my responses: %+w", err))
	}
	return ctx.JSON(200, res)
}

// (DELETE /responses/{responseID})
func (h Handler) DeleteResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	return ctx.NoContent(200)
}

// (GET /responses/{responseID})
func (h Handler) GetResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	res := openapi.Response{}

	r := controller.NewResponse()
	res, err := r.GetResponse(ctx, responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to get my response: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response: %+w", err))
	}
	return ctx.JSON(200, res)
}

// (PATCH /responses/{responseID})
func (h Handler) EditResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	return ctx.NoContent(200)
}