package model

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
	administratorTestQuestionnaireDatas []administratorsTestQuestionnairesTestData
)

func TestAministrators(t *testing.T) {
	t.Parallel()

	setupAdministratorTest(t)

	t.Run("InsertAdministrators", insertAdministratorsTest)
	t.Run("DeleteAdministrators", deleteAdministratorsTest)
	t.Run("GetAdministrators", getAdministratorsTest)
	t.Run("CheckQuestionnaireAdmin", checkQuestionnaireAdminTest)
}

func setupAdministratorTest(t *testing.T) {
	administratorTestQuestionnaireDatas = []administratorsTestQuestionnairesTestData{
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

	for i, questionnaireData := range administratorTestQuestionnaireDatas {
		err := db.
			Session(&gorm.Session{NewDB: true}).
			Create(&administratorTestQuestionnaireDatas[i].questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %v", questionnaireData, err)
		}

		for _, administrator := range questionnaireData.administrators {
			err = db.
				Session(&gorm.Session{NewDB: true}).
				Create(&Administrators{
					QuestionnaireID: administratorTestQuestionnaireDatas[i].questionnaire.ID,
					UserTraqid:      administrator,
				}).Error
			if err != nil {
				t.Errorf("failed to create administrator(%s): %v", administrator, err)
			}
		}
	}
}

func insertAdministratorsTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

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
		err := db.
			Session(&gorm.Session{NewDB: true}).
			Where("id = ?", invalidQuestionnaireID).
			First(&Questionnaires{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		if err != nil {
			t.Errorf("failed to get questionnaire(make invalid questionnaireID): %v", err)
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
		err := db.
			Session(&gorm.Session{NewDB: true}).
			Create(&testCase.args.questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %v", testCase.args.questionnaire, err)
		}

		err = administratorImpl.InsertAdministrators(ctx, testCase.args.questionnaire.ID, testCase.args.administrators)

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
			err = db.
				Session(&gorm.Session{NewDB: true}).
				Where("questionnaire_id = ? AND user_traqid = ?", testCase.args.questionnaire.ID, administrator).
				First(&actualAdministrators).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("no administrator(%s): %v", administrator, err)
			}
		}
	}
}

func deleteAdministratorsTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

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
		err := db.
			Session(&gorm.Session{NewDB: true}).
			Create(&testCase.args.questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %v", testCase.args.questionnaire, err)
		}

		err = administratorImpl.DeleteAdministrators(ctx, testCase.args.questionnaire.ID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		var administrators []Administrators
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("questionnaire_id = ?", testCase.args.questionnaire.ID).
			Find(&administrators).Error
		if err != nil {
			t.Errorf("failed to get administrators(%s): %v", testCase.description, err)
		}

		assertion.Len(administrators, 0, testCase.description, "administrator length")
	}
}

func getAdministratorsTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		questionnaireIDs []int
	}
	type expect struct {
		administrators []Administrators
		isErr          bool
		err            error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "questionnaire_id_num: 2",
			args: args{
				questionnaireIDs: []int{administratorTestQuestionnaireDatas[0].questionnaire.ID, administratorTestQuestionnaireDatas[1].questionnaire.ID},
			},
			expect: expect{
				administrators: []Administrators{
					{
						QuestionnaireID: administratorTestQuestionnaireDatas[0].questionnaire.ID,
						UserTraqid:      administratorsTestUserIDs[0],
					},
				},
			},
		},
		{
			description: "questionnaire_id_num: 1",
			args: args{
				questionnaireIDs: []int{administratorTestQuestionnaireDatas[0].questionnaire.ID},
			},
			expect: expect{
				administrators: []Administrators{
					{
						QuestionnaireID: administratorTestQuestionnaireDatas[0].questionnaire.ID,
						UserTraqid:      administratorsTestUserIDs[0],
					},
				},
			},
		},
		{
			description: "questionnaire_id_num: 1, no administrator",
			args: args{
				questionnaireIDs: []int{administratorTestQuestionnaireDatas[1].questionnaire.ID},
			},
			expect: expect{
				administrators: []Administrators{},
			},
		},
		{
			description: "questionnaire_id_num: 0",
			args: args{
				questionnaireIDs: []int{},
			},
			expect: expect{
				administrators: []Administrators{},
			},
		},
	}

	for _, testCase := range testCases {
		actualAdministrators, err := administratorImpl.GetAdministrators(ctx, testCase.args.questionnaireIDs)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		assertion.ElementsMatch(actualAdministrators, testCase.expect.administrators, testCase.description, "element")
	}
}

func checkQuestionnaireAdminTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		userID          string
		questionnaireID int
	}
	type expect struct {
		isAdmin bool
		isErr   bool
		err     error
	}
	type test struct {
		description string
		args
		expect
	}

	invalidQuestionnaireID := 1000
	for {
		err := db.
			Session(&gorm.Session{NewDB: true}).
			Where("id = ?", invalidQuestionnaireID).
			First(&Questionnaires{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		if err != nil {
			t.Errorf("failed to get questionnaire(make invalid questionnaireID): %v", err)
			break
		}

		invalidQuestionnaireID *= 10
	}

	testCases := []test{
		{
			description: "questionnaireID: valid, is_admin: true",
			args: args{
				userID:          administratorsTestUserIDs[0],
				questionnaireID: administratorTestQuestionnaireDatas[0].questionnaire.ID,
			},
			expect: expect{
				isAdmin: true,
			},
		},
		{
			description: "questionnaireID: valid, is_admin: false",
			args: args{
				userID:          invalidAdministratorTestUserID,
				questionnaireID: administratorTestQuestionnaireDatas[0].questionnaire.ID,
			},
			expect: expect{
				isAdmin: false,
			},
		},
		{
			description: "questionnaireID: invalid",
			args: args{
				userID:          administratorsTestUserIDs[0],
				questionnaireID: invalidQuestionnaireID,
			},
			expect: expect{
				isAdmin: false,
			},
		},
	}

	for _, testCase := range testCases {
		actualIsAdmin, err := administratorImpl.CheckQuestionnaireAdmin(ctx, testCase.args.userID, testCase.args.questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.isAdmin, actualIsAdmin, testCase.description, "isAdmin")
	}
}
