package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

type responseRequestBody struct {
	SubmittedAt null.Time      `json:"submitted_at"`
	ModifiedAt  null.Time      `json:"modified_at"`
	Body        []responseBody `json:"body"`
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

func TestGetResponse(t *testing.T) {

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()
	responseIDSuccess := 1
	questionnaireIDSuccess := 1
	questionIDSuccess := 1
	responseIDFailure := 0
	responseIDNotFound := -1

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
