package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValidationTest struct {
	validation Validations
	body       string
	result     string
}

func TestCheckNumberValidation(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{0, "", "0", "10"},
			body:       "5",
			result:     "",
		},
		{
			validation: Validations{0, "", "", "10"},
			body:       "20",
			result:     "failed to meet the boundary value. the number must be less than MaxBound (number: 20, MaxBound: 10)",
		},
		{
			validation: Validations{0, "", "20", "10"},
			body:       "20",
			result:     "failed to check the boundary value. MinBound must be less than MaxBound (MinBound: 20, MaxBound: 10)",
		},
	}

	for _, validationTest := range validationTests {
		validation := validationTest.validation
		body := validationTest.body
		result := validationTest.result
		if err := CheckNumberValidation(validation, body); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}

func TestCheckTextValidation(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{0, "^\\d*\\.\\d*$", "", ""},
			body:       "4.55",
			result:     "",
		},
		{
			validation: Validations{0, "^\\d*\\.\\d*$", "", ""},
			body:       "a4.55",
			result:     "failed to match the pattern (Responce: a4.55, RegexPattern: ^\\d*\\.\\d*$)",
		},
		{
			validation: Validations{0, "*", "", ""},
			body:       "hoge",
			result:     "error parsing regexp: missing argument to repetition operator: `*`",
		},
	}
	for _, validationTest := range validationTests {
		validation := validationTest.validation
		body := validationTest.body
		result := validationTest.result
		if err := CheckTextValidation(validation, body); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}

func TestCheckNumberValid(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{0, "", "1a", "10"},
			body:       "",
			result:     "failed to check the boundary value. MinBound is not a numerical value: strconv.Atoi: parsing \"1a\": invalid syntax",
		},
		{
			validation: Validations{0, "", "10", "0"},
			body:       "",
			result:     "failed to check the boundary value. MinBound must be less than MaxBound (MinBound: 10, MaxBound: 0)",
		},
	}

	for _, validationTest := range validationTests {
		validation := validationTest.validation
		result := validationTest.result
		if err := CheckNumberValid(validation.MinBound, validation.MaxBound); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}
