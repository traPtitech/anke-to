// todo: add testcase for response_vieable_by & is_published
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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/model/mock_model"
	"github.com/traPtitech/anke-to/openapi"
	"github.com/traPtitech/anke-to/traq/mock_traq"
)

var (
	sampleResponseBodyText           = openapi.ResponseBody{}
	sampleResponseBodyTextLong       = openapi.ResponseBody{}
	sampleResponseBodyNumber         = openapi.ResponseBody{}
	sampleResponseBodySingleChoice   = openapi.ResponseBody{}
	sampleResponseBodyMultipleChoice = openapi.ResponseBody{}
	sampleResponseBodyScale          = openapi.ResponseBody{}
	sampleResponse                   = openapi.NewResponse{}
)

func setupSampleResponse() {
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
		Answer:       0,
		QuestionType: "SingleChoice",
	})
	sampleResponseBodyMultipleChoice.FromResponseBodyMultipleChoice(openapi.ResponseBodyMultipleChoice{
		Answer:       []int{0, 1},
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

func TestGetMyResponses(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)

	mockTargetGroup := mock_model.NewMockITargetGroup(ctrl)
	mockTargetUser := mock_model.NewMockITargetUser(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockAdministratorGroup := mock_model.NewMockIAdministratorGroup(ctrl)
	mockAdministratorUser := mock_model.NewMockIAdministratorUser(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockTransaction := mock_model.NewMockITransaction(ctrl)
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	r := NewResponse(mockQuestionnaire, mockRespondent, mockResponse, mockTarget, mockQuestion, mockValidation, mockScaleLabel)
	q := NewQuestionnaire(mockQuestionnaire, mockTarget, mockTargetGroup, mockTargetUser, mockAdministrator, mockAdministratorGroup, mockAdministratorUser, mockQuestion, mockOption, mockScaleLabel, mockValidation, mockTransaction, mockRespondent, mockWebhook, r)

	setupSampleQuestionnaire()
	setupSampleResponse()

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
	sortSubmittedAt := (openapi.ResponseSortInQuery)("submitted_at")
	sortSubmittedAtDesc := (openapi.ResponseSortInQuery)("-submitted_at")
	// sortTitle := (openapi.ResponseSortInQuery)("title")
	// sortTitleDesc := (openapi.ResponseSortInQuery)("-title")
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
			description: "sort created_at",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortSubmittedAt,
				},
			},
		},
		{
			description: "sort -created_at",
			args: args{
				userID: userOne,
				params: openapi.GetMyResponsesParams{
					Sort: &sortSubmittedAtDesc,
				},
			},
		},
		// {
		// 	description: "sort title",
		// 	args: args{
		// 		userID:          userOne,
		// 		params: openapi.GetQuestionnaireResponsesParams{
		// 			Sort: &sortTitle,
		// 		},
		// 	},
		// },
		// {
		// 	description: "sort -title",
		// 	args: args{
		// 		userID:          userOne,
		// 		params: openapi.GetQuestionnaireResponsesParams{
		// 			Sort: &sortTitleDesc,
		// 		},
		// 	},
		// },
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

		responseList, err := r.GetMyResponses(ctx, testCase.args.params, testCase.args.userID)

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

		if testCase.args.params.Sort != nil {
			if *testCase.args.params.Sort == "submitted_at" {
				var preCreatedAt time.Time
				for _, response := range responseList {
					if !preCreatedAt.IsZero() {
						assertion.True(preCreatedAt.Before(response.SubmittedAt), testCase.description, "created_at")
					}
					preCreatedAt = response.SubmittedAt
				}
			} else if *testCase.args.params.Sort == "-submitted_at" {
				var preCreatedAt time.Time
				for _, response := range responseList {
					if !preCreatedAt.IsZero() {
						assertion.True(preCreatedAt.After(response.SubmittedAt), testCase.description, "-created_at")
					}
					preCreatedAt = response.SubmittedAt
				}
				// } else if *testCase.args.params.Sort == "title" {
				// 	var preTitle string
				// 	for _, response := range responseList {
				// 		if preTitle != "" {
				// 			assertion.True(preTitle > response.Title, testCase.description, "title")
				// 		}
				// 		preTitle = response.Title
				// 	}
				// } else if *testCase.args.params.Sort == "-title" {
				// 	var preTitle string
				// 	for _, response := range responseList {
				// 		if preTitle != "" {
				// 			assertion.True(preTitle < response.Title, testCase.description, "-title")
				// 		}
				// 		preTitle = response.Title
				// 	}
			} else if *testCase.args.params.Sort == "modified_at" {
				var preModifiedAt time.Time
				for _, response := range responseList {
					if !preModifiedAt.IsZero() {
						assertion.True(preModifiedAt.Before(response.ModifiedAt), testCase.description, "modified_at")
					}
					preModifiedAt = response.ModifiedAt
				}
			} else if *testCase.args.params.Sort == "-modified_at" {
				var preModifiedAt time.Time
				for _, response := range responseList {
					if !preModifiedAt.IsZero() {
						assertion.True(preModifiedAt.After(response.ModifiedAt), testCase.description, "-modified_at")
					}
					preModifiedAt = response.ModifiedAt
				}
			}
		}

		if testCase.expect.responseIdList != nil {
			var responseIdList []int
			for _, response := range responseList {
				responseIdList = append(responseIdList, response.ResponseId)
			}
			sort.Slice(*testCase.expect.responseIdList, func(i, j int) bool {
				return (*testCase.expect.responseIdList)[i] < (*testCase.expect.responseIdList)[j]
			})
			sort.Slice(responseIdList, func(i, j int) bool { return responseIdList[i] < responseIdList[j] })
			assertion.Equal(testCase.expect.responseIdList, responseIdList, testCase.description, "responseIdList")
		}

		for _, response := range responseList {
			if *response.IsAnonymous {
				assertion.Equal(response.Respondent, nil, testCase.description, "anonymous response with respondent")
			} else {
				assertion.NotEqual(response.Respondent, nil, testCase.description, "response with no respondent")
			}
		}
	}
}

func TestGetResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)

	mockTargetGroup := mock_model.NewMockITargetGroup(ctrl)
	mockTargetUser := mock_model.NewMockITargetUser(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockAdministratorGroup := mock_model.NewMockIAdministratorGroup(ctrl)
	mockAdministratorUser := mock_model.NewMockIAdministratorUser(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockTransaction := mock_model.NewMockITransaction(ctrl)
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	r := NewResponse(mockQuestionnaire, mockRespondent, mockResponse, mockTarget, mockQuestion, mockValidation, mockScaleLabel)
	q := NewQuestionnaire(mockQuestionnaire, mockTarget, mockTargetGroup, mockTargetUser, mockAdministrator, mockAdministratorGroup, mockAdministratorUser, mockQuestion, mockOption, mockScaleLabel, mockValidation, mockTransaction, mockRespondent, mockWebhook, r)

	setupSampleQuestionnaire()
	setupSampleResponse()

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
				_, err := mockRespondent.GetRespondent(c, responseID)
				if err == errors.New("record not found") {
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
		response, err := q.GetResponse(ctx, responseID)

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
			assertion.Equal(response0, response, testCase.description, "response")
		} else {
			assertion.Equal(response1, response, testCase.description, "response")
		}

		if testCase.args.isAnonymousQuestionnaire {
			assertion.Equal(response.Respondent, nil, testCase.description, "anonymous questionnaire with respondent")
		}
	}
}

func TestDeleteResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)

	mockTargetGroup := mock_model.NewMockITargetGroup(ctrl)
	mockTargetUser := mock_model.NewMockITargetUser(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockAdministratorGroup := mock_model.NewMockIAdministratorGroup(ctrl)
	mockAdministratorUser := mock_model.NewMockIAdministratorUser(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockTransaction := mock_model.NewMockITransaction(ctrl)
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	r := NewResponse(mockQuestionnaire, mockRespondent, mockResponse, mockTarget, mockQuestion, mockValidation, mockScaleLabel)
	q := NewQuestionnaire(mockQuestionnaire, mockTarget, mockTargetGroup, mockTargetUser, mockAdministrator, mockAdministratorGroup, mockAdministratorUser, mockQuestion, mockOption, mockScaleLabel, mockValidation, mockTransaction, mockRespondent, mockWebhook, r)

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

	setupSampleQuestionnaire()
	setupSampleResponse()

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
			responseID = response.ResponseId
		} else {
			responseID = 10000
			valid := true
			for valid {
				c := context.Background()
				_, err := mockRespondent.GetRespondent(c, responseID)
				if err == errors.New("record not found") {
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
		_, err = mockRespondent.GetRespondent(c, responseID)

		if err == nil {
			assertion.Fail("response not deleted")
		} else if err != errors.New("record not found") {
			assertion.Fail("unexpected error during getting respondent")
		}
	}
}

func TestEditResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()
	r := NewResponse()

	questionnaire := struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	// ctxの作成

	questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
	require.NoError(t, err)

	questionnaireID := questionnairePosted.QuestionnaireId

	responses := []struct {
		questionnaireID int
		params          openapi.PostQuestionnaireResponseJSONRequestBody
		userID          string
	}{
		// todo
	}

	responseIDs := []int{}
	for _, response := range responses {
		// todo: ctxの作成
		responsePosted, err := q.PostQuestionnaireResponse(ctx, response.questionnaireID, response.params, response.userID)
		require.NoError(t, err)
		responseIDs = append(responseIDs, responsePosted.ResponseId)
	}

	type args struct {
		responseID openapi.ResponseIDInPath
		req        openapi.EditResponseJSONRequestBody
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
		// todo: テストケースを追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		err := r.EditResponse(ctx, testCase.args.responseID, testCase.args.req)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
	}
}