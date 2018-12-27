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
	e.GET("/api/questionnaires", router.GetQuestionnaires)
	e.POST("/api/questionnaires", router.PostQuestionnaire)
	e.GET("/api/questionnaires/:id", router.GetQuestionnaire)
	e.PATCH("/api/questionnaires/:id", router.EditQuestionnaire)
	e.DELETE("/api/questionnaires/:id", router.DeleteQuestionnaire)
	e.GET("/api/questionnaires/:id/questions", router.GetQuestions)

	e.POST("/api/questions", router.PostQuestion)
	e.PATCH("/api/questions/:id", router.EditQuestion)
	e.DELETE("/api/questions/:id", router.DeleteQuestion)

	e.POST("/api/responses", router.PostResponse)
	e.GET("/api/responses/:id", router.GetResponse)
	e.PATCH("/api/responses/:id", router.EditResponse)
	e.DELETE("/api/responses/:id", router.DeleteResponse)

	//e.GET("/api/users", )
	e.GET("/api/users/me", router.GetUsersMe)
	e.GET("/api/users/me/responses", router.GetMyResponses)
	e.GET("/api/users/me/responses/:questionnaireID", router.GetMyResponsesByID)
	e.GET("/api/users/me/targeted", router.GetTargetedQuestionnaire)
	e.GET("/api/users/me/administrates", router.GetMyQuestionnaire)

	e.GET("/api/results/:questionnaireID", router.GetResponsesByID)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
