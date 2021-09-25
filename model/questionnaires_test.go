package model

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

const questionnairesTestUserID = "questionnairesUser"
const questionnairesTestUserID2 = "questionnairesUser2"
const invalidQuestionnairesTestUserID = "invalidQuestionnairesUser"

var questionnairesNow = time.Now()

type QuestionnairesTestData struct {
	questionnaire  *Questionnaires
	targets        []string
	administrators []string
	respondents    []*QuestionnairesTestRespondent
}

type QuestionnairesTestRespondent struct {
	respondent  *Respondents
	isSubmitted bool
}

var (
	datas                   = []*QuestionnairesTestData{}
	deletedQuestionnaireIDs = []int{}
	userTargetMap           = map[string][]int{}
	userAdministratorMap    = map[string][]int{}
	userRespondentMap       = map[string][]int{}
)

func TestQuestionnaires(t *testing.T) {
	t.Parallel()

	setupQuestionnairesTest(t)

	t.Run("InsertQuestionnaire", insertQuestionnaireTest)
	t.Run("UpdateQuestionnaire", updateQuestionnaireTest)
	t.Run("DeleteQuestionnaire", deleteQuestionnaireTest)
	t.Run("GetQuestionnaires", getQuestionnairesTest)
	t.Run("GetAdminQuestionnaires", getAdminQuestionnairesTest)
	t.Run("GetQuestionnaireInfo", getQuestionnaireInfoTest)
	t.Run("GetTargettedQuestionnaires", getTargettedQuestionnairesTest)
	t.Run("GetQuestionnaireLimit", getQuestionnaireLimitTest)
	t.Run("GetQuestionnaireLimitByResponseID", getQuestionnaireLimitByResponseIDTest)
	t.Run("GetResponseReadPrivilegeInfoByResponseID", getResponseReadPrivilegeInfoByResponseIDTest)
	t.Run("GetResponseReadPrivilegeInfoByQuestionnaireID", getResponseReadPrivilegeInfoByQuestionnaireIDTest)
}

