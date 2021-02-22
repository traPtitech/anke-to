package model

import (
	"sort"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
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
	questionnaireDatas       = []QuestionsTestQuestionnaireData{}
	questionDatas            = []QuestionsTestData{}
	questionQuestionnaireMap = map[int][]*Questions{}
)

func TestQuestions(t *testing.T) {
	t.Parallel()

	setupQuestionsTest(t)

	t.Run("InsertQuestion", insertQuestionTest)
	t.Run("UpdateQuestion", updateQuestionTest)
	t.Run("DeleteQuestion", deleteQuestionTest)
	t.Run("GetQuestions", getQuestionsTest)
	t.Run("CheckQuestionAdmin", checkQuestionAdminTest)
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
				QuestionNum:     0,
				Type:            "Text",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     1,
				Type:            "TextArea",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     2,
				Type:            "Number",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     3,
				Type:            "MultipleChoice",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     4,
				Type:            "Checkbox",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     5,
				Type:            "Dropdown",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     6,
				Type:            "LinearScale",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     7,
				Type:            "Date",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     8,
				Type:            "Time",
				Body:            "",
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[1].ID,
				PageNum:         1,
				QuestionNum:     8,
				Type:            "Text",
				Body:            "",
				DeletedAt: mysql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			},
		},
		{
			Questions: Questions{
				QuestionnaireID: questionnaireDatas[0].ID,
				PageNum:         1,
				QuestionNum:     0,
				Type:            "Text",
				Body:            "",
			},
		},
	}

	for i, questionData := range questionDatas {
		err := db.Create(&questionDatas[i].Questions).Error
		if err != nil {
			t.Errorf("failed to create questionnaire(%+v): %w", questionData, err)
		}

		if !questionData.Questions.DeletedAt.Valid {
			questions, ok := questionQuestionnaireMap[questionData.Questions.QuestionnaireID]
			if !ok {
				questionQuestionnaireMap[questionData.Questions.QuestionnaireID] = []*Questions{}
			}
			questionQuestionnaireMap[questionData.Questions.QuestionnaireID] = append(questions, &questionDatas[i].Questions)
		}
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
			description: "type:TextArea, required: false, question_num: 0",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         1,
					QuestionNum:     0,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
				},
			},
		},
		{
			description: "type:TextArea, required: false, page_num: 0",
			args: args{
				Questions: Questions{
					QuestionnaireID: questionnaireDatas[0].ID,
					PageNum:         0,
					QuestionNum:     1,
					Type:            "TextArea",
					Body:            "自由記述欄",
					IsRequired:      false,
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
		after `json:"-"`
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
			description: "questionNum: 1->0",
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
					QuestionNum:     0,
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

func deleteQuestionTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionID int
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

	invalidQuestionID := 1000
	for {
		err := db.Where("id = ?", invalidQuestionID).First(&Questions{}).Error
		if gorm.IsRecordNotFoundError(err) {
			break
		}
		if err != nil {
			t.Errorf("failed to get questionnaire(make invalid questionnaireID): %w", err)
			break
		}

		invalidQuestionID *= 10
	}

	testQuestions := []*Questions{
		{
			QuestionnaireID: questionnaireDatas[0].ID,
			PageNum:         1,
			QuestionNum:     1,
			Type:            "TextArea",
			Body:            "自由記述欄",
			IsRequired:      false,
		},
	}

	for _, question := range testQuestions {
		err := db.Create(question).Error
		if err != nil {
			t.Errorf("failed to insert question: %w", err)
		}
	}

	testCases := []test{
		{
			description: "questionID: valid",
			args: args{
				questionID: testQuestions[0].ID,
			},
		},
		{
			description: "questionID: invalid",
			args: args{
				questionID: invalidQuestionID,
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordDeleted,
			},
		},
	}

	for _, testCase := range testCases {
		err := questionImpl.DeleteQuestion(testCase.args.questionID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		actualQuestion := Questions{}
		err = db.Unscoped().Where("id = ?", testCase.args.questionID).First(&actualQuestion).Error
		if err != nil {
			t.Errorf("failed to get question(%s): %w", testCase.description, err)
		}

		assertion.True(actualQuestion.DeletedAt.Valid, testCase.description, "deleted_at")
	}
}

func getQuestionsTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionnaireID int
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
			description: "questionID: valid",
			args: args{
				questionnaireID: questionnaireDatas[1].ID,
			},
		},
		{
			description: "questionID: invalid",
			args: args{
				questionnaireID: invalidQuestionnaireID,
			},
		},
	}

	for _, testCase := range testCases {
		questions, err := questionImpl.GetQuestions(testCase.args.questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(testCase.expect.err, err, testCase.description, "error")
		}
		if err != nil {
			continue
		}

		actualQuestionIDs := make([]int, 0, len(questions))
		for _, question := range questions {
			actualQuestionIDs = append(actualQuestionIDs, question.ID)
		}

		expectQuestions := questionQuestionnaireMap[testCase.questionnaireID]
		if expectQuestions == nil {
			expectQuestions = []*Questions{}
		}
		expectQuestionIDs := make([]int, 0, len(expectQuestions))
		for _, question := range expectQuestions {
			expectQuestionIDs = append(expectQuestionIDs, question.ID)
		}

		assertion.Subset(actualQuestionIDs, expectQuestionIDs, testCase.description, "elements")

		assertion.True(sort.SliceIsSorted(questions, func(i, j int) bool { return questions[i].QuestionNum <= questions[j].QuestionNum }), testCase.description, "sort")

		for i, actualQuestion := range questions {
			expectQuestion := expectQuestions[i]
			assertion.Equal(expectQuestion.QuestionnaireID, actualQuestion.QuestionnaireID, testCase.description, "questionnaire_id")
			assertion.Equal(expectQuestion.PageNum, actualQuestion.PageNum, testCase.description, "page_num")
			assertion.Equal(expectQuestion.QuestionNum, actualQuestion.QuestionNum, testCase.description, "question_num")
			assertion.Equal(expectQuestion.Type, actualQuestion.Type, testCase.description, "type")
			assertion.Equal(expectQuestion.Body, actualQuestion.Body, testCase.description, "body")
			assertion.Equal(expectQuestion.IsRequired, actualQuestion.IsRequired, testCase.description, "is_required")

			assertion.WithinDuration(expectQuestion.CreatedAt, actualQuestion.CreatedAt, time.Second, testCase.description, "created_at")
			assertion.Equal(false, actualQuestion.DeletedAt.Valid, testCase.description, "deleted_at")
		}
	}
}

func checkQuestionAdminTest(t *testing.T) {
	t.Helper()
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		userID     string
		questionID int
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
	testCases := []test{
		{
			description: "userID: valid, admin: true",
			args: args{
				userID:     questionsTestUserID,
				questionID: questionDatas[10].ID,
			},
			expect: expect{
				isAdmin: true,
			},
		},
		{
			description: "userID: valid, admin: false",
			args: args{
				userID:     questionsTestUserID,
				questionID: questionDatas[0].ID,
			},
			expect: expect{
				isAdmin: false,
			},
		},
		{
			description: "userID: invalid, admin: true",
			args: args{
				userID:     invalidQuestionsTestUserID,
				questionID: questionDatas[10].ID,
			},
			expect: expect{
				isAdmin: false,
			},
		},
		{
			description: "userID: invalid, admin: false",
			args: args{
				userID:     invalidQuestionsTestUserID,
				questionID: questionDatas[10].ID,
			},
			expect: expect{
				isAdmin: false,
			},
		},
	}

	for _, testCase := range testCases {
		actualIsAdmin, err := questionImpl.CheckQuestionAdmin(testCase.args.userID, testCase.args.questionID)

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
