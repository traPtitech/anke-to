package model

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/labstack/echo"
)

type Validations struct {
	ID              int    `json:"questionID" db:"question_id"`
	RegexPattern string `json:"regex_pattern" db:"regex_pattern"`
	MinBound  string `json:"min_bound"  db:"min_bound"`
	MaxBound  string `json:"max_bound"  db:"max_bound"`
}

func GetValidations(c echo.Context, questionID int) (Validations, error) {
	validation := Validations{}
	if err := db.Get(&validation, "SELECT * FROM validations WHERE question_id = ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return Validations{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return validation, nil
}

func InsertValidations(c echo.Context, lastID int, validation Validations) error {
	if _, err := db.Exec(
		"INSERT INTO validations (question_id, regex_pattern, min_bound, max_bound) VALUES (?, ?, ?, ?)",
		lastID, validation.RegexPattern, validation.MinBound, validation.MaxBound); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func UpdateValidations(c echo.Context, questionID int, validation Validations) error {
	if _, err := db.Exec(
		`INSERT INTO validations (question_id, regex_pattern, min_bound, max_bound) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE regex_pattern = ?, min_bound = ?, max_bound = ?`,
		questionID,
		validation.RegexPattern, validation.MinBound, validation.MaxBound,
		validation.RegexPattern, validation.MinBound, validation.MaxBound); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteValidations(c echo.Context, questionID int) error {
	if _, err := db.Exec(
		"DELETE FROM validations WHERE question_id= ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func CheckNumberValid(MinBound, MaxBound string) error {
	min_bound, err := strconv.Atoi(MinBound)
	if err != nil {
		return err
	}
	max_bound, err := strconv.Atoi(MaxBound)
	if err != nil {
		return err
	}
	if min_bound > max_bound {
		return fmt.Errorf("failed: min_bound is greater than max_bound")
	}
	return nil
}
