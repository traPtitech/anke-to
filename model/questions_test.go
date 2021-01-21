package model

import (
	"testing"
)

const questionsTestUserID = "questionsUser"
const invalidQuestionsTestUserID = "invalidQuestionsUser"

type QuestionsTestQuestionnaireData struct {
	Questionnaires
	administrators []string
}

type QuestionsTestData struct {
	Question
}

var (
	questionnaireDatas = []QuestionsTestQuestionnaireData{}
	datas              = []QuestionsTestData{}
)

func TestQuestions(t *testing.T) {
	t.Parallel()
}

func setupQuestionsTest(t *testing.T) {
	questionnaireDatas = []QuestionsTestQuestionnaireData{
		{
			Questionnaires: Questionnaires{
				Title:       "第1回集会らん☆ぷろ募集アンケート",
				Description: "第1回集会らん☆ぷろ参加者募集",
			},
			administrators: []string{questionsTestUserID},
		},
		{
			Questionnaires: Questionnaires{
				Title:       "第1回集会らん☆ぷろ募集アンケート",
				Description: "第1回集会らん☆ぷろ参加者募集",
			},
			administrators: []string{},
		},
	}

	for i, questionnaireData := range questionnaireDatas {
		err := db.Create(&questionnaireDatas[i].Questionnaires).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %w", questionnaireData, err)
		}

		for _, administrator := range questionnaireData.administrators {
			err = db.Create(&Administrators{
				QuestionnaireID: questionnaireDatas[i].Questionnaires.ID,
				UserTraqid:      administrator,
			}).Error
			if err != nil {
				t.Errorf("failed to create administrator(%s): %w", administrator, err)
			}
		}
	}

	datas = []QuestionsTestData{
		{
			Question: Question{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     1,
				Type:            "Text",
				Body:            "",
			},
		},
		{
			Question: Question{
				Type: "TextArea",
			},
		},
		{
			Question: Question{
				Type: "Number",
			},
		},
		{
			Question: Question{
				Type: "MultipleChoice",
			},
		},
		{
			Question: Question{
				Type: "Checkbox",
			},
		},
		{
			Question: Question{
				Type: "Dropdown",
			},
		},
		{
			Question: Question{
				Type: "LinearScale",
			},
		},
		{
			Question: Question{
				Type: "Date",
			},
		},
		{
			Question: Question{
				Type: "Time",
			},
		},
	}
}

func insertQuestionTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	//assertion := assert.New(t)

	type args struct {
		questionnaireID int
		pageNum         int
		questionNum     int
		questionType    string
		body            string
		isRequired      bool
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
}
