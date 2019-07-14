package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"
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

	router.SetRouting(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
