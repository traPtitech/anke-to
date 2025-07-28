package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
	"gopkg.in/guregu/null.v4"
)

var (
	sampleResponseBodyText            = openapi.ResponseBody{}
	sampleResponseBodyTextLong        = openapi.ResponseBody{}
	sampleResponseBodyNumber          = openapi.ResponseBody{}
	sampleResponseBodySingleChoice    = openapi.ResponseBody{}
	sampleResponseBodyMultipleChoice  = openapi.ResponseBody{}
	sampleResponseBodyScale           = openapi.ResponseBody{}
	sampleResponse                    = openapi.NewResponse{}
	AddQuestionID2SampleResponseMutex sync.Mutex
)

func setupSampleResponse() {
	if len(sampleResponse.Body) > 0 {
		return
	}
	sampleResponseBodyText.FromResponseBodyText(openapi.ResponseBodyText{
		Answer:       "テキスト",
		QuestionType: "Text",
	})
	sampleResponseBodyTextLong.FromResponseBodyTextLong(openapi.ResponseBodyTextLong{
		Answer:       "ロングテキスト",
		QuestionType: "TextLong",
	})
	sampleResponseBodyNumber.FromResponseBodyNumber(openapi.ResponseBodyNumber{
		Answer:       0,
		QuestionType: "Number",
	})
	sampleResponseBodySingleChoice.FromResponseBodySingleChoice(openapi.ResponseBodySingleChoice{
		Answer:       1,
		QuestionType: "SingleChoice",
	})
	sampleResponseBodyMultipleChoice.FromResponseBodyMultipleChoice(openapi.ResponseBodyMultipleChoice{
		Answer:       []int{1, 2},
		QuestionType: "MultipleChoice",
	})
	sampleResponseBodyScale.FromResponseBodyScale(openapi.ResponseBodyScale{
		Answer:       1,
		QuestionType: "Scale",
	})
	sampleResponse = openapi.NewResponse{
		Body: []openapi.ResponseBody{
			sampleResponseBodyText,
			sampleResponseBodyTextLong,
			sampleResponseBodyNumber,
			sampleResponseBodySingleChoice,
			sampleResponseBodyMultipleChoice,
			sampleResponseBodyScale,
		},
		IsDraft: false,
	}
}

// sampleResponseのResponseBodyでquestionnaireIDに基づいてquestionIDを設定する
func AddQuestionID2SampleResponse(questionnaireID int) {
	questions, err := q.IQuestion.GetQuestions(context.Background(), questionnaireID)
	if err != nil {
		panic(fmt.Sprintf("failed to get questions: %v", err))
	}
	sampleResponseBodyText.QuestionId = questions[0].ID
	sampleResponseBodyTextLong.QuestionId = questions[1].ID
	sampleResponseBodyNumber.QuestionId = questions[2].ID
	sampleResponseBodySingleChoice.QuestionId = questions[3].ID
	sampleResponseBodyMultipleChoice.QuestionId = questions[4].ID
	sampleResponseBodyScale.QuestionId = questions[5].ID
	sampleResponse = openapi.NewResponse{
		Body: []openapi.ResponseBody{
			sampleResponseBodyText,
			sampleResponseBodyTextLong,
			sampleResponseBodyNumber,
			sampleResponseBodySingleChoice,
			sampleResponseBodyMultipleChoice,
			sampleResponseBodyScale,
		},
		IsDraft: false,
	}
	for i := range sampleResponse.Body {
		sampleResponse.Body[i].QuestionId = questions[i].ID
	}
}

