package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

// SetRouting ルーティングの設定
func SetRouting(port string) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	api, err := InjectAPIServer()
	if err != nil {
		log.Panicln(err)
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


	e.Use(api.SessionMiddleware())
	echoAPI := e.Group("/api", api.SetValidatorMiddleware, api.SetUserIDMiddleware, api.TraPMemberAuthenticate)
	{
		apiQuestionnnaires := echoAPI.Group("/questionnaires")
		{
			apiQuestionnnaires.GET("", api.GetQuestionnaires, api.TrapRateLimitMiddlewareFunc())
			apiQuestionnnaires.POST("", api.PostQuestionnaire)
			apiQuestionnnaires.GET("/:questionnaireID", api.GetQuestionnaire)
			apiQuestionnnaires.PATCH("/:questionnaireID", api.EditQuestionnaire, api.QuestionnaireAdministratorAuthenticate)
			apiQuestionnnaires.DELETE("/:questionnaireID", api.DeleteQuestionnaire, api.QuestionnaireAdministratorAuthenticate)
			apiQuestionnnaires.GET("/:questionnaireID/questions", api.GetQuestions)
			apiQuestionnnaires.POST("/:questionnaireID/questions", api.PostQuestionByQuestionnaireID)
		}

		apiQuestions := echoAPI.Group("/questions")
		{
			apiQuestions.PATCH("/:questionID", api.EditQuestion, api.QuestionAdministratorAuthenticate)
			apiQuestions.DELETE("/:questionID", api.DeleteQuestion, api.QuestionAdministratorAuthenticate)
		}

		apiResponses := echoAPI.Group("/responses")
		{
			apiResponses.POST("", api.PostResponse)
			apiResponses.GET("/:responseID", api.GetResponse, api.ResponseReadAuthenticate)
			apiResponses.PATCH("/:responseID", api.EditResponse, api.RespondentAuthenticate)
			apiResponses.DELETE("/:responseID", api.DeleteResponse, api.RespondentAuthenticate)
		}

		apiUsers := echoAPI.Group("/users")
		{

			apiUsersMe := apiUsers.Group("/me")
			{
				apiUsersMe.GET("", api.GetUsersMe)
				apiUsersMe.GET("/responses", api.GetMyResponses)
				apiUsersMe.GET("/responses/:questionnaireID", api.GetMyResponsesByID)
				apiUsersMe.GET("/targeted", api.GetTargetedQuestionnaire)
				apiUsersMe.GET("/administrates", api.GetMyQuestionnaire)
			}
			apiUsers.GET("/:traQID/targeted", api.GetTargettedQuestionnairesBytraQID)
		}

		apiResults := echoAPI.Group("/results")
		{
			apiResults.GET("/:questionnaireID", api.GetResults, api.ResultAuthenticate)
		}

		apiOauth := echoAPI.Group("/oauth")
		{
			apiOauth.GET("/callback", api.Callback)
			apiOauth.GET("/generate/code", api.GetCode)
		}
	}

	e.Logger.Fatal(e.Start(port))
}
