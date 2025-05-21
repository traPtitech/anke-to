package model

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestInsertTargetUsers(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		targetUsers []string
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
			description: "no user",
			args: args{
				targetUsers: []string{},
			},
		},
		{
			description: "one user",
			args: args{
				targetUsers: []string{userOne},
			},
		},
		{
			description: "two user",
			args: args{
				targetUsers: []string{userOne, userTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = targetUserImpl.InsertTargetUsers(ctx, questionnaireID, testCase.args.targetUsers)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualTargetUsers []TargetUsers
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualTargetUsers).Error
		require.NoError(t, err)

		actualTargetUserIDs := make([]string, len(actualTargetUsers))
		for i, targetUser := range actualTargetUsers {
			actualTargetUserIDs[i] = targetUser.UserTraqid
		}

		sort.Slice(testCase.args.targetUsers, func(i, j int) bool { return testCase.args.targetUsers[i] < testCase.args.targetUsers[j] })
		sort.Slice(actualTargetUserIDs, func(i, j int) bool { return actualTargetUserIDs[i] < actualTargetUserIDs[j] })
		assertion.Equal(testCase.args.targetUsers, actualTargetUserIDs, testCase.description, "target users")
	}
}

func TestDeleteTargetUsers(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		targetUsers []string
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
			description: "no user",
			args: args{
				targetUsers: []string{},
			},
		},
		{
			description: "one user",
			args: args{
				targetUsers: []string{userOne},
			},
		},
		{
			description: "two user",
			args: args{
				targetUsers: []string{userOne, userTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = targetUserImpl.InsertTargetUsers(ctx, questionnaireID, testCase.args.targetUsers)
		require.NoError(t, err)

		err = targetUserImpl.DeleteTargetUsers(ctx, questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualTargetUsers []TargetUsers
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualTargetUsers).Error
		require.NoError(t, err)

		assertion.Equal([]TargetUsers{}, actualTargetUsers, testCase.description, "target users")
	}
}

func TestGetTargetUsers(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		dualQuestionnaire bool
		targetUsers       []string
		targetUsers2      []string
	}
	type expect struct {
		targetUsers []string
		isErr       bool
		err         error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "no user",
			args: args{
				targetUsers: []string{},
			},
			expect: expect{
				targetUsers: []string{},
			},
		},
		{
			description: "one user",
			args: args{
				targetUsers: []string{userOne},
			},
			expect: expect{
				targetUsers: []string{userOne},
			},
		},
		{
			description: "two user",
			args: args{
				targetUsers: []string{userOne, userTwo},
			},
			expect: expect{
				targetUsers: []string{userOne, userTwo},
			},
		},
		{
			description: "two user",
			args: args{
				targetUsers: []string{userOne, userTwo},
			},
			expect: expect{
				targetUsers: []string{userOne, userTwo},
			},
		},
		{
			description: "dual questionnaire",
			args: args{
				dualQuestionnaire: true,
				targetUsers:       []string{userOne, userTwo},
				targetUsers2:      []string{userTwo, userThree},
			},
			expect: expect{
				targetUsers: []string{userOne, userTwo, userTwo, userThree},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)
		err = targetUserImpl.InsertTargetUsers(ctx, questionnaireID, testCase.args.targetUsers)
		require.NoError(t, err)

		var questionnaireID2 int
		if testCase.dualQuestionnaire {
			questionnaireID2, err = questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
			require.NoError(t, err)
			err = targetUserImpl.InsertTargetUsers(ctx, questionnaireID2, testCase.args.targetUsers2)
			require.NoError(t, err)
		}

		var actualTargetUsers []TargetUsers
		if !testCase.dualQuestionnaire {
			actualTargetUsers, err = targetUserImpl.GetTargetUsers(ctx, []int{questionnaireID})
		} else {
			actualTargetUsers, err = targetUserImpl.GetTargetUsers(ctx, []int{questionnaireID, questionnaireID2})
		}
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}

		actualTargetUserIDs := make([]string, len(actualTargetUsers))
		for i, targetUser := range actualTargetUsers {
			actualTargetUserIDs[i] = targetUser.UserTraqid
		}

		sort.Slice(testCase.expect.targetUsers, func(i, j int) bool { return testCase.expect.targetUsers[i] < testCase.expect.targetUsers[j] })
		sort.Slice(actualTargetUserIDs, func(i, j int) bool { return actualTargetUserIDs[i] < actualTargetUserIDs[j] })
		assertion.Equal(testCase.expect.targetUsers, actualTargetUserIDs, testCase.description, "target users")
	}
}
