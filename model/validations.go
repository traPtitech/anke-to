package model

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo"
)

type Validations struct {
	ID           int    `json:"questionID" db:"question_id"`
	RegexPattern string `json:"regex_pattern" db:"regex_pattern"`
	MinBound     string `json:"min_bound"  db:"min_bound"`
	MaxBound     string `json:"max_bound"  db:"max_bound"`
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
	var min_bound, max_bound int
	if MinBound != "" {
		min, err := strconv.Atoi(MinBound)
		min_bound = min
		if err != nil {
			return err
		}
	}
	if MaxBound != "" {
		max, err := strconv.Atoi(MaxBound)
		max_bound = max
		if err != nil {
			return err
		}
	}

	if MinBound != "" && MaxBound != "" {
		if min_bound > max_bound {
			return fmt.Errorf("failed: min_bound is greater than max_bound")
		}
	}

	return nil
}

func CheckNumberValidation(MinBound, MaxBound, Body string) error {
	if Body == "" {
		return nil
	}
	number, err := strconv.Atoi(Body)
	if err != nil {
		return err
	}

	if MinBound != "" {
		min_bound, _ := strconv.Atoi(MinBound)
		if min_bound > number {
			return fmt.Errorf("failed: value too small")
		}
	}
	if MaxBound != "" {
		max_bound, _ := strconv.Atoi(MaxBound)
		if max_bound < number {
			return fmt.Errorf("failed: value too large")
		}
	}

	return nil
}

func CheckTextValidation(RegexPattern, Response string) error {
	r, _ := regexp.Compile(RegexPattern)
	if !r.MatchString(Response) && Response != "" {
		return fmt.Errorf("failed: %s does not match the pattern%s", Response, r)
	}

	return nil
}
