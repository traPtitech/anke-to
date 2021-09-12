package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
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

	// アンケートの回答を確認する権限が無ければエラーを返す
	if err := r.checkResponseConfirmable(c, questionnaireID); err != nil {
		return err
	}

	respondentDetails, err := r.GetRespondentDetails(questionnaireID, sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, respondentDetails)
}

// アンケートの回答を確認できるか
func (r *Result) checkResponseConfirmable(c echo.Context, questionnaireID int) error {
	resSharedTo, err := r.GetResShared(questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch resSharedTo {
	case "administrators":
		userID, err := getUserID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		isAdmin, err := r.CheckQuestionnaireAdmin(userID, questionnaireID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	case "respondents":
		userID, err := getUserID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		isAdmin, err := r.CheckQuestionnaireAdmin(userID, questionnaireID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			isRespondent, err := r.CheckRespondent(userID, questionnaireID)
			if err != nil {
				return err
			}
			if !isRespondent {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("only admins and respondents can see this responses"))
			}
		}
	}
	return nil
}
