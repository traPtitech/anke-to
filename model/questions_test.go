package model

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

const questionsTestUserID = "questionsUser"
const invalidQuestionsTestUserID = "invalidQuestionsUser"

type QuestionsTestQuestionnaireData struct {
	Questionnaires
	administrators []string
}

type QuestionsTestData struct {
	Questions
}

var (
	questionnaireDatas = []QuestionsTestQuestionnaireData{}
	questionDatas      = []QuestionsTestData{}
)

func TestQuestions(t *testing.T) {
	t.Parallel()

	setupQuestionsTest(t)

	t.Run("InsertQuestion", insertQuestionTest)
	t.Run("UpdateQuestion", updateQuestionTest)
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

	questionDatas = []QuestionsTestData{
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     1,
				Type:            "Text",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				Type: "TextArea",
			},
		},
		{
			Questions: Questions{
				Type: "Number",
			},
		},
		{
			Questions: Questions{
				Type: "MultipleChoice",
			},
		},
		{
			Questions: Questions{
				Type: "Checkbox",
			},
		},
		{
			Questions: Questions{
				Type: "Dropdown",
			},
		},
		{
			Questions: Questions{
				Type: "LinearScale",
			},
		},
		{
			Questions: Questions{
				Type: "Date",
			},
		},
		{
			Questions: Questions{
				Type: "Time",
			},
		},
	}
}

func insertQuestionTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		Questions
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
			description: "type:TextArea, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:Number, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "Number",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:MultipleChoice, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "MultipleChoice",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:Checkbox, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "Checkbox",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:Dropdown, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "Dropdown",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:LinearScale, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "LinearScale",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:Date, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "Date",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:Time, required: false",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "Time",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:TextArea, required: true",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      true,
				},
			},
		},
		{
			description: "invalid questionnaireID",
			args: args{
				Questions: Questions{
					QuestionnaireID: invalidQuestionnaireID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		createdAt := time.Now()
		questionID, err := questionImpl.InsertQuestion(testCase.args.QuestionnaireID, testCase.args.PageNum, testCase.args.QuestionNum, testCase.args.Type, testCase.args.Body, testCase.args.IsRequired)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		question := Questions{}
		err = db.Where("id = ?", questionID).First(&question).Error
		if err != nil {
			t.Errorf("failed to get question(%s): %w", testCase.description, err)
		}

		assertion.Equal(testCase.args.QuestionnaireID, question.QuestionnaireID, testCase.description, "questionnaire_id")
		assertion.Equal(testCase.args.PageNum, question.PageNum, testCase.description, "page_num")
		assertion.Equal(testCase.args.QuestionNum, question.QuestionNum, testCase.description, "question_num")
		assertion.Equal(testCase.args.Type, question.Type, testCase.description, "type")
		assertion.Equal(testCase.args.Body, question.Body, testCase.description, "body")
		assertion.Equal(testCase.args.IsRequired, question.IsRequired, testCase.description, "is_required")

		assertion.WithinDuration(createdAt, question.CreatedAt, 2*time.Second, testCase.description, "created_at")
		assertion.Equal(false, question.DeletedAt.Valid, testCase.description, "deleted_at")
	}
}

func updateQuestionTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type before struct {
		Questions
	}
	type after struct {
		Questions
	}
	type expect struct {
		isErr bool
		err   error
	}
	type test struct {
		description string
		before
		after
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
			description: "type:TextArea->Number",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "Number",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "questionnaireID: valid->valid",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[1].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "questionnaireID: valid->invalid",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: invalidQuestionnaireID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "pageNum: 1->2",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         2,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "questionNum: 1->2",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     2,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "body: 自由記述欄->自由記述欄1",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     2,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "isRequired: false->true",
			before: before{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
			after: after{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     2,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      true,
				},
			},
		},
	}

	for _, testCase := range testCases {
		question := &testCase.before.Questions
		err := db.Create(question).Error
		if err != nil {
			t.Errorf("failed to insert question(%s): %w", testCase.description, err)
		}

		err = questionImpl.UpdateQuestion(testCase.after.QuestionnaireID, testCase.after.PageNum, testCase.after.QuestionNum, testCase.after.Type, testCase.after.Body, testCase.after.IsRequired, question.ID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		actualQuestion := Questions{}
		err = db.Where("id = ?", question.ID).First(&actualQuestion).Error
		if err != nil {
			t.Errorf("failed to get question(%s): %w", testCase.description, err)
		}

		assertion.Equal(testCase.after.QuestionnaireID, actualQuestion.QuestionnaireID, testCase.description, "questionnaire_id")
		assertion.Equal(testCase.after.PageNum, actualQuestion.PageNum, testCase.description, "page_num")
		assertion.Equal(testCase.after.QuestionNum, actualQuestion.QuestionNum, testCase.description, "question_num")
		assertion.Equal(testCase.after.Type, actualQuestion.Type, testCase.description, "type")
		assertion.Equal(testCase.after.Body, actualQuestion.Body, testCase.description, "body")
		assertion.Equal(testCase.after.IsRequired, actualQuestion.IsRequired, testCase.description, "is_required")

		assertion.WithinDuration(question.CreatedAt, question.CreatedAt, time.Second, testCase.description, "created_at")
		assertion.Equal(false, question.DeletedAt.Valid, testCase.description, "deleted_at")
	}
}
