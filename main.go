package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"

	"cloud.google.com/go/logging"
)

func main() {

	logger, err := model.GetLogger()
	if err != nil {
		panic(err)
	}

	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}

	if logger == nil {
		db.LogMode(true)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
	}))

	// Middleware
	e.Use(middleware.Recover())

	if logger != nil {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: logger.StandardLogger(logging.Info).Writer(),
		}))
		e.Logger.SetOutput(logger.StandardLogger(logging.Error).Writer())
	} else {
		e.Use(middleware.Logger())
	}

	// Static Files
	e.Static("/", "client/dist")
	e.Static("/js", "client/dist/js")
	e.Static("/img", "client/dist/img")
	e.Static("/fonts", "client/dist/fonts")
	e.Static("/css", "client/dist/css")

	e.File("/app.js", "client/dist/app.js")
	e.File("/favicon.ico", "client/dist/favicon.ico")
	e.File("*", "client/dist/index.html")

	router.SetRouting(e)

	port := os.Getenv("PORT")
	// Start server
	e.Logger.Fatal(e.Start(port))
}
