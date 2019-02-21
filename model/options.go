package model

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetOptions(c echo.Context, questionID int) ([]string, error) {
	options := []string{}
	if err := db.Select(
		&options, "SELECT body FROM options WHERE question_id = ? ORDER BY option_num",
		questionID); err != nil {
		c.Logger().Error(err)
		return []string{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return options, nil
}

func InsertOption(c echo.Context, lastID int, num int, body string) error {
	if _, err := db.Exec(
		"INSERT INTO options (question_id, option_num, body) VALUES (?, ?, ?)",
		lastID, num, body); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func UpdateOptions(c echo.Context, options []string, questionID int) error {
	for i, v := range options {
		if _, err := db.Exec(
			`INSERT INTO options (question_id, option_num, body) VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE option_num = ?, body = ?`,
			questionID, i+1, v, i+1, v); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	if _, err := db.Exec(
		"DELETE FROM options WHERE question_id= ? AND option_num > ?",
		questionID, len(options)); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteOptions(c echo.Context, questionID int) error {
	if _, err := db.Exec(
		"DELETE FROM options WHERE question_id= ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
