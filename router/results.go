package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
)

// Result Resultの構造体
type Result struct {
	model.IRespondent
	model.IQuestionnaire
	model.IAdministrator
}

// NewResult Resultのコンストラクタ
func NewResult(respondent model.IRespondent, questionnaire model.IQuestionnaire, administrator model.IAdministrator) *Result {
	return &Result{
		IRespondent:    respondent,
		IQuestionnaire: questionnaire,
		IAdministrator: administrator,
	}
}

// GetResults GET /results/:questionnaireID
func (r *Result) GetResults(c echo.Context) error {
	sort := c.QueryParam("sort")
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	respondentDetails, err := r.GetRespondentDetails(questionnaireID, sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, respondentDetails)
}
