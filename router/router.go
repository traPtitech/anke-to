package router

import (
	"github.com/labstack/echo"
)

// SetRouting ルーティングの設定
func SetRouting(e *echo.Echo) {
	api := e.Group("/api", UserAuthenticate)
	{
		apiQuestionnnaires := api.Group("/questionnaires")
		{
			apiQuestionnnaires.GET("", GetQuestionnaires)
			apiQuestionnnaires.POST("", PostQuestionnaire)
			apiQuestionnnaires.GET("/:questionnaireID", GetQuestionnaire)
			apiQuestionnnaires.PATCH("/:questionnaireID", EditQuestionnaire, QuestionnaireAdministratorAuthenticate)
			apiQuestionnnaires.DELETE("/:questionnaireID", DeleteQuestionnaire, QuestionnaireAdministratorAuthenticate)
			apiQuestionnnaires.GET("/:questionnaireID/questions", GetQuestions)
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
			apiResponses.GET("/:responseID", GetResponse)
			apiResponses.PATCH("/:responseID", EditResponse, RespondentAuthenticate)
			apiResponses.DELETE("/:responseID", DeleteResponse, RespondentAuthenticate)
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
			apiUsers.GET("/:traQID/targeted", GetTargettedQuestionnairesBytraQID)
		}

		apiResults := api.Group("/results")
		{
			apiResults.GET("/:questionnaireID", GetResults)
		}
	}
}
