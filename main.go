package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapiMiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"

	"github.com/traPtitech/anke-to/tuning"
)

func main() {
	env, ok := os.LookupEnv("ANKE-TO_ENV")
	if !ok {
		env = "production"
	}
	logOn := env == "pprof" || env == "dev"

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			tuning.Inititial()
			return
		case "bench":
			tuning.Bench()
			return
		}
	}

	err := model.EstablishConnection(!logOn)
	if err != nil {
		panic(err)
	}

	_, err = model.Migrate()
	if err != nil {
		panic(err)
	}

	if env == "pprof" {
		runtime.SetBlockProfileRate(1)
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("no PORT")
	}

	e := echo.New()
	api := InjectAPIServer()

	api.Reminder.Wg.Add(1)
	go func() {
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		swagger, err := openapi.GetSwagger()
		if err != nil {
			panic(err)
		}
		e.Use(oapiMiddleware.OapiRequestValidator(swagger))

		e.Use(api.Middleware.SetUserIDMiddleware)

		mws := NewMiddlewareSwitcher()
		mws.AddGroupConfig("", api.Middleware.TraPMemberAuthenticate)

		mws.AddRouteConfig("/questionnaires", http.MethodGet, api.Middleware.TrapRateLimitMiddlewareFunc())
		mws.AddRouteConfig("/questionnaires/:questionnaireID", http.MethodGet, api.Middleware.QuestionnaireReadAuthenticate)
		mws.AddRouteConfig("/questionnaires/:questionnaireID", http.MethodPatch, api.Middleware.QuestionnaireAdministratorAuthenticate)
		mws.AddRouteConfig("/questionnaires/:questionnaireID", http.MethodDelete, api.Middleware.QuestionnaireAdministratorAuthenticate)
		mws.AddRouteConfig("/questionnaires/:questionnaireID/responses", http.MethodPost, api.Middleware.QuestionnaireReadAuthenticate)
		mws.AddRouteConfig("/questionnaires/:questionnaireID/responses", http.MethodGet, api.Middleware.ResultAuthenticate)

		mws.AddRouteConfig("/responses/:responseID", http.MethodGet, api.Middleware.ResponseReadAuthenticate)
		mws.AddRouteConfig("/responses/:responseID", http.MethodPatch, api.Middleware.RespondentAuthenticate)
		mws.AddRouteConfig("/responses/:responseID", http.MethodDelete, api.Middleware.RespondentAuthenticate)

		e.Use(mws.ApplyMiddlewares)

		openapi.RegisterHandlers(e, api)

		e.Logger.Fatal(e.Start(port))

		api.Reminder.Wg.Done()
	}()

	api.Reminder.Wg.Add(1)
	go func() {
		api.Reminder.ReminderWorker()
		api.Reminder.Wg.Done()
	}()

	api.Reminder.Wg.Wait()

	// SetRouting(port)
}
