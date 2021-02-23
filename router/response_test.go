package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

type responseResponseBody struct {
	QuestionnaireID int            `json:"questionnaireID"`
	SubmittedAt     null.Time      `json:"submitted_at"`
	ModifiedAt      null.Time      `json:"modified_at"`
	Body            []responseBody `json:"body"`
}

var (
	errMock = errors.New("Mock Error")
)

type users string
type httpMethods string
type contentTypes string

const (
	rootPath               = "/api"
	userHeader             = "X-Showcase-User"
	userUnAuthorized       = "-"
	userOne          users = "mazrean"
	userTwo          users = "ryoha"
	//userThree        users        = "YumizSui"
	methodGet  httpMethods = http.MethodGet
	methodPost httpMethods = http.MethodPost
	//methodPatch      httpMethods  = http.MethodPatch
	//methodDelete      httpMethods  = http.MethodDelete
	typeNone contentTypes = ""
	typeJSON contentTypes = echo.MIMEApplicationJSON
)

func makePath(path string) string {
	return rootPath + path
}

func createRecorder(e *echo.Echo, user users, method httpMethods, path string, contentType contentTypes, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(string(method), path, strings.NewReader(body))
	if contentType != typeNone {
		req.Header.Set(echo.HeaderContentType, string(contentType))
	}
	req.Header.Set(userHeader, string(user))

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
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
	r := NewResponse(
		mockQuestionnaire,
		mockValidation,
		mockScaleLabel,
		mockRespondent,
		mockResponse,
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

	type args struct {
		responseID int
	}
	type expect struct {
		isErr    bool
		code     int
		response responseResponseBody
	}

	type test struct {
		description string
		args
		expect
	}
	testCases := []test{
		{
			description: "success",
			args: args{
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
			args: args{
				responseID: responseIDFailure,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "NotFound",
			args: args{
				responseID: responseIDNotFound,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusNotFound,
			},
		},
	}

	e := echo.New()
	e.GET("/api/responses/:responseID", r.GetResponse)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/api/responses/", testCase.args.responseID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
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
	// testList := []struct {
	// 	description string
	// 	responseID  int
	// 	request     responseRequestBody
	// 	expectCode  int
	// }{
	// 	{
	// 		description: "valid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Text",
	// 					body:           null.StringFrom("回答"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "response does not exist",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Text",
	// 					body:           null.StringFrom("回答"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusNotFound,
	// 	},
	// 	{
	// 		description: "null submittedat",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.NewTime(time.Now(), false),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Text",
	// 					body:           null.StringFrom("回答"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "empty body",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body:        []responseBody{},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "empty body",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body:        []responseBody{},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "question does not exist",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Text",
	// 					body:           null.StringFrom("回答"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "number valid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Number",
	// 					body:           null.StringFrom("10"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "number invalid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Number",
	// 					body:           null.StringFrom("-1000"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusBadRequest,
	// 	},
	// 	{
	// 		description: "text valid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Text",
	// 					body:           null.StringFrom("1000"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "text invalid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Number",
	// 					body:           null.StringFrom("100a"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusBadRequest,
	// 	},
	// 	{
	// 		description: "LinearScale valid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Number",
	// 					body:           null.StringFrom("1"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusOK,
	// 	},
	// 	{
	// 		description: "LinearScale invalid",
	// 		responseID:  -1,
	// 		request: responseRequestBody{
	// 			submittedAt: null.TimeFrom(time.Now()),
	// 			body: []responseBody{
	// 				{
	// 					questionID:     -1,
	// 					questionType:   "Number",
	// 					body:           null.StringFrom("-1"),
	// 					optionResponse: []string{},
	// 				},
	// 			},
	// 		},
	// 		expectCode: http.StatusBadRequest,
	// 	},
	// }
	// fmt.Println(testList)
}

func TestDeleteResponse(t *testing.T) {
	// testList := []struct {
	// 	description string
	// 	responseID  int
	// 	request     responseRequestBody
	// 	expectCode  int
	// }{
	// 	{
	// 		description: "valid",
	// 		responseID:  -1,
	// 		expectCode:  http.StatusOK,
	// 	},
	// 	{
	// 		description: "response not exist",
	// 		responseID:  -1,
	// 		expectCode:  http.StatusNotFound,
	// 	},
	// }
	// fmt.Println(testList)
}

// func (p *responseBody) createResponseBody() string {
// 	optionResponses := make([]string, 0, len(p.optionResponse))
// 	for _, optionResponse := range p.optionResponse {
// 		optionResponses = append(optionResponses, fmt.Sprintf("\"%s\"", optionResponse))
// 	}

// 	return fmt.Sprintf(
// 		`{
//   "questionID": %v,
//   "question_type": %v,
//   "response": "%s",
//   "option_response": [%s],
// }`, p.questionID, p.questionType, p.body.String, strings.Join(optionResponses, ",\n    "))
// }

// func (p *responseRequestBody) createResponseRequestBody() string {
// 	bodies := make([]string, 0, len(p.body))
// 	for _, body := range p.body {
// 		bodies = append(bodies, body.createResponseBody())
// 	}

// 	return fmt.Sprintf(
// 		`{
//   "questionnaireID": %v,
//   "submitted_at": %v,
//   "body": [%s],
// }`, p.id, p.submittedAt, strings.Join(bodies, ",\n    "))
// }
