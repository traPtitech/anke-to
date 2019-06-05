package router

import (
	"github.com/labstack/echo"
)

func SetRouting(e *echo.Echo) {

	api := e.Group("/api")
	{
		apiQuestionnnaires := api.Group("/questionnaires")
		{
			apiQuestionnnaires.GET("", GetQuestionnaires)
			apiQuestionnnaires.POST("", PostQuestionnaire)
			apiQuestionnnaires.GET("/:id", GetQuestionnaire)
			apiQuestionnnaires.PATCH("/:id", EditQuestionnaire)
			apiQuestionnnaires.DELETE("/:id", DeleteQuestionnaire)
			apiQuestionnnaires.GET("/:id/questions", GetQuestions)
		}

		apiQuestions := api.Group("/questions")
		{
			apiQuestions.POST("", PostQuestion)
			apiQuestions.PATCH("/:id", EditQuestion)
			apiQuestions.DELETE("/:id", DeleteQuestion)
		}

		apiResponses := api.Group("/responses")
		{
			apiResponses.POST("", PostResponse)
			apiResponses.GET("/:id", GetResponse)
			apiResponses.PATCH("/:id", EditResponse)
			apiResponses.DELETE("/:id", DeleteResponse)
		}

		apiUsers := api.Group("/users")
		{
			/*
				TODO
				apiUsers.GET("")
			*/
			apiUsersMe := apiUsers.Group("/me")
			{
				apiUsersMe.GET("", GetUsersMe)
				apiUsersMe.GET("/responses", GetMyResponses)
				apiUsersMe.GET("/responses/:questionnaireID", GetMyResponsesByID)
				apiUsersMe.GET("/targeted", GetTargetedQuestionnaire)
				apiUsersMe.GET("/administrates", GetMyQuestionnaire)
			}
		}

		apiResults := api.Group("/results")
		{
			apiResults.GET("/:questionnaireID", GetResponsesByID)
		}

	}
}
