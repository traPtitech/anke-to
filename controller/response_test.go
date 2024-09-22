package controller

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/openapi"
)

func TestGetMyResponses(t *testing.T) {
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
		params openapi.GetMyResponsesParams
		userID string
	}
	type expect struct {
		isErr                          bool
		err                            error
		responsesWithQuestionnaireInfo openapi.ResponsesWithQuestionnaireInfo
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

		responsesWithQuestionnaireInfo, err := r.GetMyResponses(ctx, testCase.args.params, testCase.args.userID)

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

		assertion.Equal(testCase.expect.responsesWithQuestionnaireInfo, responsesWithQuestionnaireInfo, testCase.description, "responsesWithQuestionnaireInfo")
	}
}

func TestGetResponse(t *testing.T) {
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
	}
	type expect struct {
		isErr    bool
		err      error
		response openapi.Response
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

		response, err := r.GetResponse(ctx, testCase.args.responseID)

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

		assertion.Equal(testCase.expect.response, response, testCase.description, "response")
	}
}

func TestDeleteResponse(t *testing.T) {
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
		userID     string
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

		err := r.DeleteResponse(ctx, testCase.args.responseID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
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
