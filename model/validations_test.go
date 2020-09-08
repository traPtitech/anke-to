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

func TestCheckNumberValid(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{0, "", "1a", "10"},
			body:       "",
			result:     "MinBound is not a numerical value: strconv.Atoi: parsing \"1a\": invalid syntax",
		},
		{
			validation: Validations{0, "", "10", "0"},
			body:       "",
			result:     "MinBound 10 is greater than MaxBound 0",
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
			result:     "the number 20 is larger than MaxBound 10",
		},
		{
			validation: Validations{0, "", "20", "10"},
			body:       "20",
			result:     "MinBound 20 is greater than MaxBound 10",
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
			result:     "a4.55 does not match the pattern ^\\d*\\.\\d*$",
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
