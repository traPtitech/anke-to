package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/model/mock_model"
	"gopkg.in/guregu/null.v3"
)

type responseBody struct {
	QuestionID     int         `json:"questionID"`
	QuestionType   string      `json:"question_type"`
	Body           null.String `json:"response"`
	OptionResponse []string    `json:"option_response"`
}

func TestPostResponse(t *testing.T) {
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
	)
	// Questionnaire
	// GetQuestionnaireLimit
	// success
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(questionnaireIDSuccess).
		Return(null.TimeFrom(nowTime.Add(time.Minute)), nil).AnyTimes()
	// failure
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(questionnaireIDFailure).
		Return(null.NewTime(time.Time{}, false), gorm.ErrRecordNotFound).AnyTimes()
	// limit
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(questionnaireIDLimit).
		Return(null.TimeFrom(nowTime.Add(-time.Minute)), nil).AnyTimes()

	// Validation
	// GetValidations
	// success
	mockValidation.EXPECT().
		GetValidations([]int{questionIDSuccess}).
		Return([]model.Validations{validation}, nil).AnyTimes()
	// failure
	mockValidation.EXPECT().
		GetValidations([]int{questionIDFailure}).
		Return([]model.Validations{}, nil).AnyTimes()
	// nothing
	mockValidation.EXPECT().
		GetValidations([]int{}).
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
		GetScaleLabels([]int{questionIDSuccess}).
		Return([]model.ScaleLabels{scalelabel}, nil).AnyTimes()
	// failure
	mockScaleLabel.EXPECT().
		GetScaleLabels([]int{questionIDFailure}).
		Return([]model.ScaleLabels{}, nil).AnyTimes()
	// nothing
	mockScaleLabel.EXPECT().
		GetScaleLabels([]int{}).
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
		InsertRespondent(string(userOne), questionnaireIDSuccess, gomock.Any()).
		Return(responseIDSuccess, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		InsertRespondent(string(userOne), questionnaireIDFailure, gomock.Any()).
		Return(responseIDFailure, nil).AnyTimes()

	// Response
	// InsertResponses
	// success
	mockResponse.EXPECT().
		InsertResponses(responseIDSuccess, gomock.Any()).
		Return(nil).AnyTimes()
	// failure
	mockResponse.EXPECT().
		InsertResponses(responseIDFailure, gomock.Any()).
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
				code:  http.StatusNotFound,
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
	e.POST("/api/responses", r.PostResponse, m.UserAuthenticate)

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
	)

	// Respondent
	// InsertRespondent
	// success
	mockRespondent.EXPECT().
		GetRespondentDetail(responseIDSuccess).
		Return(respondentDetail, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		GetRespondentDetail(responseIDFailure).
		Return(model.RespondentDetail{}, errMock).AnyTimes()
	// NotFound
	mockRespondent.EXPECT().
		GetRespondentDetail(responseIDNotFound).
		Return(model.RespondentDetail{}, gorm.ErrRecordNotFound).AnyTimes()

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
	e.GET("/api/responses/:responseID", r.GetResponse, m.UserAuthenticate)

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
	)
	// Questionnaire
	// GetQuestionnaireLimit
	// success
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(questionnaireIDSuccess).
		Return(null.TimeFrom(nowTime.Add(time.Minute)), nil).AnyTimes()
	// failure
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(questionnaireIDFailure).
		Return(null.NewTime(time.Time{}, false), errMock).AnyTimes()
	// limit
	mockQuestionnaire.EXPECT().
		GetQuestionnaireLimit(questionnaireIDLimit).
		Return(null.TimeFrom(nowTime.Add(-time.Minute)), nil).AnyTimes()

	// Validation
	// GetValidations
	// success
	mockValidation.EXPECT().
		GetValidations([]int{questionIDSuccess}).
		Return([]model.Validations{validation}, nil).AnyTimes()
	// failure
	mockValidation.EXPECT().
		GetValidations([]int{questionIDFailure}).
		Return([]model.Validations{}, nil).AnyTimes()
	// nothing
	mockValidation.EXPECT().
		GetValidations([]int{}).
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
		GetScaleLabels([]int{questionIDSuccess}).
		Return([]model.ScaleLabels{scalelabel}, nil).AnyTimes()
	// failure
	mockScaleLabel.EXPECT().
		GetScaleLabels([]int{questionIDFailure}).
		Return([]model.ScaleLabels{}, nil).AnyTimes()
	// nothing
	mockScaleLabel.EXPECT().
		GetScaleLabels([]int{}).
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
		InsertRespondent(string(userOne), questionnaireIDSuccess, gomock.Any()).
		Return(responseIDSuccess, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		InsertRespondent(string(userOne), questionnaireIDFailure, gomock.Any()).
		Return(responseIDFailure, nil).AnyTimes()
	// CheckRespondentByResponseID
	// success
	mockRespondent.EXPECT().
		CheckRespondentByResponseID(gomock.Any(), responseIDSuccess).
		Return(true, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		CheckRespondentByResponseID(gomock.Any(), responseIDFailure).
		Return(false, nil).AnyTimes()
	// UpdateSubmittedAt
	// success
	mockRespondent.EXPECT().
		UpdateSubmittedAt(gomock.Any()).
		Return(nil).AnyTimes()

	// Response
	// InsertResponses
	// success
	mockResponse.EXPECT().
		InsertResponses(responseIDSuccess, gomock.Any()).
		Return(nil).AnyTimes()
	// failure
	mockResponse.EXPECT().
		InsertResponses(responseIDFailure, gomock.Any()).
		Return(errMock).AnyTimes()
	// DeleteResponse
	// success
	mockResponse.EXPECT().
		DeleteResponse(responseIDSuccess).
		Return(nil).AnyTimes()
	// failure
	mockResponse.EXPECT().
		DeleteResponse(responseIDFailure).
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
			description: "response doe not exist",
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
				code:  http.StatusForbidden,
			},
		},
	}

	e := echo.New()
	e.PATCH("/api/responses/:responseID", r.EditResponse, m.UserAuthenticate, m.RespondentAuthenticate)

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

	responseIDSuccess := 1
	responseIDFailure := 0

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
	)

	// Respondent
	// InsertRespondent
	// success
	mockRespondent.EXPECT().
		DeleteRespondent(gomock.Any(), responseIDSuccess).
		Return(nil).AnyTimes()
	// success
	mockRespondent.EXPECT().
		DeleteRespondent(gomock.Any(), responseIDSuccess).
		Return(model.ErrNoRecordDeleted).AnyTimes()
	// CheckRespondentByResponseID
	// success
	mockRespondent.EXPECT().
		CheckRespondentByResponseID(gomock.Any(), responseIDSuccess).
		Return(true, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		CheckRespondentByResponseID(gomock.Any(), responseIDFailure).
		Return(false, nil).AnyTimes()

	type request struct {
		user       users
		responseID int
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
				responseID: responseIDSuccess,
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
			},
		},
		{
			description: "response does not exist",
			request: request{
				responseID: responseIDFailure,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusForbidden,
			},
		},
	}

	e := echo.New()
	e.DELETE("/api/responses/:responseID", r.DeleteResponse, m.UserAuthenticate, m.RespondentAuthenticate)

	for _, testCase := range testCases {
		rec := createRecorder(e, testCase.request.user, methodDelete, fmt.Sprint(rootPath, "/responses/", testCase.request.responseID), typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
	}
}