func setupQuestionnairesTest(t *testing.T) {
	datas = []*QuestionnairesTestData{
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(questionnairesNow, true),
				ResSharedTo:  "public",
				CreatedAt:    questionnairesNow,
				ModifiedAt:   questionnairesNow,
			},
			targets:        []string{},
			administrators: []string{},
			respondents:    []*QuestionnairesTestRespondent{},
		},
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "respondents",
				CreatedAt:    questionnairesNow,
				ModifiedAt:   questionnairesNow,
			},
			targets:        []string{questionnairesTestUserID},
			administrators: []string{},
			respondents:    []*QuestionnairesTestRespondent{},
		},
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケート",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "administrators",
				CreatedAt:    questionnairesNow.Add(time.Second),
				ModifiedAt:   questionnairesNow.Add(2 * time.Second),
			},
			targets:        []string{},
			administrators: []string{questionnairesTestUserID},
			respondents:    []*QuestionnairesTestRespondent{},
		},
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケート",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    questionnairesNow,
				ModifiedAt:   questionnairesNow,
			},
			targets:        []string{},
			administrators: []string{},
			respondents: []*QuestionnairesTestRespondent{
				{
					respondent: &Respondents{
						UserTraqid: questionnairesTestUserID,
					},
					isSubmitted: true,
				},
			},
		},
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケート",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    questionnairesNow,
				ModifiedAt:   questionnairesNow,
			},
			targets:        []string{},
			administrators: []string{},
			respondents: []*QuestionnairesTestRespondent{
				{
					respondent: &Respondents{
						UserTraqid: questionnairesTestUserID,
					},
				},
			},
		},
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    questionnairesNow.Add(2 * time.Second),
				ModifiedAt:   questionnairesNow.Add(3 * time.Second),
			},
			targets:        []string{questionnairesTestUserID},
			administrators: []string{questionnairesTestUserID},
			respondents:    []*QuestionnairesTestRespondent{},
		},
		{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    questionnairesNow,
				ModifiedAt:   questionnairesNow,
				DeletedAt:    null.NewTime(questionnairesNow, true),
			},
			targets:        []string{},
			administrators: []string{},
			respondents:    []*QuestionnairesTestRespondent{},
		},
	}
	for i := 0; i < 20; i++ {
		datas = append(datas, &QuestionnairesTestData{
			questionnaire: &Questionnaires{
				Title:        "第1回集会らん☆ぷろ募集アンケート",
				Description:  "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit: null.NewTime(time.Time{}, false),
				ResSharedTo:  "public",
				CreatedAt:    questionnairesNow.Add(time.Duration(len(datas)) * time.Second),
				ModifiedAt:   questionnairesNow,
			},
			targets:        []string{},
			administrators: []string{},
			respondents:    []*QuestionnairesTestRespondent{},
		})
	}
	datas = append(datas, &QuestionnairesTestData{
		questionnaire: &Questionnaires{
			Title:        "第1回集会らん☆ぷろ募集アンケートGetQuestionnaireTest",
			Description:  "第1回集会らん☆ぷろ参加者募集",
			ResTimeLimit: null.NewTime(time.Time{}, false),
			ResSharedTo:  "public",
			CreatedAt:    questionnairesNow.Add(2 * time.Second),
			ModifiedAt:   questionnairesNow.Add(3 * time.Second),
		},
		targets:        []string{questionnairesTestUserID},
		administrators: []string{questionnairesTestUserID},
		respondents:    []*QuestionnairesTestRespondent{},
	}, &QuestionnairesTestData{
		questionnaire: &Questionnaires{
			Title:        "第1回集会らん☆ぷろ募集アンケート",
			Description:  "第1回集会らん☆ぷろ参加者募集",
			ResTimeLimit: null.NewTime(time.Time{}, false),
			ResSharedTo:  "public",
			CreatedAt:    questionnairesNow,
			ModifiedAt:   questionnairesNow,
		},
		targets:        []string{},
		administrators: []string{questionnairesTestUserID},
		respondents: []*QuestionnairesTestRespondent{
			{
				respondent: &Respondents{
					UserTraqid: questionnairesTestUserID,
				},
				isSubmitted: true,
			},
		},
	}, &QuestionnairesTestData{
		questionnaire: &Questionnaires{
			Title:        "第1回集会らん☆ぷろ募集アンケート",
			Description:  "第1回集会らん☆ぷろ参加者募集",
			ResTimeLimit: null.NewTime(time.Time{}, false),
			ResSharedTo:  "public",
			CreatedAt:    questionnairesNow,
			ModifiedAt:   questionnairesNow,
		},
		targets:        []string{},
		administrators: []string{},
		respondents: []*QuestionnairesTestRespondent{
			{
				respondent: &Respondents{
					UserTraqid: questionnairesTestUserID,
				},
				isSubmitted: true,
			},
			{
				respondent: &Respondents{
					UserTraqid: questionnairesTestUserID2,
				},
			},
		},
	}, &QuestionnairesTestData{
		questionnaire: &Questionnaires{
			Title:        "第1回集会らん☆ぷろ募集アンケート",
			Description:  "第1回集会らん☆ぷろ参加者募集",
			ResTimeLimit: null.NewTime(questionnairesNow, true),
			ResSharedTo:  "public",
			CreatedAt:    questionnairesNow,
			ModifiedAt:   questionnairesNow,
		},
		targets:        []string{},
		administrators: []string{questionnairesTestUserID},
		respondents: []*QuestionnairesTestRespondent{
			{
				respondent: &Respondents{
					UserTraqid: questionnairesTestUserID,
				},
				isSubmitted: true,
			},
		},
	}, &QuestionnairesTestData{
		questionnaire: &Questionnaires{
			Title:        "第1回集会らん☆ぷろ募集アンケート",
			Description:  "第1回集会らん☆ぷろ参加者募集",
			ResTimeLimit: null.NewTime(time.Time{}, false),
			ResSharedTo:  "public",
			CreatedAt:    questionnairesNow,
			ModifiedAt:   questionnairesNow,
		},
		targets:        []string{},
		administrators: []string{questionnairesTestUserID},
		respondents: []*QuestionnairesTestRespondent{
			{
				respondent: &Respondents{
					UserTraqid: questionnairesTestUserID,
				},
				isSubmitted: true,
			},
		},
	})

	for i, data := range datas {
		if data.questionnaire.DeletedAt.Valid {
			deletedQuestionnaireIDs = append(deletedQuestionnaireIDs, data.questionnaire.ID)
		}

		err := db.Create(data.questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %w", data, err)
		}

		for _, target := range data.targets {
			questionnaires, ok := userTargetMap[target]
			if !ok {
				questionnaires = []int{}
			}
			userTargetMap[target] = append(questionnaires, datas[i].questionnaire.ID)

			err := db.Create(&Targets{
				QuestionnaireID: data.questionnaire.ID,
				UserTraqid:      target,
			}).Error
			if err != nil {
				t.Errorf("failed to create target: %w", err)
			}
		}

		for _, administrator := range data.administrators {
			questionnaires, ok := userAdministratorMap[administrator]
			if !ok {
				questionnaires = []int{}
			}
			userAdministratorMap[administrator] = append(questionnaires, datas[i].questionnaire.ID)

			err := db.Create(&Administrators{
				QuestionnaireID: data.questionnaire.ID,
				UserTraqid:      administrator,
			}).Error
			if err != nil {
				t.Errorf("failed to create target: %w", err)
			}
		}

		for _, respondentData := range data.respondents {
			if respondentData.isSubmitted {
				questionnaires, ok := userRespondentMap[respondentData.respondent.UserTraqid]
				if !ok {
					questionnaires = []int{}
				}
				userRespondentMap[respondentData.respondent.UserTraqid] = append(questionnaires, datas[i].questionnaire.ID)
			}

			respondentData.respondent.QuestionnaireID = data.questionnaire.ID
			if respondentData.isSubmitted {
				respondentData.respondent.SubmittedAt = null.NewTime(time.Now(), true)
			}

			err := db.Transaction(func(tx *gorm.DB) error {
				err := db.Create(respondentData.respondent).Error
				if err != nil {
					return fmt.Errorf("failed to create respondent: %w", err)
				}

				err = db.Order("response_id DESC").Last(respondentData.respondent).Error
				if err != nil {
					return fmt.Errorf("failed to get respondent: %w", err)
				}

				return nil
			})
			if err != nil {
				t.Errorf("failed in the transaction: %w", err)
			}
		}
	}
}

