package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/openapi"
)

// (GET /responses/myResponses)
func (h Handler) GetMyResponses(ctx echo.Context, params openapi.GetMyResponsesParams) error {
	res := []openapi.ResponsesWithQuestionnaireInfo{}

	return ctx.JSON(200, res)
}

// (DELETE /responses/{responseID})
func (h Handler) DeleteResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	return ctx.NoContent(200)
}

// (GET /responses/{responseID})
func (h Handler) GetResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	res := openapi.Response{}

	return ctx.JSON(200, res)
}

// (PATCH /responses/{responseID})
func (h Handler) EditResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
	return ctx.NoContent(200)
}
