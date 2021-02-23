package model

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestInsertResponses(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	questionID, err := questionImpl.InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
	require.NoError(t, err)

	type args struct {
		validID       bool
		responseMetas []*ResponseMeta
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
				validID: true,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionID, Data: "リマインダーBOTを作った話"},
				},
			},
		},
		{
			description: "long Data",
			args: args{
				validID: true,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionID, Data: strings.Repeat("a", 200)},
				},
			},
		},
		{
			description: "too long Data",
			args: args{
				validID: true,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionID, Data: strings.Repeat("a", 200000)},
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "questionID not exist",
			args: args{
				validID: true,
				responseMetas: []*ResponseMeta{
					{QuestionID: -1, Data: "リマインダーBOTを作った話"},
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "responseID not exist",
			args: args{
				validID: false,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionID, Data: "リマインダーBOTを作った話"},
				},
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		responseID, err := respondentImpl.InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)
		if !testCase.args.validID {
			responseID = -1
		}
		err = responseImpl.InsertResponses(responseID, testCase.args.responseMetas)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		response := Responses{}
		err = db.Where("response_id = ?", responseID).First(&response).Error
		if err != nil {
			t.Errorf("failed to get questionnaire(%s): %w", testCase.description, err)
		}

		assertion.Equal(responseID, response.ResponseID, testCase.description, "responseID")
		assertion.Equal(questionID, response.QuestionID, testCase.description, "questionID")
		assertion.Equal(testCase.args.responseMetas[0].Data, response.Body.ValueOrZero(), testCase.description, "Body")
		assertion.WithinDuration(time.Now(), response.ModifiedAt, 2*time.Second, testCase.description, "ModifiedAt")
		assertion.Equal(time.Time{}, response.DeletedAt.ValueOrZero(), 2*time.Second, testCase.description, "DeletedAt")
	}
}

func TestDeleteResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	questionID, err := questionImpl.InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
	require.NoError(t, err)

	type args struct {
		validID       bool
		responseMetas []*ResponseMeta
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
				validID: true,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionID, Data: "リマインダーBOTを作った話"},
				},
			},
		},
		{
			description: "responseID not exist",
			args: args{
				validID: false,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionID, Data: "リマインダーBOTを作った話"},
				},
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordDeleted,
			},
		},
	}

	for _, testCase := range testCases {
		responseID, err := respondentImpl.InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)
		err = responseImpl.InsertResponses(responseID, testCase.args.responseMetas)
		require.NoError(t, err)
		if !testCase.args.validID {
			responseID = -1
		}

		err = responseImpl.DeleteResponse(responseID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		response := Responses{}
		err = db.
			Unscoped().
			Where("response_id = ?", responseID).
			First(&response).Error
		if err != nil {
			t.Errorf("failed to get responses(%s): %w", testCase.description, err)
		}

		assertion.WithinDuration(time.Now(), response.DeletedAt.ValueOrZero(), 2*time.Second)
	}
}
