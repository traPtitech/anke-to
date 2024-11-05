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
	"github.com/traPtitech/anke-to/controller"
	"github.com/traPtitech/anke-to/handler"
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

	controller.Wg.Add(1)
	go func() {
		e := echo.New()
		swagger, err := openapi.GetSwagger()
		if err != nil {
			panic(err)
		}
		e.Use(oapiMiddleware.OapiRequestValidator(swagger))
		e.Use(handler.SetUserIDMiddleware)
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		openapi.RegisterHandlers(e, handler.Handler{})
		e.Logger.Fatal(e.Start(port))
		controller.Wg.Done()
	}()

	controller.Wg.Add(1)
	go func ()  {
		controller.ReminderInit()
		controller.Wg.Done()
	}()

	controller.Wg.Add(1)
	go func() {
		controller.ReminderWorker()
		controller.Wg.Done()
	}()

	controller.Wg.Wait()
	

	// SetRouting(port)
}