func insertQuestionnaireTest(t *testing.T) {
	t.Helper()
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
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(testCase.args.title, testCase.args.description, testCase.args.resTimeLimit, testCase.args.resSharedTo)

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

func updateQuestionnaireTest(t *testing.T) {
	t.Helper()
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
		err = questionnaireImpl.UpdateQuestionnaire(after.title, after.description, after.resTimeLimit, after.resSharedTo, questionnaireID)

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
		err := questionnaireImpl.UpdateQuestionnaire(arg.title, arg.description, arg.resTimeLimit, arg.resSharedTo, invalidQuestionnaireID)
		if !errors.Is(err, ErrNoRecordUpdated) {
			if err == nil {
				t.Errorf("Succeeded with invalid questionnaireID")
			} else {
				t.Errorf("failed to update questionnaire(invalid questionnireID): %w", err)
			}
		}
	}
}

func deleteQuestionnaireTest(t *testing.T) {
	t.Helper()
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
		err = questionnaireImpl.DeleteQuestionnaire(questionnaireID)

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

	err := questionnaireImpl.DeleteQuestionnaire(invalidQuestionnaireID)
	if !errors.Is(err, ErrNoRecordDeleted) {
		if err == nil {
			t.Errorf("Succeeded with invalid questionnaireID")
		} else {
			t.Errorf("failed to update questionnaire(invalid questionnireID): %w", err)
		}
	}
}

