package model

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
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
			result:     "strconv.Atoi: parsing \"1a\": invalid syntax",
		},
		{
			validation: Validations{0, "", "10", "0"},
			body:       "",
			result:     "failed: minBoundNum is greater than maxBoundNum",
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
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	validationTests := []ValidationTest{
		{
			validation: Validations{0, "", "0", "10"},
			body:       "5",
			result:     "",
		},
		{
			validation: Validations{0, "", "", "10"},
			body:       "20",
			result:     "code=400, message=Bad Request",
		},
		{
			validation: Validations{0, "", "20", "10"},
			body:       "20",
			result:     "code=500, message=Internal Server Error",
		},
	}

	for _, validationTest := range validationTests {
		validation := validationTest.validation
		body := validationTest.body
		result := validationTest.result
		if err := CheckNumberValidation(c, validation, body); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}

func TestCheckTextValidation(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	validationTests := []ValidationTest{
		{
			validation: Validations{0, "^\\d*\\.\\d*$", "", ""},
			body:       "4.55",
			result:     "",
		},
		{
			validation: Validations{0, "^\\d*\\.\\d*$", "", ""},
			body:       "a4.55",
			result:     "code=400, message=Bad Request",
		},
		{
			validation: Validations{0, "*", "", ""},
			body:       "hoge",
			result:     "code=500, message=Internal Server Error",
		},
	}
	for _, validationTest := range validationTests {
		validation := validationTest.validation
		body := validationTest.body
		result := validationTest.result
		if err := CheckTextValidation(c, validation, body); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}
