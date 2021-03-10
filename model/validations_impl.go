package model

import (
	"fmt"
	"regexp"
	"strconv"
)

// Validation ValidationRepositoryの実装
type Validation struct{}

// NewValidation Validationのコンストラクター
func NewValidation() *Validation {
	return new(Validation)
}

//Validations validationsテーブルの構造体
type Validations struct {
	QuestionID   int    `json:"questionID"    gorm:"type:int(11) PRIMARY KEY;"`
	RegexPattern string `json:"regex_pattern" gorm:"type:text;default:NULL;"`
	MinBound     string `json:"min_bound"     gorm:"type:text;default:NULL;"`
	MaxBound     string `json:"max_bound"     gorm:"type:text;default:NULL;"`
}

// InsertValidation IDを指定してvalidationsを挿入する
func (*Validation) InsertValidation(lastID int, validation Validations) error {
	validation.QuestionID = lastID
	if err := db.Create(&validation).Error; err != nil {
		return fmt.Errorf("failed to insert the validation (lastID: %d): %w", lastID, err)
	}
	return nil
}

// UpdateValidation questionIDを指定してvalidationを更新する
func (*Validation) UpdateValidation(questionID int, validation Validations) error {
	result := db.
		Model(&Validations{}).
		Where("question_id = ?", questionID).
		Update(map[string]interface{}{
			"question_id":   questionID,
			"regex_pattern": validation.RegexPattern,
			"min_bound":     validation.MinBound,
			"max_bound":     validation.MaxBound})
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to update the validation (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update a validation record: %w", ErrNoRecordUpdated)
	}
	return nil
}

// DeleteValidation questionIDを指定してvalidationを削除する
func (*Validation) DeleteValidation(questionID int) error {
	result := db.
		Where("question_id = ?", questionID).
		Delete(&Validations{})
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to delete the validation (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete a validation : %w", ErrNoRecordDeleted)
	}
	return nil
}

// GetValidations qustionIDのリストから対応するvalidationsのリストを取得する
func (*Validation) GetValidations(qustionIDs []int) ([]Validations, error) {
	validations := []Validations{}
	err := db.
		Where("question_id IN (?)", qustionIDs).
		Find(&validations).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get the validations : %w", err)
	}

	return validations, nil
}

// CheckNumberValidation BodyがMinBound,MaxBoundを満たしているか
func (v *Validation) CheckNumberValidation(validation Validations, Body string) error {
	if err := v.CheckNumberValid(validation.MinBound, validation.MaxBound); err != nil {
		return err
	}

	if Body == "" {
		return nil
	}
	number, err := strconv.ParseFloat(Body, 64)
	if err != nil {
		return ErrInvalidNumber
	}

	if validation.MinBound != "" {
		minBoundNum, _ := strconv.ParseFloat(validation.MinBound, 64)
		if minBoundNum > number {
			return fmt.Errorf("failed to meet the boundary value. the number must be greater than MinBound (number: %g, MinBound: %g): %w", number, minBoundNum, ErrNumberBoundary)
		}
	}
	if validation.MaxBound != "" {
		maxBoundNum, _ := strconv.ParseFloat(validation.MaxBound, 64)
		if maxBoundNum < number {
			return fmt.Errorf("failed to meet the boundary value. the number must be less than MaxBound (number: %g, MaxBound: %g): %w", number, maxBoundNum, ErrNumberBoundary)
		}
	}

	return nil
}

// CheckTextValidation ResponseがRegexPatternにマッチしているか
func (*Validation) CheckTextValidation(validation Validations, Response string) error {
	r, err := regexp.Compile(validation.RegexPattern)
	if err != nil {
		return fmt.Errorf("failed to compile the pattern (RegexPattern: %s): %w", r, ErrInvalidRegex)
	}
	if !r.MatchString(Response) && Response != "" {
		return fmt.Errorf("failed to match the pattern (Response: %s, RegexPattern: %s): %w", Response, r, ErrTextMatching)
	}

	return nil
}

// CheckNumberValid MinBound,MaxBoundが指定されていれば，有効な入力か確認する
func (*Validation) CheckNumberValid(MinBound, MaxBound string) error {
	var minBoundNum, maxBoundNum float64
	if MinBound != "" {
		min, err := strconv.ParseFloat(MinBound, 64)
		minBoundNum = min
		if err != nil {
			return fmt.Errorf("failed to check the boundary value. MinBound is not a numerical value: %w", ErrInvalidNumber)
		}
	}
	if MaxBound != "" {
		max, err := strconv.ParseFloat(MaxBound, 64)
		maxBoundNum = max
		if err != nil {
			return fmt.Errorf("failed to check the boundary value. MaxBound is not a numerical value: %w", ErrInvalidNumber)
		}
	}

	if MinBound != "" && MaxBound != "" {
		if minBoundNum > maxBoundNum {
			return fmt.Errorf("failed to check the boundary value. MinBound must be less than MaxBound (MinBound: %g, MaxBound: %g): %w", minBoundNum, maxBoundNum, ErrInvalidNumber)

		}
	}

	return nil
}