func getQuestionnairesTest(t *testing.T) {
	t.Helper()

	assertion := assert.New(t)

	sortFuncMap := map[string]func(questionnaires []QuestionnaireInfo) func(i, j int) bool{
		"created_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.Before(questionnaires[j].CreatedAt) || (questionnaires[i].CreatedAt.Equal(questionnaires[j].CreatedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"-created_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.After(questionnaires[j].CreatedAt) || (questionnaires[i].CreatedAt.Equal(questionnaires[j].CreatedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"title": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].Title < questionnaires[j].Title || (questionnaires[i].Title == questionnaires[j].Title && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"-title": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].Title > questionnaires[j].Title || (questionnaires[i].Title == questionnaires[j].Title && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"modified_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ModifiedAt.Before(questionnaires[j].ModifiedAt) || (questionnaires[i].ModifiedAt.Equal(questionnaires[j].ModifiedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"-modified_at": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ModifiedAt.After(questionnaires[j].ModifiedAt) || (questionnaires[i].ModifiedAt.Equal(questionnaires[j].ModifiedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"": func(questionnaires []QuestionnaireInfo) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ID > questionnaires[j].ID
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

	testCases := []test{
		{
			description: "userID:valid, sort:no, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:created_at, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "created_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:-created_at, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "-created_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:title, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "title",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:-title, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "-title",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:modified_at, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "modified_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:-modified_at, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "-modified_at",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
		},
		{
			description: "userID:valid, sort:no, search:GetQuestionnaireTest$, page:1",
			args: args{
				userID:      questionnairesTestUserID,
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
				userID:      questionnairesTestUserID,
				sort:        "",
				search:      "",
				pageNum:     2,
				nontargeted: false,
			},
		},
		{
			description: "too large page",
			args: args{
				userID:      questionnairesTestUserID,
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
				userID:      questionnairesTestUserID,
				sort:        "",
				search:      "",
				pageNum:     1,
				nontargeted: true,
			},
		},
		{
			description: "userID:valid, sort:no, search:notFoundQuestionnaire, page:1",
			args: args{
				userID:      questionnairesTestUserID,
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
		{
			description: "userID:valid, sort:invalid, search:no, page:1",
			args: args{
				userID:      questionnairesTestUserID,
				sort:        "hogehoge",
				search:      "",
				pageNum:     1,
				nontargeted: false,
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidSortParam,
			},
		},
	}

	for _, testCase := range testCases {
		questionnaires, pageMax, err := questionnaireImpl.GetQuestionnaires(testCase.args.userID, testCase.args.sort, testCase.args.search, testCase.args.pageNum, testCase.args.nontargeted)

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
			for _, targettedQuestionnaireID := range userTargetMap[questionnairesTestUserID] {
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
		sort.Slice(copyQuestionnaires, sortFuncMap[testCase.args.sort](copyQuestionnaires))
		expectQuestionnaireIDs := make([]int, 0, len(copyQuestionnaires))
		for _, questionnaire := range copyQuestionnaires {
			expectQuestionnaireIDs = append(expectQuestionnaireIDs, questionnaire.ID)
		}
		assertion.Equal(expectQuestionnaireIDs, actualQuestionnaireIDs, testCase.description, "sort")
	}
}

func getAdminQuestionnairesTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		userID string
	}
	type expect struct {
		isCheckLen bool
		length     int
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
			description: "user:valid",
			args: args{
				userID: questionnairesTestUserID,
			},
		},
		{
			description: "empty response",
			args: args{
				userID: invalidQuestionnairesTestUserID,
			},
			expect: expect{
				isCheckLen: true,
				length:     0,
			},
		},
	}

	for _, testCase := range testCases {
		questionnaires, err := questionnaireImpl.GetAdminQuestionnaires(testCase.userID)

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

		actualQuestionnaireIDs := make([]int, 0, len(questionnaires))
		actualIDQuestionnaireMap := make(map[int]Questionnaires, len(questionnaires))
		for _, questionnaire := range questionnaires {
			actualQuestionnaireIDs = append(actualQuestionnaireIDs, questionnaire.ID)
			actualIDQuestionnaireMap[questionnaire.ID] = questionnaire
		}

		if testCase.expect.isCheckLen {
			assertion.Len(questionnaires, testCase.expect.length, testCase.description, "length")
		}

		assertion.Subset(userAdministratorMap[testCase.args.userID], actualQuestionnaireIDs, testCase.description, "administrate")

		expectQuestionnaires := []Questionnaires{}
		err = db.
			Where("id IN (?)", actualQuestionnaireIDs).
			Find(&expectQuestionnaires).Error
		if err != nil {
			t.Errorf("failed to get questionnaires(%s): %w", testCase.description, err)
		}

		for _, expectQuestionnaire := range expectQuestionnaires {
			actualQuestionnaire := actualIDQuestionnaireMap[expectQuestionnaire.ID]

			assertion.Equal(expectQuestionnaire, actualQuestionnaire, testCase.description, "questionnaire")
		}
	}
}

func getQuestionnaireInfoTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionnaireID int
	}
	type expect struct {
		questionnaire  Questionnaires
		targets        []string
		administrators []string
		respondents    []string
		isErr          bool
		err            error
	}
	type test struct {
		description string
		args
		expect
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

	testCases := []test{
		{
			description: "respondents: no, targets: no, administrator: no",
			args: args{
				questionnaireID: datas[0].questionnaire.ID,
			},
			expect: expect{
				questionnaire:  *datas[0].questionnaire,
				targets:        []string{},
				administrators: []string{},
				respondents:    []string{},
			},
		},
		{
			description: "respondents: no, targets: valid, administrator: no",
			args: args{
				questionnaireID: datas[1].questionnaire.ID,
			},
			expect: expect{
				questionnaire:  *datas[1].questionnaire,
				targets:        []string{questionnairesTestUserID},
				administrators: []string{},
				respondents:    []string{},
			},
		},
		{
			description: "respondents: no, targets: no, administrator: valid",
			args: args{
				questionnaireID: datas[2].questionnaire.ID,
			},
			expect: expect{
				questionnaire:  *datas[2].questionnaire,
				targets:        []string{},
				administrators: []string{questionnairesTestUserID},
				respondents:    []string{},
			},
		},
		{
			description: "respondents: submitted, targets: no, administrator: no",
			args: args{
				questionnaireID: datas[3].questionnaire.ID,
			},
			expect: expect{
				questionnaire:  *datas[3].questionnaire,
				targets:        []string{},
				administrators: []string{},
				respondents:    []string{questionnairesTestUserID},
			},
		},
		{
			description: "respondents: saved, targets: no, administrator: no",
			args: args{
				questionnaireID: datas[4].questionnaire.ID,
			},
			expect: expect{
				questionnaire:  *datas[4].questionnaire,
				targets:        []string{},
				administrators: []string{},
				respondents:    []string{},
			},
		},
		{
			description: "questionnaireID: invalid",
			args: args{
				questionnaireID: invalidQuestionnaireID,
			},
			expect: expect{
				isErr: true,
				err:   gorm.ErrRecordNotFound,
			},
		},
	}

	for _, testCase := range testCases {
		actualQuestionnaire, actualTargets, actualAdministrators, actualRespondents, err := questionnaireImpl.GetQuestionnaireInfo(testCase.questionnaireID)

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

		assertion.Equal(testCase.expect.questionnaire.ID, actualQuestionnaire.ID, testCase.description, "questionnaire(ID)")
		assertion.Equal(testCase.expect.questionnaire.Title, actualQuestionnaire.Title, testCase.description, "questionnaire(Title)")
		assertion.Equal(testCase.expect.questionnaire.Description, actualQuestionnaire.Description, testCase.description, "questionnaire(Description)")
		assertion.Equal(testCase.expect.questionnaire.ResSharedTo, actualQuestionnaire.ResSharedTo, testCase.description, "questionnaire(ResSharedTo)")
		assertion.WithinDuration(testCase.expect.questionnaire.ResTimeLimit.ValueOrZero(), actualQuestionnaire.ResTimeLimit.ValueOrZero(), 2*time.Second, testCase.description, "questionnaire(ResTimeLimit)")
		assertion.WithinDuration(testCase.expect.questionnaire.CreatedAt, actualQuestionnaire.CreatedAt, 2*time.Second, testCase.description, "questionnaire(CreatedAt)")
		assertion.WithinDuration(testCase.expect.questionnaire.ModifiedAt, actualQuestionnaire.ModifiedAt, 2*time.Second, testCase.description, "questionnaire(ModifiedAt)")
		assertion.WithinDuration(testCase.expect.questionnaire.DeletedAt.ValueOrZero(), actualQuestionnaire.DeletedAt.ValueOrZero(), 2*time.Second, testCase.description, "questionnaire(DeletedAt)")

		sort.Slice(testCase.targets, func(i, j int) bool { return testCase.targets[i] < testCase.targets[j] })
		sort.Slice(actualTargets, func(i, j int) bool { return actualTargets[i] < actualTargets[j] })
		assertion.Equal(testCase.targets, actualTargets, testCase.description, "targets")

		sort.Slice(testCase.administrators, func(i, j int) bool { return testCase.administrators[i] < testCase.administrators[j] })
		sort.Slice(actualAdministrators, func(i, j int) bool { return actualAdministrators[i] < actualAdministrators[j] })
		assertion.Equal(testCase.administrators, actualAdministrators, testCase.description, "administrators")

		sort.Slice(testCase.respondents, func(i, j int) bool { return testCase.respondents[i] < testCase.respondents[j] })
		sort.Slice(actualRespondents, func(i, j int) bool { return actualRespondents[i] < actualRespondents[j] })
		assertion.Equal(testCase.respondents, actualRespondents, testCase.description, "respondents")
	}
}

func getTargettedQuestionnairesTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	sortFuncMap := map[string]func(questionnaires []TargettedQuestionnaire) func(i, j int) bool{
		"created_at": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.Before(questionnaires[j].CreatedAt) || (questionnaires[i].CreatedAt.Equal(questionnaires[j].CreatedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"-created_at": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].CreatedAt.After(questionnaires[j].CreatedAt) || (questionnaires[i].CreatedAt.Equal(questionnaires[j].CreatedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"title": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].Title < questionnaires[j].Title || (questionnaires[i].Title == questionnaires[j].Title && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"-title": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].Title > questionnaires[j].Title || (questionnaires[i].Title == questionnaires[j].Title && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"modified_at": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ModifiedAt.Before(questionnaires[j].ModifiedAt) || (questionnaires[i].ModifiedAt.Equal(questionnaires[j].ModifiedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"-modified_at": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ModifiedAt.After(questionnaires[j].ModifiedAt) || (questionnaires[i].ModifiedAt.Equal(questionnaires[j].ModifiedAt) && questionnaires[i].ID > questionnaires[j].ID)
			}
		},
		"": func(questionnaires []TargettedQuestionnaire) func(i, j int) bool {
			return func(i, j int) bool {
				return questionnaires[i].ID > questionnaires[j].ID
			}
		},
	}

	type args struct {
		userID   string
		answered string
		sort     string
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
			description: "userID: valid, answered: no, sort: no",
			args: args{
				userID: questionnairesTestUserID,
			},
		},
		{
			description: "userID: valid, answered: answered, sort: no",
			args: args{
				userID:   questionnairesTestUserID,
				answered: "answered",
			},
		},
		{
			description: "userID: valid, answered: unanswered, sort: no",
			args: args{
				userID:   questionnairesTestUserID,
				answered: "unanswered",
			},
		},
		{
			description: "userID: valid, answered: no, sort: created_at",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "created_at",
			},
		},
		{
			description: "userID: valid, answered: no, sort: -created_at",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "-created_at",
			},
		},
		{
			description: "userID: valid, answered: no, sort: title",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "title",
			},
		},
		{
			description: "userID: valid, answered: no, sort: -title",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "-title",
			},
		},
		{
			description: "userID: valid, answered: no, sort: modified_at",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "modified_at",
			},
		},
		{
			description: "userID: valid, answered: no, sort: -modified_at",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "-modified_at",
			},
		},
		{
			description: "userID: valid, answered: invalid, sort: no",
			args: args{
				userID:   questionnairesTestUserID,
				answered: "invalidAnswered",
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidAnsweredParam,
			},
		},
		{
			description: "userID: valid, answered: no, sort: invalid",
			args: args{
				userID: questionnairesTestUserID,
				sort:   "invalidSort",
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidSortParam,
			},
		},
	}

	for _, testCase := range testCases {
		questionnaires, err := questionnaireImpl.GetTargettedQuestionnaires(testCase.args.userID, testCase.args.answered, testCase.args.sort)

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

		actualQuestionnaireIDs := []int{}
		for _, questionnaire := range questionnaires {
			actualQuestionnaireIDs = append(actualQuestionnaireIDs, questionnaire.ID)
		}
		assertion.Subset(userTargetMap[questionnairesTestUserID], actualQuestionnaireIDs, testCase.description, "contain(targetted)")
		for _, deletedQuestionnaireID := range deletedQuestionnaireIDs {
			assertion.NotContains(actualQuestionnaireIDs, deletedQuestionnaireID, testCase.description, "not contain(deleted)")
		}
		for _, deletedQuestionnaireID := range deletedQuestionnaireIDs {
			assertion.NotContains(actualQuestionnaireIDs, deletedQuestionnaireID, testCase.description, "not contain(deleted)")
		}

		switch testCase.args.answered {
		case "answered":
			for _, questionnaire := range questionnaires {
				assertion.True(questionnaire.RespondedAt.Valid, testCase.description, "answered")
			}
			assertion.Subset(userRespondentMap[questionnairesTestUserID], actualQuestionnaireIDs, testCase.description, "contain(responded)")
		case "unanswered":
			for _, questionnaire := range questionnaires {
				assertion.True(!questionnaire.RespondedAt.Valid, testCase.description, "unanswered")
			}
			for _, respondedQuestionnaireID := range userRespondentMap[questionnairesTestUserID] {
				assertion.NotContains(actualQuestionnaireIDs, respondedQuestionnaireID, testCase.description, "not contain(deleted)")
			}
		}

		copyQuestionnaires := make([]TargettedQuestionnaire, len(questionnaires))
		copy(copyQuestionnaires, questionnaires)
		sort.Slice(copyQuestionnaires, sortFuncMap[testCase.args.sort](copyQuestionnaires))
		expectQuestionnaireIDs := []int{}
		for _, questionnaire := range copyQuestionnaires {
			expectQuestionnaireIDs = append(expectQuestionnaireIDs, questionnaire.ID)
		}
		assertion.Equal(expectQuestionnaireIDs, actualQuestionnaireIDs, testCase.description, "sort")
	}
}

func getQuestionnaireLimitTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

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

	type args struct {
		questionnaireID int
	}
	type expect struct {
		limit null.Time
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
			description: "limit: not null",
			args: args{
				questionnaireID: datas[0].questionnaire.ID,
			},
			expect: expect{
				limit: datas[0].questionnaire.ResTimeLimit,
			},
		},
		{
			description: "limit: null",
			args: args{
				questionnaireID: datas[1].questionnaire.ID,
			},
			expect: expect{
				limit: datas[1].questionnaire.ResTimeLimit,
			},
		},
		{
			description: "questionnaireID: invalid",
			args: args{
				questionnaireID: invalidQuestionnaireID,
			},
			expect: expect{
				isErr: true,
				err:   gorm.ErrRecordNotFound,
			},
		},
	}

	for _, testCase := range testCases {
		actualLimit, err := questionnaireImpl.GetQuestionnaireLimit(testCase.args.questionnaireID)

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

		assertion.WithinDuration(testCase.limit.ValueOrZero(), actualLimit.ValueOrZero(), 2*time.Second, testCase.description, "limit")
	}
}

func getQuestionnaireLimitByResponseIDTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	invalidResponseID := 1000
	for {
		err := db.Where("response_id = ?", invalidResponseID).First(&Respondents{}).Error
		if gorm.IsRecordNotFoundError(err) {
			break
		}
		if err != nil {
			t.Errorf("failed to get response(make invalid responseID): %w", err)
			break
		}

		invalidResponseID *= 10
	}

	type args struct {
		responseID int
	}
	type expect struct {
		limit null.Time
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
			description: "回答が存在しないのでエラー",
			args: args{
				responseID: invalidResponseID,
			},
			expect: expect{
				isErr: true,
				err:   gorm.ErrRecordNotFound,
			},
		},
		{
			description: "limit: not null",
			args: args{
				responseID: datas[30].respondents[0].respondent.ResponseID,
			},
			expect: expect{
				limit: datas[30].questionnaire.ResTimeLimit,
			},
		},
		{
			description: "limit: null",
			args: args{
				responseID: datas[31].respondents[0].respondent.ResponseID,
			},
			expect: expect{
				limit: datas[31].questionnaire.ResTimeLimit,
			},
		},
	}

	for _, testCase := range testCases {
		actualLimit, err := questionnaireImpl.GetQuestionnaireLimitByResponseID(testCase.args.responseID)

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

		if testCase.expect.limit.Valid {
			assertion.WithinDuration(testCase.limit.ValueOrZero(), actualLimit.ValueOrZero(), 2*time.Second, testCase.description, "limit")
		} else {
			assertion.False(actualLimit.Valid, testCase.description, "limit null")
		}
	}
}

