package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/model/mock_model"
	"gopkg.in/guregu/null.v3"
)

func TestGetResults(t *testing.T) {
	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	result := NewResult(mockRespondent, mockQuestionnaire, mockAdministrator)

	type request struct {
		sortParam                 string
		questionnaireIDParam      string
		questionnaireIDValid      bool
		questionnaireID           int
		respondentDetails         []model.RespondentDetail
		getRespondentDetailsError error
	}
	type response struct {
		statusCode int
		body       string
	}
	type test struct {
		description string
		request
		response
	}

	textResponse := []model.RespondentDetail{
		{
			ResponseID:      1,
			TraqID:          "mazrean",
			QuestionnaireID: 1,
			SubmittedAt:     null.NewTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			ModifiedAt:      time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Responses: []model.ResponseBody{
				{
					QuestionID:     1,
					QuestionType:   "Text",
					Body:           null.NewString("テスト", true),
					OptionResponse: nil,
				},
			},
		},
	}
	sb := strings.Builder{}
	err := json.NewEncoder(&sb).Encode(textResponse)
	if err != nil {
		t.Errorf("failed to encode text response: %v", err)
		return
	}

	testCases := []test{
		{
			description: "questionnaireIDが数字でないので400",
			request: request{
				sortParam:            "traqid",
				questionnaireIDParam: "abc",
			},
			response: response{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "questionnaireIDが空文字なので400",
			request: request{
				sortParam:            "traqid",
				questionnaireIDParam: "",
			},
			response: response{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "GetRespondentDetailsがエラーなので500",
			request: request{
				sortParam:                 "traqid",
				questionnaireIDValid:      true,
				questionnaireIDParam:      "1",
				questionnaireID:           1,
				getRespondentDetailsError: fmt.Errorf("error"),
			},
			response: response{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "respondentDetailsがnilでも200",
			request: request{
				sortParam:            "traqid",
				questionnaireIDValid: true,
				questionnaireIDParam: "1",
				questionnaireID:      1,
			},
			response: response{
				statusCode: http.StatusOK,
				body:       "null\n",
			},
		},
		{
			description: "respondentDetailsがそのまま帰り200",
			request: request{
				sortParam:            "traqid",
				questionnaireIDValid: true,
				questionnaireIDParam: "1",
				questionnaireID:      1,
				respondentDetails: []model.RespondentDetail{
					{
						ResponseID:      1,
						TraqID:          "mazrean",
						QuestionnaireID: 1,
						SubmittedAt:     null.NewTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), true),
						ModifiedAt:      time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
						Responses: []model.ResponseBody{
							{
								QuestionID:     1,
								QuestionType:   "Text",
								Body:           null.NewString("テスト", true),
								OptionResponse: nil,
							},
						},
					},
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body:       sb.String(),
			},
		},
	}

	for _, testCase := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/results/%s?sort=%s", testCase.request.questionnaireIDParam, testCase.request.sortParam), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/results/:questionnaireID")
		c.SetParamNames("questionnaireID", "sort")
		c.SetParamValues(testCase.request.questionnaireIDParam, testCase.request.sortParam)

		if testCase.request.questionnaireIDValid {
			mockRespondent.
				EXPECT().
				GetRespondentDetails(testCase.request.questionnaireID, testCase.request.sortParam).
				Return(testCase.request.respondentDetails, testCase.request.getRespondentDetailsError)
		}

		e.HTTPErrorHandler(result.GetResults(c), c)
		assertion.Equalf(testCase.response.statusCode, rec.Code, testCase.description, "statusCode")
		if testCase.response.statusCode == http.StatusOK {
			assertion.Equalf(testCase.response.body, rec.Body.String(), testCase.description, "body")
		}
	}
}
