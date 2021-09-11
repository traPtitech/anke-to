package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/model/mock_model"
	"gopkg.in/guregu/null.v3"
)

func TestGetResults(t *testing.T) {

	type resultResponseBody struct {
		QuestionnaireID int            `json:"questionnaireID"`
		SubmittedAt     null.Time      `json:"submitted_at"`
		ModifiedAt      time.Time      `json:"modified_at"`
		Body            []responseBody `json:"body"`
	}

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()
	questionnaireIDSuccess := 1
	// responseIDSuccess := 1
	questionIDSuccess := 1

	questionnaireIDFailure := 0
	questionnaireIDNotFound := -1
	response1 := model.ResponseBody{
		QuestionID:     questionIDSuccess,
		QuestionType:   "Text",
		Body:           null.StringFrom("回答1"),
		OptionResponse: []string{},
	}
	response2 := model.ResponseBody{
		QuestionID:     questionIDSuccess,
		QuestionType:   "Text",
		Body:           null.StringFrom("回答2"),
		OptionResponse: []string{},
	}
	response3 := model.ResponseBody{
		QuestionID:     questionIDSuccess,
		QuestionType:   "Text",
		Body:           null.StringFrom("回答3"),
		OptionResponse: []string{},
	}
	respondentDetails := []model.RespondentDetail{
		model.RespondentDetail{
			QuestionnaireID: questionnaireIDSuccess,
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
			Responses:       []model.ResponseBody{response1},
		},
		model.RespondentDetail{
			QuestionnaireID: questionnaireIDSuccess,
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
			Responses:       []model.ResponseBody{response2},
		},
		model.RespondentDetail{
			QuestionnaireID: questionnaireIDSuccess,
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
			Responses:       []model.ResponseBody{response3},
		},
	}

	resultResponseBodies := []resultResponseBody{
		{
			QuestionnaireID: questionnaireIDSuccess,
			Body: []responseBody{
				responseBody{
					QuestionID:     response1.QuestionID,
					QuestionType:   response1.QuestionType,
					Body:           response1.Body,
					OptionResponse: response1.OptionResponse,
				},
			},
			SubmittedAt: null.TimeFrom(nowTime),
			ModifiedAt:  nowTime,
		},
		{
			QuestionnaireID: questionnaireIDSuccess,
			Body: []responseBody{
				responseBody{
					QuestionID:     response2.QuestionID,
					QuestionType:   response2.QuestionType,
					Body:           response2.Body,
					OptionResponse: response2.OptionResponse,
				},
			},
			SubmittedAt: null.TimeFrom(nowTime),
			ModifiedAt:  nowTime,
		},
		{
			QuestionnaireID: questionnaireIDSuccess,
			Body: []responseBody{
				responseBody{
					QuestionID:     response3.QuestionID,
					QuestionType:   response3.QuestionType,
					Body:           response3.Body,
					OptionResponse: response3.OptionResponse,
				},
			},
			SubmittedAt: null.TimeFrom(nowTime),
			ModifiedAt:  nowTime,
		},
	}

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	mockQuestion := mock_model.NewMockIQuestion(ctrl)

	r := NewResult(
		mockRespondent,
		mockQuestionnaire,
		mockAdministrator,
	)

	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
	)

	// Respondent
	// InsertRespondent
	mockRespondent.EXPECT().
		GetRespondentDetails(questionnaireIDSuccess, gomock.Any()).
		Return(respondentDetails, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		GetRespondentDetails(questionnaireIDFailure, gomock.Any()).
		Return(nil, errMock).AnyTimes()
	// NotFound
	mockRespondent.EXPECT().
		GetRespondentDetails(questionnaireIDNotFound, gomock.Any()).
		Return([]model.RespondentDetail{}, nil).AnyTimes()

	// Questionnaire
	// GetResShared
	// failure
	mockQuestionnaire.EXPECT().
		GetResShared(questionnaireIDFailure).
		Return("", errMock).AnyTimes()
	// NotFound
	mockQuestionnaire.EXPECT().
		GetResShared(questionnaireIDNotFound).
		Return("", gorm.ErrRecordNotFound).AnyTimes()

	type request struct {
		user            users
		questionnaireID int
	}
	type call struct {
		resSharedTo  string
		isAdmin      bool
		isRespondent bool
	}
	type expect struct {
		isErr    bool
		code     int
		response []resultResponseBody
	}

	type test struct {
		description string
		request
		call
		expect
	}
	testCases := []test{
		{
			description: "public",
			request: request{
				questionnaireID: questionnaireIDSuccess,
			},
			call: call{
				resSharedTo: "public",
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: resultResponseBodies,
			},
		},
		{
			description: "administrators admin",
			request: request{
				questionnaireID: questionnaireIDSuccess,
			},
			call: call{
				resSharedTo: "administrators",
				isAdmin:     true,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: resultResponseBodies,
			},
		},
		{
			description: "administrators not Admin",
			request: request{
				questionnaireID: questionnaireIDSuccess,
			},
			call: call{
				resSharedTo: "administrators",
				isAdmin:     false,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusUnauthorized,
			},
		},
		{
			description: "respondents admin",
			request: request{
				questionnaireID: questionnaireIDSuccess,
			},
			call: call{
				resSharedTo: "respondents",
				isAdmin:     true,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: resultResponseBodies,
			},
		},
		{
			description: "respondents respondent",
			request: request{
				questionnaireID: questionnaireIDSuccess,
			},
			call: call{
				resSharedTo:  "respondents",
				isRespondent: true,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: resultResponseBodies,
			},
		},
		{
			description: "respondents not admin or respondent",
			request: request{
				questionnaireID: questionnaireIDSuccess,
			},
			call: call{
				resSharedTo:  "respondents",
				isAdmin:      false,
				isRespondent: false,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusUnauthorized,
			},
		},
		{
			description: "failure",
			request: request{
				questionnaireID: questionnaireIDFailure,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "NotFound",
			request: request{
				questionnaireID: questionnaireIDNotFound,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusNotFound,
			},
		},
	}

	e := echo.New()
	e.GET("/api/results/:questionnaireID", r.GetResults, m.UserAuthenticate)

	for _, testCase := range testCases {
		if testCase.request.questionnaireID == questionnaireIDSuccess {
			// GetResShared
			mockQuestionnaire.EXPECT().
				GetResShared(questionnaireIDSuccess).
				Return(testCase.call.resSharedTo, nil)
		}
		if testCase.call.resSharedTo != "" {

			if testCase.call.resSharedTo == "administrators" || testCase.call.resSharedTo == "respondents" {
				// CheckQuestionnaireAdmin
				mockAdministrator.EXPECT().
					CheckQuestionnaireAdmin(gomock.Any(), testCase.request.questionnaireID).
					Return(testCase.call.isAdmin, nil)
			}

			if testCase.call.resSharedTo == "respondents" && !testCase.call.isAdmin {
				// CheckRespondent
				mockRespondent.EXPECT().
					CheckRespondent(gomock.Any(), testCase.request.questionnaireID).
					Return(testCase.call.isRespondent, nil)
			}
		}

		rec := createRecorder(e, testCase.request.user, methodGet, fmt.Sprint(rootPath, "/results/", testCase.request.questionnaireID), typeNone, "")

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
