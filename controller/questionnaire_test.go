package controller

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/openapi"
)

func TestGetQuestionnaires(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaires := []struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	questionnaireIDs := []int{}
	for _, questionnaire := range questionnaires {
		// ctxの作成
		questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnairePosted.QuestionnaireId)
	}

	type args struct {
		userID string
		params openapi.GetQuestionnairesParams
	}
	type expect struct {
		isErr             bool
		err               error
		questionnaireList openapi.QuestionnaireList
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

		questionnaireList, err := q.GetQuestionnaires(ctx, testCase.args.userID, testCase.args.params)

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

		assertion.Equal(testCase.expect.questionnaireList, questionnaireList, testCase.description, "questionnaireList")
	}
}

func TestPostQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	type args struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}
	type expect struct {
		isErr               bool
		err                 error
		questionnaireDetail openapi.QuestionnaireDetail
	}
	type test struct {
		description string
		args
		expect
	}

	// ResponseDueDateTime := null.NewTime(time.Time{}, false)
	// ResponseDueDateTimeMinus := null.NewTime(time.Now().Add(-24*time.Hour), true)
	// ResponseDueDateTimePlus := null.NewTime(time.Now().Add(24*time.Hour), true)

	testCases := []test{
		// todo: テストケースを追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		questionnaireDetail, err := q.PostQuestionnaire(ctx, testCase.args.userID, testCase.args.params)

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

		assertion.Equal(testCase.expect.questionnaireDetail, questionnaireDetail, testCase.description, "questionnaireDetail")
	}
}

func TestGetQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaires := []struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	questionnaireIDs := []int{}
	for _, questionnaire := range questionnaires {
		// ctxの作成
		questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnairePosted.QuestionnaireId)
	}

	type args struct {
		questionnaireID int
	}
	type expect struct {
		isErr               bool
		err                 error
		questionnaireDetail openapi.QuestionnaireDetail
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

		questionnaireDetail, err := q.GetQuestionnaire(ctx, testCase.args.questionnaireID)

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

		assertion.Equal(testCase.expect.questionnaireDetail, questionnaireDetail, testCase.description, "questionnaireDetail")
	}
}

func TestEditQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaires := []struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	questionnaireIDs := []int{}
	for _, questionnaire := range questionnaires {
		// ctxの作成
		questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnairePosted.QuestionnaireId)
	}

	type args struct {
		questionnaireID int
		params          openapi.EditQuestionnaireJSONRequestBody
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

		err := q.EditQuestionnaire(ctx, testCase.args.questionnaireID, testCase.args.params)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
	}
}

func TestDeleteQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaires := []struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	questionnaireIDs := []int{}
	for _, questionnaire := range questionnaires {
		// ctxの作成
		questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnairePosted.QuestionnaireId)
	}

	type args struct {
		questionnaireID int
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

		err := q.DeleteQuestionnaire(ctx, testCase.args.questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
	}
}

func TestGetQuestionnaireMyRemindStatus(t *testing.T) {
	// todo
}

func TestEditQuestionnaireMyRemindStatus(t *testing.T) {
	// todo
}

func TestGetQuestionnaireResponses(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

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

	for _, response := range responses {
		// todo: ctxの作成
		_, err := q.PostQuestionnaireResponse(ctx, response.questionnaireID, response.params, response.userID)
		require.NoError(t, err)
	}

	type args struct {
		questionnaireID int
		params          openapi.GetQuestionnaireResponsesParams
		userID          string
	}
	type expect struct {
		isErr     bool
		err       error
		responses openapi.Responses
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

		responses, err := q.GetQuestionnaireResponses(ctx, testCase.args.questionnaireID, testCase.args.params, testCase.args.userID)

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

		assertion.Equal(testCase.expect.responses, responses, testCase.description, "responses")
	}
}

func TestPostQuestionnaireResponse(t *testing.T) {

	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	type args struct {
		questionnaireID int
		params          openapi.PostQuestionnaireResponseJSONRequestBody
		userID          string
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
		//	todo: テストケースの追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		response, err := q.PostQuestionnaireResponse(ctx, testCase.args.questionnaireID, testCase.args.params, testCase.userID)

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

func TestGetQuestionnaireResult(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

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

	for _, response := range responses {
		// todo: ctxの作成
		_, err := q.PostQuestionnaireResponse(ctx, response.questionnaireID, response.params, response.userID)
		require.NoError(t, err)
	}

	type args struct {
		questionnaireID int
		userID          string
	}
	type expect struct {
		isErr  bool
		err    error
		result openapi.Result
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

		result, err := q.GetQuestionnaireResult(ctx, testCase.args.questionnaireID, testCase.args.userID)

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

		assertion.Equal(testCase.expect.result, result, testCase.description, "result")
	}
}
