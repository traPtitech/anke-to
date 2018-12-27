package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"git.trapti.tech/SysAd/anke-to/model"
	"git.trapti.tech/SysAd/anke-to/router"
)

func main() {

	_db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}
	model.DB = _db

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

	// Routes
	e.GET("/questionnaires", router.GetQuestionnaires)
	e.POST("/questionnaires", router.PostQuestionnaire)
	e.GET("/questionnaires/:id", router.GetQuestionnaire)
	e.PATCH("/questionnaires/:id", router.EditQuestionnaire)
	e.DELETE("/questionnaires/:id", router.DeleteQuestionnaire)
	e.GET("/questionnaires/:id/questions", router.GetQuestions)

	e.POST("/questions", router.PostQuestion)
	e.PATCH("/questions/:id", router.EditQuestion)
	e.DELETE("/questions/:id", router.DeleteQuestion)

	e.POST("/responses", router.PostResponse)
	e.GET("/responses/:id", router.GetResponse)
	e.PATCH("/responses/:id", router.EditResponse)
	e.DELETE("/responses/:id", router.DeleteResponse)

	//e.GET("/users", )
	e.GET("/users/me", router.GetUsersMe)
	e.GET("/users/me/responses", router.GetMyResponses)
	e.GET("/users/me/responses/:questionnaireID", router.GetMyResponsesByID)
	e.GET("/users/me/targeted", router.GetTargetedQuestionnaire)
	e.GET("/users/me/administrates", router.GetMyQuestionnaire)

	e.GET("/results/:questionnaireID", router.GetResponsesByID)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
