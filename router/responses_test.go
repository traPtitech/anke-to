package router

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/model/mock_model"
	"gopkg.in/guregu/null.v3"
)

type responseBody struct {
	QuestionID     int         `json:"questionID" validate:"min=0"`
	QuestionType   string      `json:"question_type" validate:"required,oneof=Text TextArea Number MultipleChoice Checkbox LinearScale"`
	Body           null.String `json:"response"  validate:"required"`
	OptionResponse []string    `json:"option_response"  validate:"required_if=QuestionType Checkbox,required_if=QuestionType MultipleChoice,dive,max=50"`
}

func TestPostResponseValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		request     *Responses
		isErr       bool
	}{
		{
			description: "一般的なリクエストなのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
		},
		{
			description: "IDが0でもエラーなし",
			request: &Responses{
				ID:          0,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
		},
		{
			description: "BodyのQuestionIDが0でもエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     0,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
		},
		{
			description: "ResponsesのIDが負なのでエラー",
			request: &Responses{
				ID:          -1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
			isErr: true,
		},
		{
			description: "Temporarilyがtrueでもエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: true,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
			isErr: false,
		},
		{
			description: "Bodyがnilなのでエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body:        nil,
			},
			isErr: true,
		},
		{
			description: "BodyのQuestionIDが負なのでエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     -1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
			isErr: true,
		},
		{
			description: "TextタイプでoptionResponseが50文字以上でエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: []string{"012345678901234567890123456789012345678901234567890"},
					},
				},
			},
			isErr: true,
		},
		{
			description: "TextタイプでoptionResponseが50文字ピッタリはエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Text",
						Body:           null.String{},
						OptionResponse: []string{"01234567890123456789012345678901234567890123456789"},
					},
				},
			},
		},
		{
			description: "一般的なTextAreaタイプの回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "TextArea",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
		},
		{
			description: "TextAreaタイプでoptionResponseが50文字以上でもエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "TextArea",
						Body:           null.String{},
						OptionResponse: []string{"012345678901234567890123456789012345678901234567890"},
					},
				},
			},
			isErr: true,
		},
		{
			description: "TextAreaタイプでoptionResponseが50文字ピッタリはエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "TextArea",
						Body:           null.String{},
						OptionResponse: []string{"01234567890123456789012345678901234567890123456789"},
					},
				},
			},
		},
		{
			description: "一般的なNumberタイプの回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Number",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
		},
		{
			description: "NumberタイプでoptionResponseが50文字以上でもエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Number",
						Body:           null.String{},
						OptionResponse: []string{"012345678901234567890123456789012345678901234567890"},
					},
				},
			},
			isErr: true,
		},
		{
			description: "NumberタイプでoptionResponseが50文字ピッタリでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Number",
						Body:           null.String{},
						OptionResponse: []string{"01234567890123456789012345678901234567890123456789"},
					},
				},
			},
		},
		{
			description: "Checkboxタイプで一般的な回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Checkbox",
						Body:           null.String{},
						OptionResponse: []string{"a", "b"},
					},
				},
			},
		},
		{
			description: "CheckboxタイプでOptionResponseがnilな回答なのでエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Checkbox",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
			isErr: true,
		},
		{
			description: "CheckboxタイプでOptionResponseが50文字以上な回答なのでエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Checkbox",
						Body:           null.String{},
						OptionResponse: []string{"012345678901234567890123456789012345678901234567890"},
					},
				},
			},
			isErr: true,
		},
		{
			description: "CheckboxタイプでOptionResponseが50文字ピッタリな回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "Checkbox",
						Body:           null.String{},
						OptionResponse: []string{"01234567890123456789012345678901234567890123456789"},
					},
				},
			},
		},
		{
			description: "MultipleChoiceタイプで一般的な回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "MultipleChoice",
						Body:           null.String{},
						OptionResponse: []string{"a", "b"},
					},
				},
			},
		},
		{
			description: "MultipleChoiceタイプでOptionResponseがnilな回答なのでエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "MultipleChoice",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
			isErr: true,
		},
		{
			description: "MultipleChoiceタイプでOptionResponseが50文字以上な回答なのでエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "MultipleChoice",
						Body:           null.String{},
						OptionResponse: []string{"012345678901234567890123456789012345678901234567890"},
					},
				},
			},
			isErr: true,
		},
		{
			description: "MultipleChoiceタイプでOptionResponseが50文字ピッタリな回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "MultipleChoice",
						Body:           null.String{},
						OptionResponse: []string{"01234567890123456789012345678901234567890123456789"},
					},
				},
			},
		},
		{
			description: "一般的なLinearScaleタイプの回答なのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "LinearScale",
						Body:           null.String{},
						OptionResponse: nil,
					},
				},
			},
		},
		{
			description: "LinearScaleタイプでoptionResponseが50文字以上でもエラー",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "LinearScale",
						Body:           null.String{},
						OptionResponse: []string{"012345678901234567890123456789012345678901234567890"},
					},
				},
			},
			isErr: true,
		},
		{
			description: "LinearScaleタイプでoptionResponseが50文字ピッタリなのでエラーなし",
			request: &Responses{
				ID:          1,
				Temporarily: false,
				Body: []model.ResponseBody{
					{
						QuestionID:     1,
						QuestionType:   "LinearScale",
						Body:           null.String{},
						OptionResponse: []string{"01234567890123456789012345678901234567890123456789"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		validate := validator.New()
		t.Run(test.description, func(t *testing.T) {
			err := validate.Struct(test.request)

			if test.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostResponse(t *testing.T) {
	type responseRequestBody struct {
		QuestionnaireID int            `json:"questionnaireID" validate:"min=0"`
		SubmittedAt     null.Time      `json:"submitted_at" validate:"-"`
		Body            []responseBody `json:"body" validate:"required"`
	}
	type responseResponseBody struct {
		Body            []responseBody `json:"body" validate:"required"`
		QuestionnaireID int            `json:"questionnaireID" validate:"min=0"`
		ResponseID      int            `json:"responseID" validate:"min=0"`
		SubmittedAt     null.Time      `json:"submitted_at" validate:"-"`
	}

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	questionnaireIDSuccess := 1
	questionIDSuccess := 1
	responseIDSuccess := 1

	questionnaireIDFailure := 0
	questionIDFailure := 0
	responseIDFailure := 0

	questionnaireIDLimit := 2

	validation :=
		model.Validations{
			QuestionID:   questionIDSuccess,
			RegexPattern: "^\\d*\\.\\d*$",
			MinBound:     "0",
			MaxBound:     "10",
		}
	scalelabel :=
		model.ScaleLabels{
			QuestionID:      questionIDSuccess,
			ScaleLabelRight: "そう思わない",
			ScaleLabelLeft:  "そう思う",
			ScaleMin:        1,
			ScaleMax:        5,
		}
	// questionnaireIDNotFound := -1
	// questionIDNotFound := -1
	// responseIDNotFound := -1

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)

	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)

	r := NewResponse(
		mockQuestionnaire,
		mockValidation,
		mockScaleLabel,
		mockRespondent,
		mockResponse,
	)
	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
	)
	// Questionnaire
	// GetQuestionnaireLimit
	// success
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(gomock.Any(), questionnaireIDSuccess).
		Return(null.TimeFrom(nowTime.Add(time.Minute)), nil).AnyTimes()
	// failure
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(gomock.Any(), questionnaireIDFailure).
		Return(null.NewTime(time.Time{}, false), model.ErrRecordNotFound).AnyTimes()
	// limit
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(gomock.Any(), questionnaireIDLimit).
		Return(null.TimeFrom(nowTime.Add(-time.Minute)), nil).AnyTimes()

	// Validation
	// GetValidations
	// success
	mockValidation.EXPECT().
		GetValidations(gomock.Any(), []int{questionIDSuccess}).
		Return([]model.Validations{validation}, nil).AnyTimes()
	// failure
	mockValidation.EXPECT().
		GetValidations(gomock.Any(), []int{questionIDFailure}).
		Return([]model.Validations{}, nil).AnyTimes()
	// nothing
	mockValidation.EXPECT().
		GetValidations(gomock.Any(), []int{}).
		Return([]model.Validations{}, nil).AnyTimes()
	// CheckNumberValidation
	// success
	mockValidation.EXPECT().
		CheckNumberValidation(validation, "success case").
		Return(nil).AnyTimes()
	// ErrInvalidNumber
	mockValidation.EXPECT().
		CheckNumberValidation(validation, "ErrInvalidNumber").
		Return(model.ErrInvalidNumber).AnyTimes()
	// BadRequest
	mockValidation.EXPECT().
		CheckNumberValidation(validation, "BadRequest").
		Return(errMock).AnyTimes()

	// CheckTextValidation
	// success
	mockValidation.EXPECT().
		CheckTextValidation(validation, "success case").
		Return(nil).AnyTimes()
	// ErrTextMatching
	mockValidation.EXPECT().
		CheckTextValidation(validation, "ErrTextMatching").
		Return(model.ErrTextMatching).AnyTimes()
	// InternalServerError
	mockValidation.EXPECT().
		CheckTextValidation(validation, "InternalServerError").
		Return(errMock).AnyTimes()

	// ScaleLabel
	// GetScaleLabels
	// success
	mockScaleLabel.EXPECT().
		GetScaleLabels(gomock.Any(), []int{questionIDSuccess}).
		Return([]model.ScaleLabels{scalelabel}, nil).AnyTimes()
	// failure
	mockScaleLabel.EXPECT().
		GetScaleLabels(gomock.Any(), []int{questionIDFailure}).
		Return([]model.ScaleLabels{}, nil).AnyTimes()
	// nothing
	mockScaleLabel.EXPECT().
		GetScaleLabels(gomock.Any(), []int{}).
		Return([]model.ScaleLabels{}, nil).AnyTimes()

	// CheckScaleLabel
	// success
	mockScaleLabel.EXPECT().
		CheckScaleLabel(scalelabel, "success case").
		Return(nil).AnyTimes()
	// BadRequest
	mockScaleLabel.EXPECT().
		CheckScaleLabel(scalelabel, "BadRequest").
		Return(errMock).AnyTimes()

	// Respondent
	// InsertRespondent
	// success
	mockRespondent.EXPECT().
		InsertRespondent(gomock.Any(), string(userOne), questionnaireIDSuccess, gomock.Any()).
		Return(responseIDSuccess, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		InsertRespondent(gomock.Any(), string(userOne), questionnaireIDFailure, gomock.Any()).
		Return(responseIDFailure, nil).AnyTimes()

	// Response
	// InsertResponses
	// success
	mockResponse.EXPECT().
		InsertResponses(gomock.Any(), responseIDSuccess, gomock.Any()).
		Return(nil).AnyTimes()
	// failure
	mockResponse.EXPECT().
		InsertResponses(gomock.Any(), responseIDFailure, gomock.Any()).
		Return(errMock).AnyTimes()

	// responseID, err := mockRespondent.
	// 	InsertRespondent(string(userOne), 1, null.NewTime(nowTime, true))
	// assertion.Equal(1, responseID)
	// assertion.NoError(err)

	type request struct {
		user             users
		isBadRequestBody bool
		requestBody      responseRequestBody
	}
	type expect struct {
		isErr      bool
		code       int
		responseID int
	}
	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr:      false,
				code:       http.StatusCreated,
				responseID: responseIDSuccess,
			},
		},
		{
			description: "null submittedat",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr:      false,
				code:       http.StatusCreated,
				responseID: responseIDSuccess,
			},
		},
		{
			description: "empty body",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body:            []responseBody{},
				},
			},
			expect: expect{
				isErr:      false,
				code:       http.StatusCreated,
				responseID: responseIDSuccess,
			},
		},
		{
			description: "questionnaire does not exist",
			request: request{
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDFailure,
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "bad request body",
			request: request{
				isBadRequestBody: true,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "limit exceeded",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDLimit,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body:            []responseBody{},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusMethodNotAllowed,
			},
		},
		{
			description: "valid number",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Number",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr:      false,
				code:       http.StatusCreated,
				responseID: responseIDSuccess,
			},
		},
		{
			description: "invalid number",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Number",
							Body:           null.StringFrom("ErrInvalidNumber"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "BadRequest number",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Number",
							Body:           null.StringFrom("BadRequest"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "valid text",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr:      false,
				code:       http.StatusCreated,
				responseID: responseIDSuccess,
			},
		},
		{
			description: "text does not match",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("ErrTextMatching"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "invalid text",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("InternalServerError"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "valid LinearScale",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "LinearScale",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr:      false,
				code:       http.StatusCreated,
				responseID: responseIDSuccess,
			},
		},
		{
			description: "invalid LinearScale",
			request: request{
				user: userOne,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "LinearScale",
							Body:           null.StringFrom("BadRequest"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
	}

	e := echo.New()
	e.POST("/api/responses", r.PostResponse, m.SetUserIDMiddleware, m.SetValidatorMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {
		requestByte, jsonErr := json.Marshal(testCase.request.requestBody)
		require.NoError(t, jsonErr)
		requestStr := string(requestByte) + "\n"

		if testCase.request.isBadRequestBody {
			requestStr = "badRequestBody"
		}
		rec := createRecorder(e, testCase.request.user, methodPost, makePath("/responses"), typeJSON, requestStr)

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		response := responseResponseBody{
			ResponseID:      testCase.expect.responseID,
			QuestionnaireID: testCase.request.requestBody.QuestionnaireID,
			SubmittedAt:     testCase.request.requestBody.SubmittedAt,
			Body:            testCase.request.requestBody.Body,
		}

		responseByte, jsonErr := json.Marshal(response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

func TestGetResponse(t *testing.T) {

	type responseResponseBody struct {
		QuestionnaireID int            `json:"questionnaireID"`
		SubmittedAt     null.Time      `json:"submitted_at"`
		ModifiedAt      null.Time      `json:"modified_at"`
		Body            []responseBody `json:"body"`
	}

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	responseIDSuccess := 1
	responseIDFailure := 0
	responseIDNotFound := -1

	questionnaireIDSuccess := 1
	questionIDSuccess := 1
	respondentDetail := model.RespondentDetail{
		QuestionnaireID: questionnaireIDSuccess,
		SubmittedAt:     null.TimeFrom(nowTime),
		ModifiedAt:      nowTime,
		Responses: []model.ResponseBody{
			{
				QuestionID:     questionIDSuccess,
				QuestionType:   "Text",
				Body:           null.StringFrom("回答"),
				OptionResponse: []string{},
			},
		},
	}

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)

	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)

	r := NewResponse(
		mockQuestionnaire,
		mockValidation,
		mockScaleLabel,
		mockRespondent,
		mockResponse,
	)
	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
	)

	// Respondent
	// InsertRespondent
	// success
	mockRespondent.EXPECT().
		GetRespondentDetail(gomock.Any(), responseIDSuccess).
		Return(respondentDetail, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		GetRespondentDetail(gomock.Any(), responseIDFailure).
		Return(model.RespondentDetail{}, errMock).AnyTimes()
	// NotFound
	mockRespondent.EXPECT().
		GetRespondentDetail(gomock.Any(), responseIDNotFound).
		Return(model.RespondentDetail{}, model.ErrRecordNotFound).AnyTimes()

	type request struct {
		user       users
		responseID int
	}
	type expect struct {
		isErr    bool
		code     int
		response responseResponseBody
	}

	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				responseID: responseIDSuccess,
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
				response: responseResponseBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.TimeFrom(nowTime),
					ModifiedAt:      null.TimeFrom(nowTime),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("回答"),
							OptionResponse: []string{},
						},
					},
				},
			},
		},
		{
			description: "failure",
			request: request{
				responseID: responseIDFailure,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "NotFound",
			request: request{
				responseID: responseIDNotFound,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusNotFound,
			},
		},
	}

	e := echo.New()
	e.GET("/api/responses/:responseID", r.GetResponse, m.SetUserIDMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {

		rec := createRecorder(e, testCase.request.user, methodGet, fmt.Sprint(rootPath, "/responses/", testCase.request.responseID), typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		responseByte, jsonErr := json.Marshal(testCase.expect.response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

func TestEditResponse(t *testing.T) {
	type responseRequestBody struct {
		QuestionnaireID int            `json:"questionnaireID"`
		SubmittedAt     null.Time      `json:"submitted_at"`
		Body            []responseBody `json:"body"`
	}
	type responseResponseBody struct {
		Body            []responseBody `json:"body"`
		QuestionnaireID int            `json:"questionnaireID"`
		ResponseID      int            `json:"responseID"`
		SubmittedAt     null.Time      `json:"submitted_at"`
	}

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	questionnaireIDSuccess := 1
	questionIDSuccess := 1
	responseIDSuccess := 1

	questionnaireIDFailure := 0
	questionIDFailure := 0
	responseIDFailure := 0

	questionnaireIDLimit := 2

	validation :=
		model.Validations{
			QuestionID:   questionIDSuccess,
			RegexPattern: "^\\d*\\.\\d*$",
			MinBound:     "0",
			MaxBound:     "10",
		}
	scalelabel :=
		model.ScaleLabels{
			QuestionID:      questionIDSuccess,
			ScaleLabelRight: "そう思わない",
			ScaleLabelLeft:  "そう思う",
			ScaleMin:        1,
			ScaleMax:        5,
		}
	// questionnaireIDNotFound := -1
	// questionIDNotFound := -1
	// responseIDNotFound := -1

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)

	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)

	r := NewResponse(
		mockQuestionnaire,
		mockValidation,
		mockScaleLabel,
		mockRespondent,
		mockResponse,
	)
	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
	)
	// Questionnaire
	// GetQuestionnaireLimit
	// success
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(gomock.Any(), questionnaireIDSuccess).
		Return(null.TimeFrom(nowTime.Add(time.Minute)), nil).AnyTimes()
	// failure
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(gomock.Any(), questionnaireIDFailure).
		Return(null.NewTime(time.Time{}, false), errMock).AnyTimes()
	// limit
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(gomock.Any(), questionnaireIDLimit).
		Return(null.TimeFrom(nowTime.Add(-time.Minute)), nil).AnyTimes()

	// Validation
	// GetValidations
	// success
	mockValidation.EXPECT().
		GetValidations(gomock.Any(), []int{questionIDSuccess}).
		Return([]model.Validations{validation}, nil).AnyTimes()
	// failure
	mockValidation.EXPECT().
		GetValidations(gomock.Any(), []int{questionIDFailure}).
		Return([]model.Validations{}, nil).AnyTimes()
	// nothing
	mockValidation.EXPECT().
		GetValidations(gomock.Any(), []int{}).
		Return([]model.Validations{}, nil).AnyTimes()
	// CheckNumberValidation
	// success
	mockValidation.EXPECT().
		CheckNumberValidation(validation, "success case").
		Return(nil).AnyTimes()
	// ErrInvalidNumber
	mockValidation.EXPECT().
		CheckNumberValidation(validation, "ErrInvalidNumber").
		Return(model.ErrInvalidNumber).AnyTimes()
	// BadRequest
	mockValidation.EXPECT().
		CheckNumberValidation(validation, "BadRequest").
		Return(errMock).AnyTimes()

	// CheckTextValidation
	// success
	mockValidation.EXPECT().
		CheckTextValidation(validation, "success case").
		Return(nil).AnyTimes()
	// ErrTextMatching
	mockValidation.EXPECT().
		CheckTextValidation(validation, "ErrTextMatching").
		Return(model.ErrTextMatching).AnyTimes()
	// InternalServerError
	mockValidation.EXPECT().
		CheckTextValidation(validation, "InternalServerError").
		Return(errMock).AnyTimes()

	// ScaleLabel
	// GetScaleLabels
	// success
	mockScaleLabel.EXPECT().
		GetScaleLabels(gomock.Any(), []int{questionIDSuccess}).
		Return([]model.ScaleLabels{scalelabel}, nil).AnyTimes()
	// failure
	mockScaleLabel.EXPECT().
		GetScaleLabels(gomock.Any(), []int{questionIDFailure}).
		Return([]model.ScaleLabels{}, nil).AnyTimes()
	// nothing
	mockScaleLabel.EXPECT().
		GetScaleLabels(gomock.Any(), []int{}).
		Return([]model.ScaleLabels{}, nil).AnyTimes()

	// CheckScaleLabel
	// success
	mockScaleLabel.EXPECT().
		CheckScaleLabel(scalelabel, "success case").
		Return(nil).AnyTimes()
	// BadRequest
	mockScaleLabel.EXPECT().
		CheckScaleLabel(scalelabel, "BadRequest").
		Return(errMock).AnyTimes()

	// Respondent
	// InsertRespondent
	// success
	mockRespondent.EXPECT().
		InsertRespondent(gomock.Any(), string(userOne), questionnaireIDSuccess, gomock.Any()).
		Return(responseIDSuccess, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		InsertRespondent(gomock.Any(), string(userOne), questionnaireIDFailure, gomock.Any()).
		Return(responseIDFailure, nil).AnyTimes()
	// UpdateSubmittedAt
	// success
	mockRespondent.EXPECT().
		UpdateSubmittedAt(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	// Response
	// InsertResponses
	// success
	mockResponse.EXPECT().
		InsertResponses(gomock.Any(), responseIDSuccess, gomock.Any()).
		Return(nil).AnyTimes()
	// failure
	mockResponse.EXPECT().
		InsertResponses(gomock.Any(), responseIDFailure, gomock.Any()).
		Return(errMock).AnyTimes()
	// DeleteResponse
	// success
	mockResponse.EXPECT().
		DeleteResponse(gomock.Any(), responseIDSuccess).
		Return(nil).AnyTimes()
	// failure
	mockResponse.EXPECT().
		DeleteResponse(gomock.Any(), responseIDFailure).
		Return(model.ErrNoRecordDeleted).AnyTimes()

	// responseID, err := mockRespondent.
	// 	InsertRespondent(string(userOne), 1, null.NewTime(nowTime, true))
	// assertion.Equal(1, responseID)
	// assertion.NoError(err)

	type request struct {
		user             users
		responseID       int
		isBadRequestBody bool
		requestBody      responseRequestBody
	}
	type expect struct {
		isErr bool
		code  int
	}
	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "null submittedat",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "empty body",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body:            []responseBody{},
				},
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "bad request body",
			request: request{
				isBadRequestBody: true,
				responseID:       responseIDSuccess,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "limit exceeded",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDLimit,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body:            []responseBody{},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusMethodNotAllowed,
			},
		},
		{
			description: "valid number",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Number",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "invalid number",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Number",
							Body:           null.StringFrom("ErrInvalidNumber"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "BadRequest number",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Number",
							Body:           null.StringFrom("BadRequest"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "valid text",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "text does not match",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("ErrTextMatching"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "invalid text",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("InternalServerError"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "valid LinearScale",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "LinearScale",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "invalid LinearScale",
			request: request{
				user:       userOne,
				responseID: responseIDSuccess,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.NewTime(nowTime, false),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "LinearScale",
							Body:           null.StringFrom("BadRequest"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
		{
			description: "response does not exist",
			request: request{
				user:       userOne,
				responseID: responseIDFailure,
				requestBody: responseRequestBody{
					QuestionnaireID: questionnaireIDSuccess,
					SubmittedAt:     null.TimeFrom(nowTime),
					Body: []responseBody{
						{
							QuestionID:     questionIDSuccess,
							QuestionType:   "Text",
							Body:           null.StringFrom("success case"),
							OptionResponse: []string{},
						},
					},
				},
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError, //middlewareで弾くので500で良い
			},
		},
	}

	e := echo.New()
	e.PATCH("/api/responses/:responseID", r.EditResponse, m.SetUserIDMiddleware, m.TraPMemberAuthenticate, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			responseID, err := strconv.Atoi(c.Param("responseID"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, "responseID is not number")
			}

			c.Set(responseIDKey, responseID)
			return next(c)
		}
	})

	for _, testCase := range testCases {
		requestByte, jsonErr := json.Marshal(testCase.request.requestBody)
		require.NoError(t, jsonErr)
		requestStr := string(requestByte) + "\n"

		if testCase.request.isBadRequestBody {
			requestStr = "badRequestBody"
		}
		rec := createRecorder(e, testCase.request.user, methodPatch, makePath(fmt.Sprint("/responses/", testCase.request.responseID)), typeJSON, requestStr)

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
	}
}

func TestDeleteResponse(t *testing.T) {
	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)

	r := NewResponse(
		mockQuestionnaire,
		mockValidation,
		mockScaleLabel,
		mockRespondent,
		mockResponse,
	)

	type request struct {
		QuestionnaireLimit         null.Time
		GetQuestionnaireLimitError error
		ExecutesDeletion           bool
		DeleteRespondentError      error
		DeleteResponseError        error
	}
	type expect struct {
		statusCode int
	}
	type test struct {
		description string
		request
		expect
	}

	testCases := []test{
		{
			description: "期限が設定されていない、かつDeleteRespondentがエラーなしなので200",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Time{}, false),
				GetQuestionnaireLimitError: nil,
				ExecutesDeletion:           true,
				DeleteRespondentError:      nil,
				DeleteResponseError:        nil,
			},
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			description: "期限前、かつDeleteRespondentがエラーなしなので200",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Now().AddDate(0, 0, 1), true),
				GetQuestionnaireLimitError: nil,
				ExecutesDeletion:           true,
				DeleteRespondentError:      nil,
				DeleteResponseError:        nil,
			},
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			description: "期限後なので405",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Now().AddDate(0, 0, -1), true),
				GetQuestionnaireLimitError: nil,
				ExecutesDeletion:           false,
				DeleteRespondentError:      nil,
				DeleteResponseError:        nil,
			},
			expect: expect{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			description: "GetQuestionnaireLimitByResponseIDがエラーRecordNotFoundを吐くので404",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Time{}, false),
				GetQuestionnaireLimitError: model.ErrRecordNotFound,
				ExecutesDeletion:           false,
				DeleteRespondentError:      nil,
				DeleteResponseError:        nil,
			},
			expect: expect{
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "GetQuestionnaireLimitByResponseIDがエラーを吐くので500",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Time{}, false),
				GetQuestionnaireLimitError: errors.New("error"),
				ExecutesDeletion:           false,
				DeleteRespondentError:      nil,
				DeleteResponseError:        nil,
			},
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "DeleteRespondentがエラーを吐くので500",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Time{}, false),
				GetQuestionnaireLimitError: nil,
				ExecutesDeletion:           true,
				DeleteRespondentError:      errors.New("error"),
				DeleteResponseError:        nil,
			},
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "DeleteResponseがエラーを吐くので500",
			request: request{
				QuestionnaireLimit:         null.NewTime(time.Time{}, false),
				GetQuestionnaireLimitError: nil,
				ExecutesDeletion:           true,
				DeleteRespondentError:      nil,
				DeleteResponseError:        errors.New("error"),
			},
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, testCase := range testCases {
		userID := "userID1"
		responseID := 1

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/responses/%d", responseID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/responses/:responseID")
		c.SetParamNames("responseID")
		c.SetParamValues(strconv.Itoa(responseID))
		c.Set(userIDKey, userID)
		c.Set(responseIDKey, responseID)

		mockQuestionnaire.
			EXPECT().
			GetQuestionnaireLimitByResponseID(gomock.Any(), responseID).
			Return(testCase.request.QuestionnaireLimit, testCase.request.GetQuestionnaireLimitError)
		if testCase.request.ExecutesDeletion {
			mockRespondent.
				EXPECT().
				DeleteRespondent(gomock.Any(), responseID).
				Return(testCase.request.DeleteRespondentError)
			if testCase.request.DeleteRespondentError == nil {
				mockResponse.
					EXPECT().
					DeleteResponse(c.Request().Context(), responseID).
					Return(testCase.request.DeleteResponseError)
			}
		}

		e.HTTPErrorHandler(r.DeleteResponse(c), c)

		assertion.Equal(testCase.expect.statusCode, rec.Code, testCase.description, "status code")
	}
}
