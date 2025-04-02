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

func TestInsertAdministratorGroups(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		adminGroups []uuid.UUID
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
				adminGroups: []uuid.UUID{},
			},
		},
		{
			description: "one group",
			args: args{
				adminGroups: []uuid.UUID{groupOne},
			},
		},
		{
			description: "two group",
			args: args{
				adminGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = administratorGroupImpl.InsertAdministratorGroups(ctx, questionnaireID, testCase.args.adminGroups)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualAdministratorGroups []AdministratorGroups
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualAdministratorGroups).Error
		require.NoError(t, err)

		actualAdminGroupIDs := make([]uuid.UUID, len(actualAdministratorGroups))
		for i, adminGroup := range actualAdministratorGroups {
			actualAdminGroupIDs[i] = adminGroup.GroupID
		}

		sort.Slice(testCase.args.adminGroups, func(i, j int) bool {
			return testCase.args.adminGroups[i].String() < testCase.args.adminGroups[j].String()
		})
		sort.Slice(actualAdminGroupIDs, func(i, j int) bool { return actualAdminGroupIDs[i].String() < actualAdminGroupIDs[j].String() })
		assertion.Equal(testCase.args.adminGroups, actualAdminGroupIDs, testCase.description, "admin groups")
	}
}

func TestDeleteAdministratorGroups(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		adminGroups []uuid.UUID
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
				adminGroups: []uuid.UUID{},
			},
		},
		{
			description: "one group",
			args: args{
				adminGroups: []uuid.UUID{groupOne},
			},
		},
		{
			description: "two group",
			args: args{
				adminGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = administratorGroupImpl.InsertAdministratorGroups(ctx, questionnaireID, testCase.args.adminGroups)
		require.NoError(t, err)

		err = administratorGroupImpl.DeleteAdministratorGroups(ctx, questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualAdministratorGroups []AdministratorGroups
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualAdministratorGroups).Error
		require.NoError(t, err)

		assertion.Equal([]AdministratorGroups{}, actualAdministratorGroups, testCase.description, "admin groups")
	}
}

func TestGetAdministratorGroups(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		dualQuestionnaire bool
		adminGroups       []uuid.UUID
		adminGroups2      []uuid.UUID
	}
	type expect struct {
		adminGroups []uuid.UUID
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
			description: "no group",
			args: args{
				adminGroups: []uuid.UUID{},
			},
			expect: expect{
				adminGroups: []uuid.UUID{},
			},
		},
		{
			description: "one group",
			args: args{
				adminGroups: []uuid.UUID{groupOne},
			},
			expect: expect{
				adminGroups: []uuid.UUID{groupOne},
			},
		},
		{
			description: "two group",
			args: args{
				adminGroups: []uuid.UUID{groupOne, groupTwo},
			},
			expect: expect{
				adminGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
		{
			description: "two group",
			args: args{
				adminGroups: []uuid.UUID{groupOne, groupTwo},
			},
			expect: expect{
				adminGroups: []uuid.UUID{groupOne, groupTwo},
			},
		},
		{
			description: "dual questionnaire",
			args: args{
				dualQuestionnaire: true,
				adminGroups:       []uuid.UUID{groupOne, groupTwo},
				adminGroups2:      []uuid.UUID{groupTwo, groupThree},
			},
			expect: expect{
				adminGroups: []uuid.UUID{groupOne, groupTwo, groupTwo, groupThree},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)
		err = administratorGroupImpl.InsertAdministratorGroups(ctx, questionnaireID, testCase.args.adminGroups)
		require.NoError(t, err)

		var questionnaireID2 int
		if testCase.dualQuestionnaire {
			questionnaireID2, err = questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
			require.NoError(t, err)
			err = administratorGroupImpl.InsertAdministratorGroups(ctx, questionnaireID2, testCase.args.adminGroups2)
			require.NoError(t, err)
		}

		var actualAdministratorGroups []AdministratorGroups
		if !testCase.dualQuestionnaire {
			actualAdministratorGroups, err = administratorGroupImpl.GetAdministratorGroups(ctx, []int{questionnaireID})
		} else {
			actualAdministratorGroups, err = administratorGroupImpl.GetAdministratorGroups(ctx, []int{questionnaireID, questionnaireID2})
		}
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}

		actualAdminGroupIDs := make([]uuid.UUID, len(actualAdministratorGroups))
		for i, adminGroup := range actualAdministratorGroups {
			actualAdminGroupIDs[i] = adminGroup.GroupID
		}

		sort.Slice(testCase.expect.adminGroups, func(i, j int) bool {
			return testCase.expect.adminGroups[i].String() < testCase.expect.adminGroups[j].String()
		})
		sort.Slice(actualAdminGroupIDs, func(i, j int) bool { return actualAdminGroupIDs[i].String() < actualAdminGroupIDs[j].String() })
		assertion.Equal(testCase.expect.adminGroups, actualAdminGroupIDs, testCase.description, "admin groups")
	}
}
