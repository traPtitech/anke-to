package model

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Validations struct {
	ID           int    `json:"questionID"    db:"question_id"   gorm:"column:question_id"`
	RegexPattern string `json:"regex_pattern" db:"regex_pattern" gorm:"column:regex_pattern"`
	MinBound     string `json:"min_bound"     db:"min_bound"     gorm:"column:min_bound"`
	MaxBound     string `json:"max_bound"     db:"max_bound"     gorm:"column:max_bound"`
}

func GetValidations(c echo.Context, questionID int) (Validations, error) {
	validation := Validations{}
	if err := gormDB.Where("question_id = ?", questionID).First(&validation).Error; gorm.IsRecordNotFoundError(err) {
		return Validations{}, nil
	} else if err != nil {
		c.Logger().Error(err)
		return Validations{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return validation, nil
}

func InsertValidations(c echo.Context, lastID int, validation Validations) error {
	validation.ID = lastID
	if err := gormDB.Create(&validation).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func UpdateValidations(c echo.Context, questionID int, validation Validations) error {
	validationBefore := Validations{}

	var err error
	if validationBefore, err = GetValidations(c, questionID); gorm.IsRecordNotFoundError(err) {
		return nil
	} else if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := gormDB.Model(&validationBefore).Update(&validation).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteValidations(c echo.Context, questionID int) error {
	if err := gormDB.Where("question_id = ?", questionID).Delete(&Validations{}).Error; err != nil {
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

func CheckNumberValidation(c echo.Context, validation Validations, Body string) error {
	if err := CheckNumberValid(validation.MinBound, validation.MaxBound); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if Body == "" {
		return nil
	}
	number, err := strconv.Atoi(Body)
	if err != nil {
		return err
	}

	if validation.MinBound != "" {
		min_bound, _ := strconv.Atoi(validation.MinBound)
		if min_bound > number {
			err := fmt.Errorf("failed: value too small")
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}
	if validation.MaxBound != "" {
		max_bound, _ := strconv.Atoi(validation.MaxBound)
		if max_bound < number {
			err := fmt.Errorf("failed: value too large")
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	return nil
}

func CheckTextValidation(c echo.Context, validation Validations, Response string) error {
	if _, err := regexp.Compile(validation.RegexPattern); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	r, _ := regexp.Compile(validation.RegexPattern)
	if !r.MatchString(Response) && Response != "" {
		err := fmt.Errorf("failed: %s does not match the pattern%s", Response, r)
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return nil
}
