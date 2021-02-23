package model

import "testing"

type optionsTestQuestionnaireData struct {
	questionnaire Questionnaires
	questions     []*optionsTestQuestionData
}

type optionsTestQuestionData struct {
	question Questions
	options  []*Options
}

var (
	optionsTestQuestionnaireDatas []*optionsTestQuestionnaireData
)

func TestOptions(t *testing.T) {
	t.Parallel()

	setupOptionsTest(t)
}

func setupOptionsTest(t *testing.T) {
	optionsTestQuestionnaireDatas = []*optionsTestQuestionnaireData{
		{
			questionnaire: Questionnaires{
				Title:       "第1回集会らん☆ぷろ募集アンケート",
				Description: "第1回集会らん☆ぷろ参加者募集",
			},
			questions: []*optionsTestQuestionData{
				{
					question: Questions{
						PageNum:     1,
						QuestionNum: 0,
						Type:        "MultipleChoice",
						Body:        "",
					},
					options: []*Options{
						{
							OptionNum: 0,
							Body:      "option1",
						},
						{
							OptionNum: 1,
							Body:      "option2",
						},
					},
				},
				{
					question: Questions{
						PageNum:     1,
						QuestionNum: 1,
						Type:        "MultipleChoice",
						Body:        "",
					},
					options: []*Options{},
				},
			},
		},
	}

	for _, questionnaireData := range optionsTestQuestionnaireDatas {
		err := db.Create(&questionnaireData.questionnaire).Error
		if err != nil {
			t.Errorf("failed to create questionnaire: %w", err)
		}

		for _, questionData := range questionnaireData.questions {
			questionData.question.QuestionnaireID = questionnaireData.questionnaire.ID
			err := db.Create(&questionData.question).Error
			if err != nil {
				t.Errorf("failed to create question: %w", err)
			}

			for _, option := range questionData.options {
				option.QuestionID = questionData.question.ID
				err := db.Create(option).Error
				if err != nil {
					t.Errorf("failed to create option: %w", err)
				}
			}
		}
	}
}
