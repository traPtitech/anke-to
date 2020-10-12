package model

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/jinzhu/gorm"
)

//Validations validationsテーブルの構造体
type Validations struct {
	ID           int    `json:"questionID"    gorm:"column:question_id"`
	RegexPattern string `json:"regex_pattern" gorm:"column:regex_pattern"`
	MinBound     string `json:"min_bound"     gorm:"column:min_bound"`
	MaxBound     string `json:"max_bound"     gorm:"column:max_bound"`
}

// GetValidations 指定されたquestionIDのvalidationを取得する
func GetValidations(questionID int) (Validations, error) {
	validation := Validations{}
	err := db.
		Where("question_id = ?", questionID).
		First(&validation).
		Error
	if gorm.IsRecordNotFoundError(err) {
		return Validations{}, nil
	} else if err != nil {
		return Validations{}, fmt.Errorf("failed to get the validation (questionID: %d): %w", questionID, err)
	}
	return validation, nil
}

// InsertValidations IDを指定してvalidationsを挿入する
func InsertValidations(lastID int, validation Validations) error {
	validation.ID = lastID
	if err := db.Create(&validation).Error; err != nil {
		return fmt.Errorf("failed to insert the validation (lastID: %d): %w", lastID, err)
	}
	return nil
}

// UpdateValidations questionIDを指定してvalidationを更新する
func UpdateValidations(questionID int, validation Validations) error {
	err := db.
		Model(&Validations{}).
		Where("question_id = ?", questionID).
		Update(map[string]interface{}{
			"question_id":   questionID,
			"regex_pattern": validation.RegexPattern,
			"min_bound":     validation.MinBound,
			"max_bound":     validation.MaxBound}).
		Error
	if err != nil {
		return fmt.Errorf("failed to update the validation (questionID: %d): %w", questionID, err)
	}
	return nil
}

// DeleteValidations questionIDを指定してvalidationを削除する
func DeleteValidations(questionID int) error {
	err := db.
		Where("question_id = ?", questionID).
		Delete(&Validations{}).
		Error
	if err != nil {
		return fmt.Errorf("failed to delete the validation (questionID: %d): %w", questionID, err)
	}
	return nil
}

//NumberValidError MinBound,MaxBoundの指定が有効ではない
type NumberValidError struct {
	Msg string
	Err error
}

func (e *NumberValidError) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

func (e *NumberValidError) Unwrap() error {
	return e.Err
}

// CheckNumberValid MinBound,MaxBoundが指定されていれば，有効な入力か確認する
func CheckNumberValid(MinBound, MaxBound string) error {
	var minBoundNum, maxBoundNum int
	if MinBound != "" {
		min, err := strconv.Atoi(MinBound)
		minBoundNum = min
		if err != nil {
			return &NumberValidError{"failed to check the boundary value. MinBound is not a numerical value", err}
		}
	}
	if MaxBound != "" {
		max, err := strconv.Atoi(MaxBound)
		maxBoundNum = max
		if err != nil {
			return &NumberValidError{"failed to check the boundary value. MaxBound is not a numerical value", err}
		}
	}

	if MinBound != "" && MaxBound != "" {
		if minBoundNum > maxBoundNum {
			return &NumberValidError{fmt.Sprintf("failed to check the boundary value. MinBound must be less than MaxBound (MinBound: %d, MaxBound: %d)", minBoundNum, maxBoundNum), nil}
		}
	}

	return nil
}

//NumberBoundaryError MinBound <= value <= MaxBound でない
type NumberBoundaryError struct {
	Msg string
}

func (e *NumberBoundaryError) Error() string {
	return e.Msg
}

// CheckNumberValidation BodyがMinBound,MaxBoundを満たしているか
func CheckNumberValidation(validation Validations, Body string) error {
	if err := CheckNumberValid(validation.MinBound, validation.MaxBound); err != nil {
		return err
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
			return &NumberBoundaryError{fmt.Sprintf("failed to meet the boundary value. the number must be greater than MinBound (number: %d, MinBound: %d)", number, minBoundNum)}
		}
	}
	if validation.MaxBound != "" {
		maxBoundNum, _ := strconv.Atoi(validation.MaxBound)
		if maxBoundNum < number {
			return &NumberBoundaryError{fmt.Sprintf("failed to meet the boundary value. the number must be less than MaxBound (number: %d, MaxBound: %d)", number, maxBoundNum)}
		}
	}

	return nil
}

//TextMatchError ResponceがRegexPatternにマッチしているか
type TextMatchError struct {
	Msg string
}

func (e *TextMatchError) Error() string {
	return e.Msg
}

// CheckTextValidation ResponceがRegexPatternにマッチしているか
func CheckTextValidation(validation Validations, Response string) error {
	r, err := regexp.Compile(validation.RegexPattern)
	if err != nil {
		return err
	}
	if !r.MatchString(Response) && Response != "" {
		return &TextMatchError{fmt.Sprintf("failed to match the pattern (Responce: %s, RegexPattern: %s)", Response, r)}
	}

	return nil
}
