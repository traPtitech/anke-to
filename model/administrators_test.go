package model

import "testing"

type administratorsTestQuestionnairesTestData struct {
	questionnaire  Questionnaires
	administrators []string
}

const (
	administratorsTestUserID       = "administratorsUser"
	invalidAdministratorTestUserID = "invalidAdministratorsUser"
)

var (
	administrotorTestQuestionnaireDatas []administratorsTestQuestionnairesTestData
)

func TestAministrators(t *testing.T) {
	t.Parallel()

	setupAdministratorTest(t)
}

func setupAdministratorTest(t *testing.T) {
	administrotorTestQuestionnaireDatas = []administratorsTestQuestionnairesTestData{
		{
			questionnaire: Questionnaires{
				Title:       "第1回集会らん☆ぷろ募集アンケート",
				Description: "第1回集会らん☆ぷろ参加者募集",
			},
			administrators: []string{administratorsTestUserID},
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