func TestGetMyResponses(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaire := sampleQuestionnaire
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	AddQuestionID2SampleResponseMutex.Lock()

	AddQuestionID2SampleResponse(questionnaireDetail.QuestionnaireId)

	newResponse := sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	newResponse.IsDraft = true
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userTwo)
	require.NoError(t, err)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response4, err := q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, "myResponsesSpecialUser")
	require.NoError(t, err)

	AddQuestionID2SampleResponseMutex.Unlock()

	type args struct {
		userID string
		params openapi.GetMyResponsesParams
	}
	type expect struct {
		isErr          bool
		err            error
		responseIdList *[]int
	}
	type test struct {
		description string
		args
		expect
	}

	sortInvalid := (openapi.ResponseSortInQuery)("abcde")
	sortTraqID := (openapi.ResponseSortInQuery)("traqid")
	sortTraqIDDesc := (openapi.ResponseSortInQuery)("-traqid")
	sortSubmittedAt := (openapi.ResponseSortInQuery)("submitted_at")
	sortSubmittedAtDesc := (openapi.ResponseSortInQuery)("-submitted_at")
	sortModifiedAt := (openapi.ResponseSortInQuery)("modified_at")
	sortModifiedAtDesc := (openapi.ResponseSortInQuery)("-modified_at")

	testCases := []test{
		{
			description: "valid",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{},
			},
		},
		{
			description: "invalid param sort",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortInvalid,
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "sort submitted_at",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortSubmittedAt,
				},
			},
		},
		{
			description: "sort -submitted_at",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortSubmittedAtDesc,
				},
			},
		},
		{
			description: "sort traqid",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortTraqID,
				},
			},
		},
		{
			description: "sort -traqid",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortTraqIDDesc,
				},
			},
		},
		{
			description: "sort modified_at",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortModifiedAt,
				},
			},
		},
		{
			description: "sort -modified_at",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortModifiedAtDesc,
				},
			},
		},
		{
			description: "special user",
			args: args{
				userID: "myResponsesSpecialUser",
				params: openapi.GetMyResponsesParams{},
			},
			expect: expect{
				responseIdList: &[]int{response4.ResponseId},
			},
		},
		{
			description: "user with no record",
			args: args{
				userID: "myResponsesNoRecord",
				params: openapi.GetMyResponsesParams{},
			},
			expect: expect{
				responseIdList: &[]int{},
			},
		},
	}

	for _, testCase := range testCases {
		params := url.Values{}
		if testCase.args.params.Sort != nil {
			params.Add("sort", string(*testCase.args.params.Sort))
		}
		e = echo.New()
		req = httptest.NewRequest(http.MethodGet, "/responses/myResponses"+params.Encode(), nil)
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)

		responseLists, err := r.GetMyResponses(ctx, testCase.args.params, testCase.args.userID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		for _, responseList := range responseLists {
			if testCase.args.params.Sort != nil {
				if *testCase.args.params.Sort == "submitted_at" {
					var preCreatedAt time.Time
					for _, response := range *responseList.Responses {
						if !preCreatedAt.IsZero() {
							assertion.False(preCreatedAt.After(response.SubmittedAt), testCase.description, "submitted_at")
						}
						preCreatedAt = response.SubmittedAt
					}
				} else if *testCase.args.params.Sort == "-submitted_at" {
					var preCreatedAt time.Time
					for _, response := range *responseList.Responses {
						if !preCreatedAt.IsZero() {
							assertion.False(preCreatedAt.Before(response.SubmittedAt), testCase.description, "-submitted_at")
						}
						preCreatedAt = response.SubmittedAt
					}
				} else if *testCase.args.params.Sort == "traqid" {
					var preTraqID string
					for _, response := range *responseList.Responses {
						if preTraqID != "" {
							assertion.False(preTraqID > *response.Respondent, testCase.description, "traqid")
						}
						preTraqID = *response.Respondent
					}
				} else if *testCase.args.params.Sort == "-traqid" {
					var preTraqID string
					for _, response := range *responseList.Responses {
						if preTraqID != "" {
							assertion.False(preTraqID < *response.Respondent, testCase.description, "-traqid")
						}
						preTraqID = *response.Respondent
					}
				} else if *testCase.args.params.Sort == "modified_at" {
					var preModifiedAt time.Time
					for _, response := range *responseList.Responses {
						if !preModifiedAt.IsZero() {
							assertion.False(preModifiedAt.After(response.ModifiedAt), testCase.description, "modified_at")
						}
						preModifiedAt = response.ModifiedAt
					}
				} else if *testCase.args.params.Sort == "-modified_at" {
					var preModifiedAt time.Time
					for _, response := range *responseList.Responses {
						if !preModifiedAt.IsZero() {
							assertion.False(preModifiedAt.Before(response.ModifiedAt), testCase.description, "-modified_at")
						}
						preModifiedAt = response.ModifiedAt
					}
				}
			}
		}

		if testCase.expect.responseIdList != nil {
			responseIdList := []int{}
			for _, responseList := range responseLists {
				for _, response := range *responseList.Responses {
					responseIdList = append(responseIdList, response.ResponseId)
				}
			}
			sort.Slice(*testCase.expect.responseIdList, func(i, j int) bool {
				return (*testCase.expect.responseIdList)[i] < (*testCase.expect.responseIdList)[j]
			})
			sort.Slice(responseIdList, func(i, j int) bool { return responseIdList[i] < responseIdList[j] })
			assertion.Equal(*testCase.expect.responseIdList, responseIdList, testCase.description, "responseIdList")
		}

		for _, responseList := range responseLists {
			for _, response := range *responseList.Responses {
				assertion.Equal(testCase.args.userID, *response.Respondent, testCase.description, "response with no respondent")
			}
		}
	}
}

func TestGetResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaire := sampleQuestionnaire
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	AddQuestionID2SampleResponseMutex.Lock()

	AddQuestionID2SampleResponse(questionnaireDetail.QuestionnaireId)

	newResponse := sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response0, err := q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	questionnaireAnonymous := sampleQuestionnaire
	questionnaireAnonymous.IsAnonymous = true
	e = echo.New()
	body, err = json.Marshal(questionnaireAnonymous)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireAnonymousDetail, err := q.PostQuestionnaire(ctx, questionnaireAnonymous)
	require.NoError(t, err)

	AddQuestionID2SampleResponse(questionnaireAnonymousDetail.QuestionnaireId)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireAnonymousDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response1, err := q.PostQuestionnaireResponse(ctx, questionnaireAnonymousDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	AddQuestionID2SampleResponseMutex.Unlock()

	type args struct {
		isAnonymousQuestionnaire bool
		invalidResponseID        bool
	}
	type expect struct {
		isErr bool
		err   error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "valid",
		},
		{
			description: "invalid response id",
			args: args{
				invalidResponseID: true,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "anonymous questionnaire",
			args: args{
				isAnonymousQuestionnaire: true,
			},
		},
	}

	for _, testCase := range testCases {
		var responseID int
		if !testCase.args.invalidResponseID {
			if !testCase.args.isAnonymousQuestionnaire {
				responseID = response0.ResponseId
			} else {
				responseID = response1.ResponseId
			}
		} else {
			responseID = 10000
			valid := true
			for valid {
				c := context.Background()
				_, err := IRespondent.GetRespondent(c, responseID)
				if err == model.ErrRecordNotFound {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting respondent")
				} else {
					responseID *= 10
				}
			}
		}
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/responses/%d", responseID), nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		response, err := r.GetResponse(ctx, responseID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		if !testCase.args.isAnonymousQuestionnaire {
			assertion.Equal(response0.Body, response.Body, testCase.description, "response body")
			if response.Respondent != nil {
				assertion.Equal(*response0.IsAnonymous, *response.IsAnonymous, testCase.description, "response isAnonymous")
			} else {
				assertion.Equal(response0.IsAnonymous, response.IsAnonymous, testCase.description, "response isAnonymous")
			}
			assertion.Equal(response0.IsDraft, response.IsDraft, testCase.description, "response isDraft")
			assertion.WithinDuration(response0.ModifiedAt.UTC().Truncate(time.Second), response.ModifiedAt.UTC(), time.Second, testCase.description, "response modifiedAt")
			assertion.Equal(response0.QuestionnaireId, response.QuestionnaireId, testCase.description, "response questionnaireID")
			if response.Respondent != nil {
				assertion.Equal(*response0.Respondent, *response.Respondent, testCase.description, "response respondent")
			} else {
				assertion.Equal(response0.Respondent, response.Respondent, testCase.description, "response respondent")
			}
			assertion.Equal(response0.ResponseId, response.ResponseId, testCase.description, "response responseID")
			assertion.Equal(response0.QuestionnaireId, response.QuestionnaireId, testCase.description, "response questionnaireID")
			assertion.WithinDuration(response0.SubmittedAt.UTC().Truncate(time.Second), response.SubmittedAt.UTC(), time.Second, testCase.description, "response submittedAt")
		} else {
			assertion.Equal(response1.Body, response.Body, testCase.description, "response body")
			if response.Respondent != nil {
				assertion.Equal(*response1.IsAnonymous, *response.IsAnonymous, testCase.description, "response isAnonymous")
			} else {
				assertion.Equal(response1.IsAnonymous, response.IsAnonymous, testCase.description, "response isAnonymous")
			}
			assertion.Equal(response1.IsDraft, response.IsDraft, testCase.description, "response isDraft")
			assertion.WithinDuration(response1.ModifiedAt.UTC().Truncate(time.Second), response.ModifiedAt.UTC(), time.Second, testCase.description, "response modifiedAt")
			assertion.Equal(response1.QuestionnaireId, response.QuestionnaireId, testCase.description, "response questionnaireID")
			assertion.Nil(response.Respondent, testCase.description, "anonymous questionnaire respondent")
			assertion.Equal(response1.ResponseId, response.ResponseId, testCase.description, "response responseID")
			assertion.Equal(response1.QuestionnaireId, response.QuestionnaireId, testCase.description, "response questionnaireID")
			assertion.WithinDuration(response1.SubmittedAt.UTC().Truncate(time.Second), response.SubmittedAt.UTC(), time.Second, testCase.description, "response submittedAt")
		}
	}
}

func TestDeleteResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		invalidResponseID bool
	}
	type expect struct {
		isErr bool
		err   error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "valid",
			args: args{
				invalidResponseID: false,
			},
		},
		{
			description: "invalid",
			args: args{
				invalidResponseID: true,
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		questionnaire := sampleQuestionnaire
		e := echo.New()
		body, err := json.Marshal(questionnaire)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
		require.NoError(t, err)
		var responseID int
		if !testCase.args.invalidResponseID {
			AddQuestionID2SampleResponseMutex.Lock()

			AddQuestionID2SampleResponse(questionnaireDetail.QuestionnaireId)
			newResponse := sampleResponse
			e := echo.New()
			body, err := json.Marshal(questionnaire)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			response, err := q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, "userOne")
			require.NoError(t, err)

			AddQuestionID2SampleResponseMutex.Unlock()

			responseID = response.ResponseId
		} else {
			responseID = 10000
			valid := true
			for valid {
				c := context.Background()
				_, err := IRespondent.GetRespondent(c, responseID)
				if err == model.ErrRecordNotFound {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting respondent")
				} else {
					responseID *= 10
				}
			}
		}

		e = echo.New()
		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/responses/%d", responseID), nil)
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		err = r.DeleteResponse(ctx, responseID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}

		c := context.Background()
		_, err = IRespondent.GetRespondent(c, responseID)

		if err == nil {
			assertion.Fail("response not deleted")
		} else if err != model.ErrRecordNotFound {
			assertion.Fail("unexpected error during getting respondent")
		}
	}
}

func TestEditResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	responseDueDateTimePlus := time.Now().Add(24 * time.Hour)

	questionnaire := sampleQuestionnaire
	questionnaire.ResponseDueDateTime = &responseDueDateTimePlus
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.ResponseDueDateTime = &responseDueDateTimePlus
	questionnaire.IsAnonymous = true
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireDetailAnonymous, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireDetailNoDue, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	type args struct {
		invalidResponseID   bool
		questionnaireDetail openapi.QuestionnaireDetail
		isAnonymous         bool
		params              openapi.PostQuestionnaireResponseJSONRequestBody
		userID              string
		isTimeAfterDue      bool
	}
	type expect struct {
		isErr bool
		err   error
	}
	type test struct {
		description string
		args
		expect
	}

	invalidResponseBodyText := openapi.ResponseBody{}
	invalidResponseBodyText.QuestionId = *questionnaireDetail.Questions[0].QuestionId
	invalidResponseBodyText.FromResponseBodyText(openapi.ResponseBodyText{
		Answer:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		QuestionType: "Text",
	})
	invalidResponseBodyTextLong := openapi.ResponseBody{}
	invalidResponseBodyTextLong.QuestionId = *questionnaireDetail.Questions[1].QuestionId
	invalidResponseBodyTextLong.FromResponseBodyTextLong(openapi.ResponseBodyTextLong{
		Answer:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		QuestionType: "TextLong",
	})
	invalidResponseBodyNumber := openapi.ResponseBody{}
	invalidResponseBodyNumber.QuestionId = *questionnaireDetail.Questions[2].QuestionId
	invalidResponseBodyNumber.FromResponseBodyNumber(openapi.ResponseBodyNumber{
		Answer:       101,
		QuestionType: "Number",
	})
	invalidResponseBodySingleChoice := openapi.ResponseBody{}
	invalidResponseBodySingleChoice.QuestionId = *questionnaireDetail.Questions[3].QuestionId
	invalidResponseBodySingleChoice.FromResponseBodySingleChoice(openapi.ResponseBodySingleChoice{
		Answer:       5,
		QuestionType: "SingleChoice",
	})
	invalidResponseBodyMultipleChoice := openapi.ResponseBody{}
	invalidResponseBodyMultipleChoice.QuestionId = *questionnaireDetail.Questions[4].QuestionId
	invalidResponseBodyMultipleChoice.FromResponseBodyMultipleChoice(openapi.ResponseBodyMultipleChoice{
		Answer:       []int{5},
		QuestionType: "MultipleChoice",
	})
	invalidResponseBodyScale := openapi.ResponseBody{}
	invalidResponseBodyScale.QuestionId = *questionnaireDetail.Questions[5].QuestionId
	invalidResponseBodyScale.FromResponseBodyScale(openapi.ResponseBodyScale{
		Answer:       0,
		QuestionType: "Scale",
	})

	AddQuestionID2SampleResponseMutex.Lock()

	AddQuestionID2SampleResponse(questionnaireDetail.QuestionnaireId)
	testCases := []test{
		{
			description: "valid",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params:              sampleResponse,
				userID:              userOne,
			},
		},
		{
			description: "invalid edit to draft",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: true,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response id",
			args: args{
				questionnaireDetail: questionnaireDetail,
				invalidResponseID:   true,
				params:              sampleResponse,
				userID:              userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body text",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						invalidResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body text long",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						sampleResponseBodyText,
						invalidResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body number",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						invalidResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body single choice",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						invalidResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body multiple choice",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						invalidResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body scale",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.ResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						invalidResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
	}
	tmp_testCase := test{
		description: "question type not match",
		args: args{
			questionnaireDetail: questionnaireDetail,
			params: openapi.PostQuestionnaireResponseJSONRequestBody{
				Body: []openapi.ResponseBody{
					sampleResponseBodyText,
					sampleResponseBodyScale,
					sampleResponseBodyNumber,
					sampleResponseBodySingleChoice,
					sampleResponseBodyMultipleChoice,
					sampleResponseBodyTextLong,
				},
				IsDraft: false,
			},
			userID: userOne,
		},
		expect: expect{
			isErr: true,
		},
	}
	tmp_testCase.args.params.Body[1].QuestionId = sampleResponseBodyTextLong.QuestionId
	tmp_testCase.args.params.Body[5].QuestionId = sampleResponseBodyScale.QuestionId
	testCases = append(testCases, tmp_testCase)

	AddQuestionID2SampleResponse(questionnaireDetailAnonymous.QuestionnaireId)
	testCases = append(testCases, test{
		description: "valid anonymous",
		args: args{
			questionnaireDetail: questionnaireDetailAnonymous,
			isAnonymous:         true,
			params:              sampleResponse,
			userID:              userOne,
		},
	})

	AddQuestionID2SampleResponse(questionnaireDetailNoDue.QuestionnaireId)
	testCases = append(testCases, test{
		description: "questionnaire no due",
		args: args{
			questionnaireDetail: questionnaireDetailNoDue,
			params:              sampleResponse,
			userID:              userOne,
		},
	})

	testCases = append(testCases, test{
		description: "is time after due (attention: need to be the last testCase)",
		args: args{
			questionnaireDetail: questionnaireDetail,
			params:              sampleResponse,
			userID:              userOne,
			isTimeAfterDue:      true,
		},
		expect: expect{
			isErr: true,
		},
	})

	AddQuestionID2SampleResponseMutex.Unlock()

	for _, testCase := range testCases {
		AddQuestionID2SampleResponseMutex.Lock()

		AddQuestionID2SampleResponse(testCase.args.questionnaireDetail.QuestionnaireId)
		e = echo.New()
		body, err = json.Marshal(sampleResponse)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", testCase.args.questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		response, err := q.PostQuestionnaireResponse(ctx, testCase.args.questionnaireDetail.QuestionnaireId, sampleResponse, testCase.args.userID)
		require.NoError(t, err)

		AddQuestionID2SampleResponseMutex.Unlock()

		var responseID int
		if !testCase.args.invalidResponseID {
			responseID = response.ResponseId
		} else {
			responseID = 10000
			valid := true
			for valid {
				c := context.Background()
				_, err := IRespondent.GetRespondent(c, responseID)
				if err == model.ErrRecordNotFound {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting respondent")
				} else {
					responseID *= 10
				}
			}
		}

		if testCase.args.isTimeAfterDue {
			c := context.Background()
			responseDueDateTime := null.Time{}
			responseDueDateTime.Valid = true
			responseDueDateTime.Time = time.Now().Add(-24 * time.Hour)
			err = IQuestionnaire.UpdateQuestionnaire(c, testCase.args.questionnaireDetail.Title, testCase.args.questionnaireDetail.Description, responseDueDateTime, convertResponseViewableBy(testCase.args.questionnaireDetail.ResponseViewableBy), testCase.args.questionnaireDetail.QuestionnaireId, testCase.args.questionnaireDetail.IsPublished, testCase.args.questionnaireDetail.IsAnonymous, testCase.args.questionnaireDetail.IsDuplicateAnswerAllowed)
			require.NoError(t, err)
		}

		var responseEditPost openapi.EditResponseJSONRequestBody
		responseEditPost.Body = testCase.args.params.Body
		responseEditPost.IsDraft = testCase.args.params.IsDraft
		responseEditPost.ResponseId = &responseID

		e = echo.New()
		body, err = json.Marshal(responseEditPost)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/responses/%d", responseID), bytes.NewReader(body))
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		err = r.EditResponse(ctx, responseID, responseEditPost)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/responses/%d", responseID), nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		responseEdited, err := r.GetResponse(ctx, responseID)
		require.NoError(t, err)

		assertion.Equal(response.QuestionnaireId, responseEdited.QuestionnaireId, testCase.description, "questionnaireId")
		if testCase.args.isAnonymous {
			assertion.Nil(responseEdited.Respondent, testCase.description, "anonymous response with respondent")
		} else {
			assertion.Equal(response.Respondent, responseEdited.Respondent, testCase.description, "respondent")
		}
		assertion.Equal(response.ResponseId, responseEdited.ResponseId, testCase.description, "responseId")
		assertion.Equal(response.IsAnonymous, responseEdited.IsAnonymous, testCase.description, "isAnonymous")
		assertion.WithinDuration(response.SubmittedAt.UTC().Truncate(time.Second), responseEdited.SubmittedAt.UTC(), time.Second, testCase.description, "submittedAt")
		modifiedAtDiff := time.Since(responseEdited.ModifiedAt)
		assertion.True(modifiedAtDiff > -time.Second && modifiedAtDiff < time.Minute, testCase.description, "modifiedAt", responseEdited.ModifiedAt)

		assertion.Equal(responseEditPost.Body, responseEdited.Body, testCase.description, "response body")
		assertion.Equal(responseEditPost.IsDraft, responseEdited.IsDraft, testCase.description, "response isDraft")

	}
}
