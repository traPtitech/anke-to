package model

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestInsertQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		title        string
		description  string
		resTimeLimit null.Time
		resSharedTo  string
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
			description: "time limit: no, res shared to: public",
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
		{
			description: "time limit: yes, res shared to: public",
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
		},
		{
			description: "time limit: no, res shared to: respondents",
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "respondents",
			},
		},
		{
			description: "time limit: no, res shared to: administrators",
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "administrators",
			},
		},
		{
			description: "long title",
			args: args{
				title:        strings.Repeat("a", 50),
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
		{
			description: "too long title",
			args: args{
				title:        strings.Repeat("a", 500),
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "long description",
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  strings.Repeat("a", 2000),
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
		{
			description: "too long description",
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  strings.Repeat("a", 200000),
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		questionnaireID, err := InsertQuestionnaire(testCase.args.title, testCase.args.description, testCase.args.resTimeLimit, testCase.args.resSharedTo)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		questionnaire := Questionnaires{}
		err = db.Where("id = ?", questionnaireID).First(&questionnaire).Error
		if err != nil {
			t.Errorf("failed to get questionnaire(%s): %w", testCase.description, err)
		}

		assertion.Equal(testCase.args.title, questionnaire.Title, testCase.description, "title")
		assertion.Equal(testCase.args.description, questionnaire.Description, testCase.description, "description")
		assertion.WithinDuration(testCase.args.resTimeLimit.ValueOrZero(), questionnaire.ResTimeLimit.ValueOrZero(), time.Second, testCase.description, "res_time_limit")
		assertion.Equal(testCase.args.resSharedTo, questionnaire.ResSharedTo, testCase.description, "res_shared_to")

		assertion.WithinDuration(time.Now(), questionnaire.CreatedAt, time.Second, testCase.description, "created_at")
		assertion.WithinDuration(time.Now(), questionnaire.ModifiedAt, time.Second, testCase.description, "modified_at")
	}
}

func TestUpdateQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		title        string
		description  string
		resTimeLimit null.Time
		resSharedTo  string
	}
	type expect struct {
		isErr bool
		err   error
	}

	type test struct {
		description string
		before      args
		after       args
		expect
	}

	testCases := []test{
		{
			description: "update res_shared_to",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "respondents",
			},
		},
		{
			description: "update title",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第2回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
		{
			description: "update description",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第2回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
		{
			description: "update res_shared_to(res_time_limit is valid)",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "respondents",
			},
		},
		{
			description: "update title(res_time_limit is valid)",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第2回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
		},
		{
			description: "update description(res_time_limit is valid)",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第2回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
		},
		{
			description: "update res_time_limit(null->time)",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
		},
		{
			description: "update res_time_limit(time->time)",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now().Add(time.Minute), true),
				resSharedTo:  "public",
			},
		},
		{
			description: "update res_time_limit(time->null)",
			before: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Now(), true),
				resSharedTo:  "public",
			},
			after: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
	}

	for _, testCase := range testCases {
		before := &testCase.before
		questionnaire := Questionnaires{
			Title:        before.title,
			Description:  before.description,
			ResTimeLimit: before.resTimeLimit,
			ResSharedTo:  before.resSharedTo,
		}
		err := db.Create(&questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%s): %w", testCase.description, err)
		}

		createdAt := questionnaire.CreatedAt
		questionnaireID := questionnaire.ID
		after := &testCase.after
		err = UpdateQuestionnaire(after.title, after.description, after.resTimeLimit, after.resSharedTo, questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		questionnaire = Questionnaires{}
		err = db.Where("id = ?", questionnaireID).First(&questionnaire).Error
		if err != nil {
			t.Errorf("failed to get questionnaire(%s): %w", testCase.description, err)
		}

		assertion.Equal(after.title, questionnaire.Title, testCase.description, "title")
		assertion.Equal(after.description, questionnaire.Description, testCase.description, "description")
		assertion.WithinDuration(after.resTimeLimit.ValueOrZero(), questionnaire.ResTimeLimit.ValueOrZero(), time.Second, testCase.description, "res_time_limit")
		assertion.Equal(after.resSharedTo, questionnaire.ResSharedTo, testCase.description, "res_shared_to")

		assertion.WithinDuration(createdAt, questionnaire.CreatedAt, time.Second, testCase.description, "created_at")
		assertion.WithinDuration(time.Now(), questionnaire.ModifiedAt, time.Second, testCase.description, "modified_at")
	}

	invalidQuestionnaireID := 1000
	for {
		err := db.Where("id = ?", invalidQuestionnaireID).First(&Questionnaires{}).Error
		if gorm.IsRecordNotFoundError(err) {
			break
		}
		if err != nil {
			t.Errorf("failed to get questionnaire(make invalid questionnaireID): %w", err)
			break
		}

		invalidQuestionnaireID *= 10
	}

	invalidTestCases := []args{
		{
			title:        "第1回集会らん☆ぷろ募集アンケート",
			description:  "第1回集会らん☆ぷろ参加者募集",
			resTimeLimit: null.NewTime(time.Time{}, false),
			resSharedTo:  "public",
		},
		{
			title:        "第1回集会らん☆ぷろ募集アンケート",
			description:  "第1回集会らん☆ぷろ参加者募集",
			resTimeLimit: null.NewTime(time.Now(), true),
			resSharedTo:  "public",
		},
	}

	for _, arg := range invalidTestCases {
		err := UpdateQuestionnaire(arg.title, arg.description, arg.resTimeLimit, arg.resSharedTo, invalidQuestionnaireID)
		if !errors.Is(err, ErrNoRecordUpdated) {
			if err == nil {
				t.Errorf("Succeeded with invalid questionnaireID")
			} else {
				t.Errorf("failed to update questionnaire(invalid questionnireID): %w", err)
			}
		}
	}
}
