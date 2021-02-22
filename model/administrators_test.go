package model

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type administratorsTestQuestionnairesTestData struct {
	questionnaire  Questionnaires
	administrators []string
}

const (
	invalidAdministratorTestUserID = "invalidAdministratorsUser"
)

var (
	administratorsTestUserIDs           = []string{"administratorsUser0", "administratorsUser1"}
	administrotorTestQuestionnaireDatas []administratorsTestQuestionnairesTestData
)

func TestAministrators(t *testing.T) {
	t.Parallel()

	setupAdministratorTest(t)

	t.Run("InsertAdministrators", insertAdministratorsTest)
}

func setupAdministratorTest(t *testing.T) {
	administrotorTestQuestionnaireDatas = []administratorsTestQuestionnairesTestData{
		{
			questionnaire: Questionnaires{
				Title:       "第1回集会らん☆ぷろ募集アンケート",
				Description: "第1回集会らん☆ぷろ参加者募集",
			},
			administrators: []string{administratorsTestUserIDs[0]},
		},
		{
			questionnaire: Questionnaires{
				Title:       "第1回集会らん☆ぷろ募集アンケート",
				Description: "第1回集会らん☆ぷろ参加者募集",
			},
			administrators: []string{},
		},
	}

	for i, questionnaireData := range administrotorTestQuestionnaireDatas {
		err := db.Create(&administrotorTestQuestionnaireDatas[i].questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %w", questionnaireData, err)
		}

		for _, administrator := range questionnaireData.administrators {
			err = db.Create(&Administrators{
				QuestionnaireID: administrotorTestQuestionnaireDatas[i].questionnaire.ID,
				UserTraqid:      administrator,
			}).Error
			if err != nil {
				t.Errorf("failed to create administrator(%s): %w", administrator, err)
			}
		}
	}
}

func insertAdministratorsTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionnaire  Questionnaires
		administrators []string
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
			description: "questionnaireID: valid, administrator_num: 2",
			args: args{
				questionnaire: Questionnaires{
					Title:       "第1回集会らん☆ぷろ募集アンケート",
					Description: "第1回集会らん☆ぷろ参加者募集",
				},
				administrators: []string{administratorsTestUserIDs[0], administratorsTestUserIDs[1]},
			},
		},
		{
			description: "questionnaireID: valid, administrator_num: 1",
			args: args{
				questionnaire: Questionnaires{
					Title:       "第1回集会らん☆ぷろ募集アンケート",
					Description: "第1回集会らん☆ぷろ参加者募集",
				},
				administrators: []string{administratorsTestUserIDs[0]},
			},
		},
		{
			description: "questionnaireID: valid, administrator_num: 0",
			args: args{
				questionnaire: Questionnaires{
					Title:       "第1回集会らん☆ぷろ募集アンケート",
					Description: "第1回集会らん☆ぷろ参加者募集",
				},
				administrators: []string{},
			},
		},
		{
			description: "questionnaireID: invalid, administrator_num: 1",
			args: args{
				questionnaire: Questionnaires{
					Title:       "第1回集会らん☆ぷろ募集アンケート",
					Description: "第1回集会らん☆ぷろ参加者募集",
				},
				administrators: []string{administratorsTestUserIDs[0]},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "questionnaireID: invalid, administrator_num: 0",
			args: args{
				questionnaire: Questionnaires{
					Title:       "第1回集会らん☆ぷろ募集アンケート",
					Description: "第1回集会らん☆ぷろ参加者募集",
				},
				administrators: []string{},
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		err := db.Create(&testCase.args.questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %w", testCase.args.questionnaire, err)
		}

		err = administratorImpl.InsertAdministrators(testCase.args.questionnaire.ID, testCase.args.administrators)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		for _, administrator := range testCase.administrators {
			var actualAdministrators Administrators
			err = db.Where("questionnaire_id = ? AND user_traqid = ?", testCase.args.questionnaire.ID, administrator).First(&actualAdministrators).Error

			if gorm.IsRecordNotFoundError(err) {
				t.Errorf("no administrator(%s): %w", administrator, err)
			}
		}
	}
}
