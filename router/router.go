package router

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/traPtitech/anke-to/model"
)

func UserAuthenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// トークンを持たないユーザはアクセスできない
			if model.GetUserID(c) == "-" {
				return echo.NewHTTPError(http.StatusUnauthorized, "You are not logged in")
			}

			return next(c)
		}
	}
}

func SetRouting(e *echo.Echo) {

	api := e.Group("/api", UserAuthenticate())
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
			apiUsers.GET("/:traQID/targeted", GetTargettedQuestionnairesBytraQID)
		}

		apiResults := api.Group("/results")
		{
			apiResults.GET("/:questionnaireID", GetResponsesByID)
		}

	}
}
