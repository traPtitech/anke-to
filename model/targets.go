package model

import (
	"net/http"

	"github.com/labstack/echo"
)

type TargetType int

const (
	Targeted = iota
	Nontargeted
	All
)

func GetTargets(c echo.Context, questionnaireID int) ([]string, error) {
	targets := []string{}
	if err := db.Select(&targets, "SELECT user_traqid FROM targets WHERE questionnaire_id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return targets, nil
}

func InsertTargets(c echo.Context, questionnaireID int, targets []string) error {
	for _, v := range targets {
		if _, err := db.Exec(
			"INSERT INTO targets (questionnaire_id, user_traqid) VALUES (?, ?)",
			questionnaireID, v); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

func DeleteTargets(c echo.Context, questionnaireID int) error {
	if _, err := db.Exec(
		"DELETE from targets WHERE questionnaire_id = ?",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
