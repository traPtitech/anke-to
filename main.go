package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"git.trapti.tech/SysAd/anke-to/model"
	"git.trapti.tech/SysAd/anke-to/router"
)

func main() {

	err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
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

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
