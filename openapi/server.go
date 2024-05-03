// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package openapi

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /questionnaires)
	GetQuestionnaires(ctx echo.Context, params GetQuestionnairesParams) error

	// (POST /questionnaires)
	PostQuestionnaire(ctx echo.Context) error

	// (DELETE /questionnaires/{questionnaireID})
	DeleteQuestionnaire(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (GET /questionnaires/{questionnaireID})
	GetQuestionnaire(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (PATCH /questionnaires/{questionnaireID})
	EditQuestionnaire(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (GET /questionnaires/{questionnaireID}/myRemindStatus)
	GetQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (PATCH /questionnaires/{questionnaireID}/myRemindStatus)
	EditQuestionnaireMyRemindStatus(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (GET /questionnaires/{questionnaireID}/responses)
	GetQuestionnaireResponses(ctx echo.Context, questionnaireID QuestionnaireIDInPath, params GetQuestionnaireResponsesParams) error

	// (POST /questionnaires/{questionnaireID}/responses)
	PostQuestionnaireResponse(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (GET /questionnaires/{questionnaireID}/result)
	GetQuestionnaireResult(ctx echo.Context, questionnaireID QuestionnaireIDInPath) error

	// (GET /responses/myResponses)
	GetMyResponses(ctx echo.Context, params GetMyResponsesParams) error

	// (DELETE /responses/{responseID})
	DeleteResponse(ctx echo.Context, responseID ResponseIDInPath) error

	// (GET /responses/{responseID})
	GetResponse(ctx echo.Context, responseID ResponseIDInPath) error

	// (PATCH /responses/{responseID})
	EditResponse(ctx echo.Context, responseID ResponseIDInPath) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetQuestionnaires converts echo context to params.
func (w *ServerInterfaceWrapper) GetQuestionnaires(ctx echo.Context) error {
	var err error

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetQuestionnairesParams
	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", ctx.QueryParams(), &params.Sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sort: %s", err))
	}

	// ------------- Optional query parameter "search" -------------

	err = runtime.BindQueryParameter("form", true, false, "search", ctx.QueryParams(), &params.Search)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter search: %s", err))
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "onlyTargetingMe" -------------

	err = runtime.BindQueryParameter("form", true, false, "onlyTargetingMe", ctx.QueryParams(), &params.OnlyTargetingMe)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter onlyTargetingMe: %s", err))
	}

	// ------------- Optional query parameter "onlyAdministratedByMe" -------------

	err = runtime.BindQueryParameter("form", true, false, "onlyAdministratedByMe", ctx.QueryParams(), &params.OnlyAdministratedByMe)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter onlyAdministratedByMe: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetQuestionnaires(ctx, params)
	return err
}

// PostQuestionnaire converts echo context to params.
func (w *ServerInterfaceWrapper) PostQuestionnaire(ctx echo.Context) error {
	var err error

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostQuestionnaire(ctx)
	return err
}

// DeleteQuestionnaire converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteQuestionnaire(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteQuestionnaire(ctx, questionnaireID)
	return err
}

// GetQuestionnaire converts echo context to params.
func (w *ServerInterfaceWrapper) GetQuestionnaire(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetQuestionnaire(ctx, questionnaireID)
	return err
}

// EditQuestionnaire converts echo context to params.
func (w *ServerInterfaceWrapper) EditQuestionnaire(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.EditQuestionnaire(ctx, questionnaireID)
	return err
}

// GetQuestionnaireMyRemindStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetQuestionnaireMyRemindStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetQuestionnaireMyRemindStatus(ctx, questionnaireID)
	return err
}

// EditQuestionnaireMyRemindStatus converts echo context to params.
func (w *ServerInterfaceWrapper) EditQuestionnaireMyRemindStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.EditQuestionnaireMyRemindStatus(ctx, questionnaireID)
	return err
}

