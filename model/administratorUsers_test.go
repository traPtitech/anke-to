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

func TestInsertAdministratorUsers(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		adminUsers []string
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
				adminUsers: []string{},
			},
		},
		{
			description: "one user",
			args: args{
				adminUsers: []string{userOne},
			},
		},
		{
			description: "two user",
			args: args{
				adminUsers: []string{userOne, userTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = administratorUserImpl.InsertAdministratorUsers(ctx, questionnaireID, testCase.args.adminUsers)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualAdministratorUsers []AdministratorUsers
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualAdministratorUsers).Error
		require.NoError(t, err)

		actualAdminUserIDs := make([]string, len(actualAdministratorUsers))
		for i, adminUser := range actualAdministratorUsers {
			actualAdminUserIDs[i] = adminUser.UserTraqid
		}

		sort.Slice(testCase.args.adminUsers, func(i, j int) bool { return testCase.args.adminUsers[i] < testCase.args.adminUsers[j] })
		sort.Slice(actualAdminUserIDs, func(i, j int) bool { return actualAdminUserIDs[i] < actualAdminUserIDs[j] })
		assertion.Equal(testCase.args.adminUsers, actualAdminUserIDs, testCase.description, "admin users")
	}
}

func TestDeleteAdministratorUsers(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		adminUsers []string
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
				adminUsers: []string{},
			},
		},
		{
			description: "one user",
			args: args{
				adminUsers: []string{userOne},
			},
		},
		{
			description: "two user",
			args: args{
				adminUsers: []string{userOne, userTwo},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)

		err = administratorUserImpl.InsertAdministratorUsers(ctx, questionnaireID, testCase.args.adminUsers)
		require.NoError(t, err)

		err = administratorUserImpl.DeleteAdministratorUsers(ctx, questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		var actualAdministratorUsers []AdministratorUsers
		err = db.Session(&gorm.Session{NewDB: true}).Where("questionnaire_id = ?", questionnaireID).Find(&actualAdministratorUsers).Error
		require.NoError(t, err)

		assertion.Equal([]AdministratorUsers{}, actualAdministratorUsers, testCase.description, "admin users")
	}
}

func TestGetAdministratorUsers(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		dualQuestionnaire bool
		adminUsers        []string
		adminUsers2       []string
	}
	type expect struct {
		adminUsers []string
		isErr      bool
		err        error
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
				adminUsers: []string{},
			},
			expect: expect{
				adminUsers: []string{},
			},
		},
		{
			description: "one user",
			args: args{
				adminUsers: []string{userOne},
			},
			expect: expect{
				adminUsers: []string{userOne},
			},
		},
		{
			description: "two user",
			args: args{
				adminUsers: []string{userOne, userTwo},
			},
			expect: expect{
				adminUsers: []string{userOne, userTwo},
			},
		},
		{
			description: "two user",
			args: args{
				adminUsers: []string{userOne, userTwo},
			},
			expect: expect{
				adminUsers: []string{userOne, userTwo},
			},
		},
		{
			description: "dual questionnaire",
			args: args{
				dualQuestionnaire: true,
				adminUsers:        []string{userOne, userTwo},
				adminUsers2:       []string{userTwo, userThree},
			},
			expect: expect{
				adminUsers: []string{userOne, userTwo, userTwo, userThree},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
		require.NoError(t, err)
		err = administratorUserImpl.InsertAdministratorUsers(ctx, questionnaireID, testCase.args.adminUsers)
		require.NoError(t, err)

		var questionnaireID2 int
		if testCase.dualQuestionnaire {
			questionnaireID2, err = questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
			require.NoError(t, err)
			err = administratorUserImpl.InsertAdministratorUsers(ctx, questionnaireID2, testCase.args.adminUsers2)
			require.NoError(t, err)
		}

		var actualAdministratorUsers []AdministratorUsers
		if !testCase.dualQuestionnaire {
			actualAdministratorUsers, err = administratorUserImpl.GetAdministratorUsers(ctx, []int{questionnaireID})
		} else {
			actualAdministratorUsers, err = administratorUserImpl.GetAdministratorUsers(ctx, []int{questionnaireID, questionnaireID2})
		}
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}

		actualAdminUserIDs := make([]string, len(actualAdministratorUsers))
		for i, adminUser := range actualAdministratorUsers {
			actualAdminUserIDs[i] = adminUser.UserTraqid
		}

		sort.Slice(testCase.expect.adminUsers, func(i, j int) bool { return testCase.expect.adminUsers[i] < testCase.expect.adminUsers[j] })
		sort.Slice(actualAdminUserIDs, func(i, j int) bool { return actualAdminUserIDs[i] < actualAdminUserIDs[j] })
		assertion.Equal(testCase.expect.adminUsers, actualAdminUserIDs, testCase.description, "admin users")
	}
}
