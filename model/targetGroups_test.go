package model

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestInsertTargetGroups(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		targetGroups []uuid.UUID
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
			description: "no group",
			args: args{
				targetGroups: []uuid.UUID{},
			},
		},
		{
			description: "one group",
			args: args{
				targetGroups: []uuid.UUID{groupOne},
			},
		},
		{
			description: "two group",
			args: args{
				targetGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = targetGroupImpl.InsertTargetGroups(ctx, questionnaireID, testCase.args.targetGroups)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualTargetGroups []TargetGroups
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualTargetGroups).Error
		require.NoError(t, err)

		actualTargetGroupIDs := make([]uuid.UUID, len(actualTargetGroups))
		for i, targetGroup := range actualTargetGroups {
			actualTargetGroupIDs[i] = targetGroup.GroupID
		}

		sort.Slice(testCase.args.targetGroups, func(i, j int) bool {
			return testCase.args.targetGroups[i].String() < testCase.args.targetGroups[j].String()
		})
		sort.Slice(actualTargetGroupIDs, func(i, j int) bool { return actualTargetGroupIDs[i].String() < actualTargetGroupIDs[j].String() })
		assertion.Equal(testCase.args.targetGroups, actualTargetGroupIDs, testCase.description, "target groups")
	}
}

func TestDeleteTargetGroups(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		targetGroups []uuid.UUID
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
			description: "no group",
			args: args{
				targetGroups: []uuid.UUID{},
			},
		},
		{
			description: "one group",
			args: args{
				targetGroups: []uuid.UUID{groupOne},
			},
		},
		{
			description: "two group",
			args: args{
				targetGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = targetGroupImpl.InsertTargetGroups(ctx, questionnaireID, testCase.args.targetGroups)
		require.NoError(t, err)

		err = targetGroupImpl.DeleteTargetGroups(ctx, questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualTargetGroups []TargetGroups
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualTargetGroups).Error
		require.NoError(t, err)

		assertion.Equal([]TargetGroups{}, actualTargetGroups, testCase.description, "target groups")
	}
}

func TestGetTargetGroups(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		dualQuestionnaire bool
		targetGroups      []uuid.UUID
		targetGroups2     []uuid.UUID
	}
	type expect struct {
		targetGroups []uuid.UUID
		isErr        bool
		err          error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "no group",
			args: args{
				targetGroups: []uuid.UUID{},
			},
			expect: expect{
				targetGroups: []uuid.UUID{},
			},
		},
		{
			description: "one group",
			args: args{
				targetGroups: []uuid.UUID{groupOne},
			},
			expect: expect{
				targetGroups: []uuid.UUID{groupOne},
			},
		},
		{
			description: "two group",
			args: args{
				targetGroups: []uuid.UUID{groupOne, groupTwo},
			},
			expect: expect{
				targetGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
		{
			description: "two group",
			args: args{
				targetGroups: []uuid.UUID{groupOne, groupTwo},
			},
			expect: expect{
				targetGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
		{
			description: "dual questionnaire",
			args: args{
				dualQuestionnaire: true,
				targetGroups:      []uuid.UUID{groupOne, groupTwo},
				targetGroups2:     []uuid.UUID{groupTwo, groupThree},
			},
			expect: expect{
				targetGroups: []uuid.UUID{groupOne, groupTwo, groupTwo, groupThree},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)
		err = targetGroupImpl.InsertTargetGroups(ctx, questionnaireID, testCase.args.targetGroups)
		require.NoError(t, err)

		var questionnaireID2 int
		if testCase.dualQuestionnaire {
			questionnaireID2, err = questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
			require.NoError(t, err)
			err = targetGroupImpl.InsertTargetGroups(ctx, questionnaireID2, testCase.args.targetGroups2)
			require.NoError(t, err)
		}

		var actualTargetGroups []TargetGroups
		if !testCase.dualQuestionnaire {
			actualTargetGroups, err = targetGroupImpl.GetTargetGroups(ctx, []int{questionnaireID})
		} else {
			actualTargetGroups, err = targetGroupImpl.GetTargetGroups(ctx, []int{questionnaireID, questionnaireID2})
		}
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}

		actualTargetGroupIDs := make([]uuid.UUID, len(actualTargetGroups))
		for i, targetGroup := range actualTargetGroups {
			actualTargetGroupIDs[i] = targetGroup.GroupID
		}

		sort.Slice(testCase.expect.targetGroups, func(i, j int) bool {
			return testCase.expect.targetGroups[i].String() < testCase.expect.targetGroups[j].String()
		})
		sort.Slice(actualTargetGroupIDs, func(i, j int) bool { return actualTargetGroupIDs[i].String() < actualTargetGroupIDs[j].String() })
		assertion.Equal(testCase.expect.targetGroups, actualTargetGroupIDs, testCase.description, "target groups")
	}
}