// GetQuestionnaireResponses converts echo context to params.
func (w *ServerInterfaceWrapper) GetQuestionnaireResponses(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetQuestionnaireResponsesParams
	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", ctx.QueryParams(), &params.Sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sort: %s", err))
	}

	// ------------- Optional query parameter "onlyMyResponse" -------------

	err = runtime.BindQueryParameter("form", true, false, "onlyMyResponse", ctx.QueryParams(), &params.OnlyMyResponse)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter onlyMyResponse: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetQuestionnaireResponses(ctx, questionnaireID, params)
	return err
}

// PostQuestionnaireResponse converts echo context to params.
func (w *ServerInterfaceWrapper) PostQuestionnaireResponse(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostQuestionnaireResponse(ctx, questionnaireID)
	return err
}

// GetQuestionnaireResult converts echo context to params.
func (w *ServerInterfaceWrapper) GetQuestionnaireResult(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "questionnaireID" -------------
	var questionnaireID QuestionnaireIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "questionnaireID", runtime.ParamLocationPath, ctx.Param("questionnaireID"), &questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter questionnaireID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetQuestionnaireResult(ctx, questionnaireID)
	return err
}

// GetMyResponses converts echo context to params.
func (w *ServerInterfaceWrapper) GetMyResponses(ctx echo.Context) error {
	var err error

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMyResponsesParams
	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", ctx.QueryParams(), &params.Sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sort: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMyResponses(ctx, params)
	return err
}

// DeleteResponse converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteResponse(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "responseID" -------------
	var responseID ResponseIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "responseID", runtime.ParamLocationPath, ctx.Param("responseID"), &responseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter responseID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteResponse(ctx, responseID)
	return err
}

// GetResponse converts echo context to params.
func (w *ServerInterfaceWrapper) GetResponse(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "responseID" -------------
	var responseID ResponseIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "responseID", runtime.ParamLocationPath, ctx.Param("responseID"), &responseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter responseID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetResponse(ctx, responseID)
	return err
}

// EditResponse converts echo context to params.
func (w *ServerInterfaceWrapper) EditResponse(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "responseID" -------------
	var responseID ResponseIDInPath

	err = runtime.BindStyledParameterWithLocation("simple", false, "responseID", runtime.ParamLocationPath, ctx.Param("responseID"), &responseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter responseID: %s", err))
	}

	ctx.Set(ApplicationScopes, []string{"read", "write"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.EditResponse(ctx, responseID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/questionnaires", wrapper.GetQuestionnaires)
	router.POST(baseURL+"/questionnaires", wrapper.PostQuestionnaire)
	router.DELETE(baseURL+"/questionnaires/:questionnaireID", wrapper.DeleteQuestionnaire)
	router.GET(baseURL+"/questionnaires/:questionnaireID", wrapper.GetQuestionnaire)
	router.PATCH(baseURL+"/questionnaires/:questionnaireID", wrapper.EditQuestionnaire)
	router.GET(baseURL+"/questionnaires/:questionnaireID/myRemindStatus", wrapper.GetQuestionnaireMyRemindStatus)
	router.PATCH(baseURL+"/questionnaires/:questionnaireID/myRemindStatus", wrapper.EditQuestionnaireMyRemindStatus)
	router.GET(baseURL+"/questionnaires/:questionnaireID/responses", wrapper.GetQuestionnaireResponses)
	router.POST(baseURL+"/questionnaires/:questionnaireID/responses", wrapper.PostQuestionnaireResponse)
	router.GET(baseURL+"/questionnaires/:questionnaireID/result", wrapper.GetQuestionnaireResult)
	router.GET(baseURL+"/responses/myResponses", wrapper.GetMyResponses)
	router.DELETE(baseURL+"/responses/:responseID", wrapper.DeleteResponse)
	router.GET(baseURL+"/responses/:responseID", wrapper.GetResponse)
	router.PATCH(baseURL+"/responses/:responseID", wrapper.EditResponse)

}
