package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/traPtitech/anke-to/model"
)

// GetResults GET /results/:questionnaireID
func GetResults(c echo.Context) error {
	sort := c.QueryParam("sort")
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// アンケートの回答を確認する権限が無ければエラーを返す
	if err := checkResponseConfirmable(c, questionnaireID); err != nil {
		return err
	}

	respondentDetails, err := model.GetRespondentDetails(c, questionnaireID, sort)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, respondentDetails)
}

// アンケートの回答を確認できるか
func checkResponseConfirmable(c echo.Context, questionnaireID int) error {
	resSharedTo, err := model.GetResShared(c, questionnaireID)
	if err != nil {
		return err
	}

	switch resSharedTo {
	case "administrators":
		AmAdmin, err := model.CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	case "respondents":
		AmAdmin, err := model.CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin {
			isRespondent, err := model.CheckRespondent(c, questionnaireID)
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
