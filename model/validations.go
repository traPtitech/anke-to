package model

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//Validations validationsテーブルの構造体
type Validations struct {
	ID           int    `json:"questionID"    gorm:"column:question_id"`
	RegexPattern string `json:"regex_pattern" gorm:"column:regex_pattern"`
	MinBound     string `json:"min_bound"     gorm:"column:min_bound"`
	MaxBound     string `json:"max_bound"     gorm:"column:max_bound"`
}

// GetValidations 指定されたquestionIDのvalidationを取得する
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

// InsertValidations IDを指定してvalidationsを挿入する
func InsertValidations(c echo.Context, lastID int, validation Validations) error {
	validation.ID = lastID
	if err := gormDB.Create(&validation).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateValidations questionIDを指定してvalidationを更新する
func UpdateValidations(c echo.Context, questionID int, validation Validations) error {
	if err := gormDB.Model(&Validations{}).Update(&validation).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// DeleteValidations questionIDを指定してvalidationを削除する
func DeleteValidations(c echo.Context, questionID int) error {
	if err := gormDB.Where("question_id = ?", questionID).Delete(&Validations{}).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// CheckNumberValid MinBound,MaxBoundが指定されていれば，有効な入力か確認する
func CheckNumberValid(MinBound, MaxBound string) error {
	var minBoundNum, maxBoundNum int
	if MinBound != "" {
		min, err := strconv.Atoi(MinBound)
		minBoundNum = min
		if err != nil {
			return err
		}
	}
	if MaxBound != "" {
		max, err := strconv.Atoi(MaxBound)
		maxBoundNum = max
		if err != nil {
			return err
		}
	}

	if MinBound != "" && MaxBound != "" {
		if minBoundNum > maxBoundNum {
			return fmt.Errorf("failed: minBoundNum is greater than maxBoundNum")
		}
	}

	return nil
}

// CheckNumberValidation BodyがMinBound,MaxBoundを満たしているか
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
		minBoundNum, _ := strconv.Atoi(validation.MinBound)
		if minBoundNum > number {
			err := fmt.Errorf("failed: value too small")
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}
	if validation.MaxBound != "" {
		maxBoundNum, _ := strconv.Atoi(validation.MaxBound)
		if maxBoundNum < number {
			err := fmt.Errorf("failed: value too large")
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	return nil
}

// CheckTextValidation BodyがRegexPatternにマッチしているか
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
