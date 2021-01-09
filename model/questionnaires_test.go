package model

import (
	"errors"
	"math"
	"sort"
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
		assertion.WithinDuration(testCase.args.resTimeLimit.ValueOrZero(), questionnaire.ResTimeLimit.ValueOrZero(), 2*time.Second, testCase.description, "res_time_limit")
		assertion.Equal(testCase.args.resSharedTo, questionnaire.ResSharedTo, testCase.description, "res_shared_to")

		assertion.WithinDuration(time.Now(), questionnaire.CreatedAt, 2*time.Second, testCase.description, "created_at")
		assertion.WithinDuration(time.Now(), questionnaire.ModifiedAt, 2*time.Second, testCase.description, "modified_at")
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
		assertion.WithinDuration(after.resTimeLimit.ValueOrZero(), questionnaire.ResTimeLimit.ValueOrZero(), 2*time.Second, testCase.description, "res_time_limit")
		assertion.Equal(after.resSharedTo, questionnaire.ResSharedTo, testCase.description, "res_shared_to")

		assertion.WithinDuration(createdAt, questionnaire.CreatedAt, 2*time.Second, testCase.description, "created_at")
		assertion.WithinDuration(time.Now(), questionnaire.ModifiedAt, 2*time.Second, testCase.description, "modified_at")
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

func TestDeleteQuestionnaire(t *testing.T) {
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
		args
		expect
	}

	testCases := []test{
		{
			args: args{
				title:        "第1回集会らん☆ぷろ募集アンケート",
				description:  "第1回集会らん☆ぷろ参加者募集",
				resTimeLimit: null.NewTime(time.Time{}, false),
				resSharedTo:  "public",
			},
		},
	}

	for _, testCase := range testCases {
		questionnaire := Questionnaires{
			Title:        testCase.args.title,
			Description:  testCase.args.description,
			ResTimeLimit: testCase.args.resTimeLimit,
			ResSharedTo:  testCase.args.resSharedTo,
		}
		err := db.Create(&questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%s): %w", testCase.description, err)
		}

		questionnaireID := questionnaire.ID
		err = DeleteQuestionnaire(questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		questionnaire = Questionnaires{}
		err = db.
			Unscoped().
			Where("id = ?", questionnaireID).
			Find(&questionnaire).Error
		if err != nil {
			t.Errorf("failed to get questionnaire(%s): %w", testCase.description, err)
		}

		assertion.WithinDuration(time.Now(), questionnaire.DeletedAt.ValueOrZero(), 2*time.Second)
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

	err := DeleteQuestionnaire(invalidQuestionnaireID)
	if !errors.Is(err, ErrNoRecordDeleted) {
		if err == nil {
			t.Errorf("Succeeded with invalid questionnaireID")
		} else {
			t.Errorf("failed to update questionnaire(invalid questionnireID): %w", err)
		}
	}
}

func TestGetQuestionnaires(t *testing.T) {
	assertion := assert.New(t)

	sortFuncMap := map[string]func(questionnaires []QuestionnaireInfo) func(i, j int) bool{
		"created_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.After(questionnaires[j].CreatedAt)
			}
		},
		"-created_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.Before(questionnaires[j].CreatedAt)
			}
		},
		"title": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].Title < questionnaires[j].Title
			}
		},
		"-title": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].Title > questionnaires[j].Title
			}
		},
		"modified_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ModifiedAt.After(questionnaires[j].ModifiedAt)
			}
		},
		"-modified_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.After(questionnaires[j].CreatedAt)
			}
		},
		"": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ID < questionnaires[j].ID
			}
		},
	}

	type args struct {
		userID      string
		sort        string
		search      string
		pageNum     int
		nontargeted bool
	}
	type expect struct {
		isErr      bool
		err        error
		isCheckLen bool
		length     int
	}
	type test struct {
		description string
		args
		expect
	}

	testUserID := "mds_boy"

	now := time.Now()
	datas := []QuestionnaireInfo{
		{
			Questionnaires: Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    now,
				ModifiedAt:   now,
			},
			IsTargeted: false,
		},
		{
			Questionnaires: Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    now.Add(time.Second),
				ModifiedAt:   now.Add(2 * time.Second),
			},
			IsTargeted: false,
		},
		{
			Questionnaires: Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    now.Add(2 * time.Second),
				ModifiedAt:   now.Add(3 * time.Second),
			},
			IsTargeted: true,
		},
		{
			Questionnaires: Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    now,
				ModifiedAt:   now,
				DeletedAt:    null.NewTime(now, true),
			},
			IsTargeted: false,
		},
	}
	for i := 0; i < 20; i++ {
		datas = append(datas, QuestionnaireInfo{
			Questionnaires: Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケート",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    now.Add(time.Duration(len(datas)) * time.Second),
				ModifiedAt:   now,
			},
			IsTargeted: false,
		})
	}
	datas = append(datas, QuestionnaireInfo{
		Questionnaires: Questionnaires{
			Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
			Description:  "第1回集会らん☆ぷろ参加者募集",
			ResTimeLimit: null.NewTime(time.Time{}, false),
			ResSharedTo:  "public",
			CreatedAt:    now.Add(2 * time.Second),
			ModifiedAt:   now.Add(3 * time.Second),
		},
		IsTargeted: true,
	})

	for i, data := range datas {
		err := db.Create(&datas[i].Questionnaires).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %w", data, err)
		}

		if data.IsTargeted {
			err := db.Create(Targets{
				QuestionnaireID: datas[i].Questionnaires.ID,
				UserTraqid:      testUserID,
			}).Error
			if err != nil {
				t.Errorf("failed to create target: %w", err)
			}
		}
	}

	deletedQuestionnaireIDs := []int{}
	targettedQuestionnaireIDs := []int{}
	for _, data := range datas {
		if data.Questionnaires.DeletedAt.Valid {
			deletedQuestionnaireIDs = append(deletedQuestionnaireIDs, data.Questionnaires.ID)
		}

		if data.IsTargeted {
			targettedQuestionnaireIDs = append(targettedQuestionnaireIDs, data.Questionnaires.ID)
		}
	}

	testCases := []test{
		{
			description: "userID:valid, sort:no, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:created_at, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "created_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:-created_at, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "-created_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:title, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "title",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:-title, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "-title",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:modified_at, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "modified_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:-modified_at, search:no, page:1",
			args: args{
				userID:      testUserID,
				sort:        "-modified_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:no, search:GetQuestionnaireTest$, page:1",
			args: args{
				userID:      testUserID,
				sort:        "",
				search:      "GetQuestionnaireTest$",
				pageNum:     1,
				nontargeted: false,
			},
			expect: expect{
				isCheckLen: true,
				length:     4,
			},
		},
		{
			description: "userID:valid, sort:no, search:no, page:2",
			args: args{
				userID:      testUserID,
				sort:        "",
				search:      "",
				pageNum:     2,
				nontargeted: false,
			},
		},
		{
			description: "too large page",
			args: args{
				userID:      testUserID,
				sort:        "",
				search:      "",
				pageNum:     100000,
				nontargeted: false,
			},
			expect: expect{
				isErr: true,
				err:   ErrTooLargePageNum,
			},
		},
		{
			description: "userID:valid, sort:no, search:no, page:1, nontargetted",
			args: args{
				userID:      testUserID,
				sort:        "",
				search:      "",
				pageNum:     1,
				nontargeted: true,
			},
		},
		{
			description: "userID:valid, sort:no, search:notFoundQuestionnaire, page:1",
			args: args{
				userID:      testUserID,
				sort:        "",
				search:      "notFoundQuestionnaire",
				pageNum:     1,
				nontargeted: true,
			},
			expect: expect{
				isCheckLen: false,
				length:     0,
			},
		},
	}

	for _, testCase := range testCases {
		questionnaires, pageMax, err := GetQuestionnaires(testCase.args.userID, testCase.args.sort, testCase.args.search, testCase.args.pageNum, testCase.args.nontargeted)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			if !errors.Is(err, testCase.expect.err) {
				t.Errorf("invalid error(%s): expected: %+v, actual: %+v", testCase.description, testCase.expect.err, err)
			}
		}
		if err != nil {
			continue
		}

		var questionnaireNum int
		err = db.
			Model(&Questionnaires{}).
			Where("deleted_at IS NULL").
			Count(&questionnaireNum).Error
		if err != nil {
			t.Errorf("failed to count questionnaire(%s): %w", testCase.description, err)
		}

		actualQuestionnaireIDs := []int{}
		for _, questionnaire := range questionnaires {
			actualQuestionnaireIDs = append(actualQuestionnaireIDs, questionnaire.ID)
		}
		if testCase.args.nontargeted {
			for _, targettedQuestionnaireID := range targettedQuestionnaireIDs {
				assertion.NotContains(actualQuestionnaireIDs, targettedQuestionnaireID, testCase.description, "not contain(targetted)")
			}
		}
		for _, deletedQuestionnaireID := range deletedQuestionnaireIDs {
			assertion.NotContains(actualQuestionnaireIDs, deletedQuestionnaireID, testCase.description, "not contain(deleted)")
		}

		for _, questionnaire := range questionnaires {
			assertion.Regexp(testCase.args.search, questionnaire.Title, testCase.description, "regexp")
		}

		if len(testCase.args.search) == 0 && !testCase.args.nontargeted {
			assertion.Equal((questionnaireNum+19)/20, pageMax, testCase.description, "pageMax")
			assertion.Len(questionnaires, int(math.Min(float64(questionnaireNum-20*(testCase.pageNum-1)), 20.0)), testCase.description, "page")
		}

		if testCase.expect.isCheckLen {
			assertion.Len(questionnaires, testCase.expect.length, testCase.description, "length")
		}

		copyQuestionnaires := make([]QuestionnaireInfo, len(questionnaires))
		copy(copyQuestionnaires, questionnaires)
		sort.SliceStable(copyQuestionnaires, sortFuncMap[testCase.args.sort](questionnaires))
		expectQuestionnaireIDs := []int{}
		for _, questionnaire := range copyQuestionnaires {
			expectQuestionnaireIDs = append(expectQuestionnaireIDs, questionnaire.ID)
		}
		assertion.ElementsMatch(expectQuestionnaireIDs, actualQuestionnaireIDs, testCase.description, "sort")
	}
}