func getResponseReadPrivilegeInfoByResponseIDTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

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

	invalidResponseID := 1000
	for {
		err := db.Where("response_id = ?", invalidResponseID).First(&Respondents{}).Error
		if gorm.IsRecordNotFoundError(err) {
			break
		}
		if err != nil {
			t.Errorf("failed to get response(make invalid responseID): %w", err)
			break
		}

		invalidResponseID *= 10
	}

	type args struct {
		userID     string
		responseID int
	}
	type expect struct {
		responseReadPrivilegeInfo *ResponseReadPrivilegeInfo
		isErr                     bool
		err                       error
	}
	type test struct {
		description string
		args
		expect
	}
	testCases := []test{
		{
			description: "回答が存在しないのでエラー",
			args: args{
				responseID: invalidResponseID,
				userID:     questionnairesTestUserID,
			},
			expect: expect{
				isErr: true,
				err:   ErrRecordNotFound,
			},
		},
		{
			description: "下書きの回答なのでエラー",
			args: args{
				responseID: datas[4].respondents[0].respondent.ResponseID,
				userID:     questionnairesTestUserID,
			},
			expect: expect{
				isErr: true,
				err:   ErrRecordNotFound,
			},
		},
		{
			description: "回答はあるが、管理者でも回答者でもない",
			args: args{
				responseID: datas[3].respondents[0].respondent.ResponseID,
				userID:     invalidQuestionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[3].questionnaire.ResSharedTo,
					IsAdministrator: false,
					IsRespondent:    false,
				},
			},
		},
		{
			description: "回答はあり、管理者ではないが回答者ではある",
			args: args{
				responseID: datas[3].respondents[0].respondent.ResponseID,
				userID:     questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[3].questionnaire.ResSharedTo,
					IsAdministrator: false,
					IsRespondent:    true,
				},
			},
		},
		{
			description: "回答はあり、管理者ではあるが回答者ではない",
			args: args{
				responseID: datas[28].respondents[0].respondent.ResponseID,
				userID:     questionnairesTestUserID2,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[28].questionnaire.ResSharedTo,
					IsAdministrator: true,
					IsRespondent:    false,
				},
			},
		},
		{
			description: "回答はあり、管理者かつ回答者ではある",
			args: args{
				responseID: datas[28].respondents[0].respondent.ResponseID,
				userID:     questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[28].questionnaire.ResSharedTo,
					IsAdministrator: true,
					IsRespondent:    true,
				},
			},
		},
		{
			description: "回答はあり、管理者でなく下書きの回答しかないため回答者ではない",
			args: args{
				responseID: datas[29].respondents[0].respondent.ResponseID,
				userID:     questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[29].questionnaire.ResSharedTo,
					IsAdministrator: false,
					IsRespondent:    false,
				},
			},
		},
	}

	for _, testCase := range testCases {
		responseReadPrivilegeInfo, err := questionnaireImpl.GetResponseReadPrivilegeInfoByResponseID(testCase.args.userID, testCase.args.responseID)

		if testCase.expect.isErr {
			if testCase.expect.err == nil {
				assertion.Errorf(err, testCase.description, "no error")
			} else {
				if !errors.Is(err, testCase.expect.err) {
					t.Errorf("invalid error(%s): expected: %+v, actual: %+v", testCase.description, testCase.expect.err, err)
				}
			}
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.responseReadPrivilegeInfo, responseReadPrivilegeInfo, testCase.description, "responseReadPrivilegeInfo")
	}
}

func getResponseReadPrivilegeInfoByQuestionnaireIDTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

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

	type args struct {
		userID          string
		questionnaireID int
	}
	type expect struct {
		responseReadPrivilegeInfo *ResponseReadPrivilegeInfo
		isErr                     bool
		err                       error
	}
	type test struct {
		description string
		args
		expect
	}
	testCases := []test{
		{
			description: "存在しないquestionnaireID",
			args: args{
				questionnaireID: invalidQuestionnaireID,
				userID:          questionnairesTestUserID,
			},
			expect: expect{
				isErr: true,
				err:   ErrRecordNotFound,
			},
		},
		{
			description: "管理者でも回答者でもない",
			args: args{
				questionnaireID: datas[0].questionnaire.ID,
				userID:          questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[0].questionnaire.ResSharedTo,
					IsAdministrator: false,
					IsRespondent:    false,
				},
			},
		},
		{
			description: "管理者ではあるが回答者ではない",
			args: args{
				questionnaireID: datas[2].questionnaire.ID,
				userID:          questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[2].questionnaire.ResSharedTo,
					IsAdministrator: true,
					IsRespondent:    false,
				},
			},
		},
		{
			description: "管理者ではないが回答者ではある",
			args: args{
				questionnaireID: datas[3].questionnaire.ID,
				userID:          questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[3].questionnaire.ResSharedTo,
					IsAdministrator: false,
					IsRespondent:    true,
				},
			},
		},
		{
			description: "管理者でなく下書きの回答しかないため回答者ではない",
			args: args{
				questionnaireID: datas[4].questionnaire.ID,
				userID:          questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[4].questionnaire.ResSharedTo,
					IsAdministrator: false,
					IsRespondent:    false,
				},
			},
		},
		{
			description: "管理者かつ回答者ではある",
			args: args{
				questionnaireID: datas[28].questionnaire.ID,
				userID:          questionnairesTestUserID,
			},
			expect: expect{
				responseReadPrivilegeInfo: &ResponseReadPrivilegeInfo{
					ResSharedTo:     datas[28].questionnaire.ResSharedTo,
					IsAdministrator: true,
					IsRespondent:    true,
				},
			},
		},
	}

	for _, testCase := range testCases {
		responseReadPrivilegeInfo, err := questionnaireImpl.GetResponseReadPrivilegeInfoByQuestionnaireID(testCase.args.userID, testCase.args.questionnaireID)

		if testCase.expect.isErr {
			if testCase.expect.err == nil {
				assertion.Errorf(err, testCase.description, "no error")
			} else {
				if !errors.Is(err, testCase.expect.err) {
					t.Errorf("invalid error(%s): expected: %+v, actual: %+v", testCase.description, testCase.expect.err, err)
				}
			}
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.responseReadPrivilegeInfo, responseReadPrivilegeInfo, testCase.description, "responseReadPrivilegeInfo")
	}
}
