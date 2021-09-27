package model

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
)

func TestInsertRespondent(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validQuestionnaireID bool
		userID               string
		submittedAt          null.Time
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
			description: "valid submittedAt not null",
			args: args{
				validQuestionnaireID: true,
				userID:               userTwo,
				submittedAt:          null.NewTime(time.Now(), true),
			},
		},
		{
			description: "valid submittedAt null",
			args: args{
				validQuestionnaireID: true,
				userID:               userTwo,
				submittedAt:          null.NewTime(time.Time{}, false),
			},
		},
		{
			description: "long userID",
			args: args{
				validQuestionnaireID: true,
				userID:               strings.Repeat("a", 30),
				submittedAt:          null.NewTime(time.Time{}, false),
			},
		},
		{
			description: "too long userID",
			args: args{
				validQuestionnaireID: true,
				userID:               strings.Repeat("a", 100),
				submittedAt:          null.NewTime(time.Time{}, false),
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "questionnaireID does not exist",
			args: args{
				validQuestionnaireID: false,
				userID:               userTwo,
				submittedAt:          null.NewTime(time.Time{}, false),
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		if !testCase.args.validQuestionnaireID {
			questionnaireID = -1
		}

		responseID, err := respondentImpl.InsertRespondent(testCase.args.userID, questionnaireID, testCase.args.submittedAt)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		respondent := Respondents{}
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("response_id = ?", responseID).
			First(&respondent).Error
		assertion.NoError(err, testCase.description, "get respondent")

		assertion.Equal(responseID, respondent.ResponseID, testCase.description, "responseID")
		assertion.Equal(questionnaireID, respondent.QuestionnaireID, testCase.description, "questionnaireID")
		assertion.Equal(testCase.args.userID, respondent.UserTraqid, testCase.description, "userID")
		assertion.WithinDuration(testCase.args.submittedAt.ValueOrZero(), respondent.SubmittedAt.ValueOrZero(), 2*time.Second, testCase.description, "submittedAt")
		assertion.WithinDuration(time.Now(), respondent.ModifiedAt, 2*time.Second, testCase.description, "modified_at")
		assertion.WithinDuration(null.NewTime(time.Time{}, false).ValueOrZero(), respondent.DeletedAt.Time, 2*time.Second, testCase.description, "deleted_at")
	}
}

func TestUpdateSubmittedAt(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validresponseID bool
		userID          string
		submittedAt     null.Time
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
			description: "valid submitted not null",
			args: args{
				validresponseID: true,
				userID:          userTwo,
				submittedAt:     null.NewTime(time.Now(), true),
			},
		},
	}

	for _, testCase := range testCases {
		responseID, err := respondentImpl.InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), false))
		require.NoError(t, err)
		if !testCase.args.validresponseID {
			responseID = -1
		}

		err = respondentImpl.UpdateSubmittedAt(responseID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		respondent := Respondents{}
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("response_id = ?", responseID).
			First(&respondent).Error
		assertion.NoError(err, testCase.description, "get respondent")

		assertion.Equal(responseID, respondent.ResponseID, testCase.description, "responseID")
		assertion.Equal(questionnaireID, respondent.QuestionnaireID, testCase.description, "questionnaireID")
		assertion.Equal(testCase.args.userID, respondent.UserTraqid, testCase.description, "userID")
		assertion.WithinDuration(testCase.args.submittedAt.ValueOrZero(), respondent.SubmittedAt.ValueOrZero(), 2*time.Second, testCase.description, "submittedAt")
		assertion.WithinDuration(time.Now(), respondent.ModifiedAt, 2*time.Second, testCase.description, "modified_at")
		assertion.WithinDuration(null.NewTime(time.Time{}, false).ValueOrZero(), respondent.DeletedAt.Time, 2*time.Second, testCase.description, "deleted_at")
	}
}

func TestDeleteRespondent(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validresponseID bool
		insertUserID    string
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
			description: "delete by respondents",
			args: args{
				validresponseID: true,
				insertUserID:    userTwo,
			},
		},
		{
			description: "responseID does not exist",
			args: args{
				validresponseID: false,
				insertUserID:    userTwo,
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordDeleted,
			},
		},
	}

	for _, testCase := range testCases {
		responseID, err := respondentImpl.InsertRespondent(testCase.args.insertUserID, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)
		if !testCase.args.validresponseID {
			responseID = -1
		}

		err = respondentImpl.DeleteRespondent(responseID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		respondent := Respondents{}
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Unscoped().
			Where("response_id = ?", responseID).
			First(&respondent).Error
		if err != nil {
			t.Errorf("failed to get respondent: %v", err)
		}

		assertion.WithinDuration(time.Now(), respondent.DeletedAt.Time, 2*time.Second, testCase.description, "deleted_at")
	}
}

func TestGetRespondentInfos(t *testing.T) {
	t.Parallel()
	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		questionnaireIDs []int
		userID           string
	}
	type expect struct {
		isErr  bool
		err    error
		length int
	}
	type test struct {
		description string
		args
		expect
	}
	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第2回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)
	questionnaireID2, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第2回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	questionnaire := Questionnaires{}
	err = db.
		Session(&gorm.Session{NewDB: true}).
		Unscoped().
		Where("id = ?", questionnaireID).
		Find(&questionnaire).Error
	require.NoError(t, err)

	respondents := []Respondents{
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userOne,
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userTwo,
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userTwo,
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
		{
			QuestionnaireID: questionnaireID2,
			UserTraqid:      "TestGetRespondentInfos",
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
	}

	respondentMap := make(map[int]Respondents)
	for _, respondent := range respondents {
		responseID, err := respondentImpl.InsertRespondent(respondent.UserTraqid, respondent.QuestionnaireID, respondent.SubmittedAt)
		require.NoError(t, err)
		respondent.ResponseID = responseID
		respondentMap[responseID] = respondent
	}

	testCases := []test{
		{
			description: "valid",
			args: args{
				userID:           userTwo,
				questionnaireIDs: []int{questionnaireID},
			},
			expect: expect{
				length: 2,
			},
		},
		{
			// downしてから
			description: "empty questionnaireIDs",
			args: args{
				userID:           "TestGetRespondentInfos",
				questionnaireIDs: []int{},
			},
			expect: expect{
				length: 1,
			},
		},
		{
			description: "no user",
			args: args{
				userID:           "test_user",
				questionnaireIDs: []int{questionnaireID},
			},
			expect: expect{
				length: 0,
			},
		},
		{
			description: "invalid questionnaireID",
			args: args{
				userID:           userTwo,
				questionnaireIDs: []int{-1},
			},
			expect: expect{
				length: 0,
			},
		},
	}

	for _, testCase := range testCases {

		respondentInfos, err := respondentImpl.GetRespondentInfos(testCase.args.userID, testCase.args.questionnaireIDs...)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}
		assertion.Equal(testCase.expect.length, len(respondentInfos), testCase.description, "length")
		if len(respondentInfos) < 1 {
			continue
		}

		for _, respondentInfo := range respondentInfos {
			expectRespondent, ok := respondentMap[respondentInfo.ResponseID]
			require.Equal(t, true, ok)
			if len(testCase.args.questionnaireIDs) > 0 {
				assertion.Equal(testCase.args.questionnaireIDs[0], respondentInfo.QuestionnaireID, testCase.description, "questionnaireID")
			} else {
				assertion.Equal(questionnaireID2, respondentInfo.QuestionnaireID, testCase.description, "questionnaireID")
			}
			assertion.Equal(expectRespondent.ResponseID, respondentInfo.ResponseID, testCase.description, "responseID")
			assertion.Equal(questionnaire.Title, respondentInfo.Title, testCase.description, "title")
			assertion.WithinDuration(questionnaire.ResTimeLimit.ValueOrZero(), respondentInfo.ResTimeLimit.ValueOrZero(), 2*time.Second, testCase.description, "ResTimeLimit")
			assertion.WithinDuration(expectRespondent.SubmittedAt.ValueOrZero(), respondentInfo.SubmittedAt.ValueOrZero(), 2*time.Second, testCase.description, "submittedAt")
			assertion.WithinDuration(time.Now(), respondentInfo.ModifiedAt, 2*time.Second, testCase.description, "modified_at")
			assertion.WithinDuration(null.NewTime(time.Time{}, false).ValueOrZero(), respondentInfo.DeletedAt.Time, 2*time.Second, testCase.description, "deleted_at")
		}
	}
}

func TestGetRespondentDetail(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	questionnaire := Questionnaires{}
	err = db.
		Session(&gorm.Session{NewDB: true}).
		Unscoped().
		Where("id = ?", questionnaireID).
		Find(&questionnaire).Error
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validresponseID bool
		responseMetas   []*ResponseMeta
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

	questionIDs := make([]int, 0, 2)

	questionID, err := questionImpl.InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
	require.NoError(t, err)
	questionIDs = append(questionIDs, questionID)

	questionID, err = questionImpl.InsertQuestion(questionnaireID, 1, 3, "MultipleChoice", "radio", true)
	require.NoError(t, err)
	questionIDs = append(questionIDs, questionID)

	testCases := []test{
		{
			description: "valid",
			args: args{
				validresponseID: true,
				responseMetas: []*ResponseMeta{
					{QuestionID: questionIDs[0], Data: "リマインダーBOTを作った話"},
					{QuestionID: questionIDs[1], Data: "選択肢1"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		responseID, err := respondentImpl.InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), false))
		require.NoError(t, err)
		if !testCase.args.validresponseID {
			responseID = -1
		} else {
			err := responseImpl.InsertResponses(responseID, testCase.args.responseMetas)
			require.NoError(t, err)
		}

		respondentDetail, err := respondentImpl.GetRespondentDetail(responseID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		assertion.Equal(questionnaireID, respondentDetail.QuestionnaireID, testCase.description, "questionnaireID")

		questionID := questionIDs[0]
		responseBody := respondentDetail.Responses[0]
		assertion.Equal(questionID, responseBody.QuestionID, testCase.description, "questionID1")
		assertion.Equal("Text", responseBody.QuestionType, testCase.description, "QuestionType1")
		assertion.Equal("リマインダーBOTを作った話", responseBody.Body.String, testCase.description, "description1")

		questionID = questionIDs[1]
		responseBody = respondentDetail.Responses[1]
		assertion.Equal(1, len(responseBody.OptionResponse), testCase.description, "OptionResponse len")
		optionResponse := responseBody.OptionResponse[0]
		assertion.Equal(questionID, responseBody.QuestionID, testCase.description, "questionID2")
		assertion.Equal("MultipleChoice", responseBody.QuestionType, testCase.description, "QuestionType2")
		assertion.Equal("選択肢1", optionResponse, testCase.description, "description2")
	}
}

func TestGetRespondentDetails(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	questionnaire := Questionnaires{}
	err = db.
		Session(&gorm.Session{NewDB: true}).
		Unscoped().
		Where("id = ?", questionnaireID).
		Find(&questionnaire).Error
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		questionnaireID int
		sort            string
	}
	type expect struct {
		isErr   bool
		err     error
		length  int
		sortIdx []int
	}
	type test struct {
		description string
		args
		expect
	}
	questions := []Questions{
		{
			QuestionnaireID: questionnaireID,
			PageNum:         1,
			QuestionNum:     1,
			Type:            "Text",
			Body:            "質問文",
			IsRequired:      true,
		},
		{
			QuestionnaireID: questionnaireID,
			PageNum:         1,
			QuestionNum:     2,
			Type:            "MultipleChoice",
			Body:            "radio",
			IsRequired:      true,
		},
		{
			QuestionnaireID: questionnaireID,
			PageNum:         1,
			QuestionNum:     3,
			Type:            "Number",
			Body:            "number",
			IsRequired:      true,
		},
	}

	questionLength := len(questions)
	questionIDs := make([]int, 0, questionLength)

	for _, question := range questions {
		questionID, err := questionImpl.InsertQuestion(question.QuestionnaireID, question.PageNum, question.QuestionNum, question.Type, question.Body, question.IsRequired)
		require.NoError(t, err)
		questionIDs = append(questionIDs, questionID)

	}

	respondents := []Respondents{
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userOne,
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userTwo,
			SubmittedAt:     null.NewTime(time.Now().Add(time.Second*3), true),
		},
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userThree,
			SubmittedAt:     null.NewTime(time.Now().Add(time.Second*2), true),
		},
		{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userOne,
			SubmittedAt:     null.NewTime(time.Now(), false),
		},
	}

	responseMetasList := [][]*ResponseMeta{
		{
			{QuestionID: questionIDs[0], Data: "リマインダーBOTを作った話1"},
			{QuestionID: questionIDs[1], Data: "選択肢1"},
			{QuestionID: questionIDs[2], Data: "10"},
		},
		{
			{QuestionID: questionIDs[0], Data: "リマインダーBOTを作った話2"},
			{QuestionID: questionIDs[1], Data: "選択肢2"},
			{QuestionID: questionIDs[2], Data: "5"},
		},
		{
			{QuestionID: questionIDs[0], Data: "リマインダーBOTを作った話3"},
			{QuestionID: questionIDs[1], Data: "選択肢3"},
			{QuestionID: questionIDs[2], Data: "0"},
		},
		{
			{QuestionID: questionIDs[0], Data: "リマインダーBOTを作った話1"},
			{QuestionID: questionIDs[1], Data: "選択肢1"},
			{QuestionID: questionIDs[2], Data: "10"},
		},
	}

	responseLength := len(respondents)
	responseIDs := make([]int, 0, responseLength)
	for i, respondent := range respondents {
		responseID, err := respondentImpl.InsertRespondent(respondent.UserTraqid, respondent.QuestionnaireID, respondent.SubmittedAt)
		require.NoError(t, err)
		responseIDs = append(responseIDs, responseID)

		err = responseImpl.InsertResponses(responseIDs[i], responseMetasList[i])
		require.NoError(t, err)

	}

	testCases := []test{
		{
			description: "traqid",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "traqid",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{0, 1, 2},
			},
		},
		{
			description: "-traqid",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "-traqid",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{2, 1, 0},
			},
		},
		{
			description: "submitted_at",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "submitted_at",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{0, 2, 1},
			},
		},
		{
			description: "-submitted_at",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "-submitted_at",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{1, 2, 0},
			},
		},
		{
			description: "questionnaire does not exist",
			args: args{
				questionnaireID: -1,
				sort:            "1",
			},
			expect: expect{
				length:  0,
				sortIdx: []int{},
			},
		},
		{
			description: "sortNum Number",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "3",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{2, 1, 0},
			},
		},
		{
			description: "sortNum Number",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "-3",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{0, 1, 2},
			},
		},
		{
			description: "sortNum Text",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "1",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{0, 1, 2},
			},
		},
		{
			description: "sortNum Text desc",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "-1",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{2, 1, 0},
			},
		},
		{
			description: "invalid sortnum",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "a",
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "empty sortnum",
			args: args{
				questionnaireID: questionnaireID,
				sort:            "",
			},
			expect: expect{
				length:  3,
				sortIdx: []int{0, 1, 2},
			},
		},
	}

	for _, testCase := range testCases {
		respondentDetails, err := respondentImpl.GetRespondentDetails(testCase.args.questionnaireID, testCase.args.sort)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.length, len(respondentDetails), testCase.description, "respondentDetails length")
		for i, respondentDetail := range respondentDetails {
			responseID := responseIDs[testCase.expect.sortIdx[i]]
			assertion.Equal(responseID, respondentDetail.ResponseID, testCase.description, "sort ID")
		}
	}
}

func TestGetRespondentsUserIDs(t *testing.T) {
	t.Parallel()
	assertion := assert.New(t)
	ctx := context.Background()

	type args struct {
		questionnaireIDs []int
	}
	type expect struct {
		isErr   bool
		err     error
		length  int
		userIDs []string
	}
	type test struct {
		description string
		args
		expect
	}
	questionnaireIDs := make([]int, 0, 3)
	for i := 0; i < 3; i++ {
		questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnaireID)
	}

	respondents := []Respondents{
		{
			QuestionnaireID: questionnaireIDs[0],
			UserTraqid:      userOne,
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
		{
			QuestionnaireID: questionnaireIDs[1],
			UserTraqid:      userTwo,
			SubmittedAt:     null.NewTime(time.Now(), true),
		},
	}

	respondentMap := make(map[int]Respondents)
	for _, respondent := range respondents {
		responseID, err := respondentImpl.InsertRespondent(respondent.UserTraqid, respondent.QuestionnaireID, respondent.SubmittedAt)
		require.NoError(t, err)
		respondent.ResponseID = responseID
		respondentMap[responseID] = respondent
	}

	testCases := []test{
		{
			description: "valid",
			args: args{
				questionnaireIDs: questionnaireIDs,
			},
			expect: expect{
				length:  2,
				userIDs: []string{userOne, userTwo},
			},
		},
		{
			description: "invalid questionnaireID",
			args: args{
				questionnaireIDs: []int{-1},
			},
			expect: expect{
				length: 0,
			},
		},
		{
			description: "nothing respondent questionnaireID",
			args: args{
				questionnaireIDs: []int{questionnaireIDs[2]},
			},
			expect: expect{
				length: 0,
			},
		},
	}

	for _, testCase := range testCases {

		respondents, err := respondentImpl.GetRespondentsUserIDs(testCase.args.questionnaireIDs)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}
		assertion.Equal(testCase.expect.length, len(respondents), testCase.description, "length")
		if len(respondents) < 1 {
			continue
		}

		for i, respondent := range respondents {
			assertion.Equal(testCase.expect.userIDs[i], respondent.UserTraqid, testCase.description, "userID")
		}
	}
}

func TestTestCheckRespondent(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	_, err = respondentImpl.InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), true))
	require.NoError(t, err)

	type args struct {
		userID          string
		questionnaireID int
	}
	type expect struct {
		isErr        bool
		err          error
		isRespondent bool
	}

	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "valid",
			args: args{
				userID:          userTwo,
				questionnaireID: questionnaireID,
			},
			expect: expect{
				isRespondent: true,
			},
		},
		{
			description: "not respondents",
			args: args{
				userID:          userThree,
				questionnaireID: questionnaireID,
			},
			expect: expect{
				isRespondent: false,
			},
		},
		{
			description: "questionnaireID does not exist",
			args: args{
				userID:          userTwo,
				questionnaireID: -1,
			},
			expect: expect{
				isRespondent: false,
			},
		},
	}

	for _, testCase := range testCases {
		isRespondent, err := respondentImpl.CheckRespondent(testCase.args.userID, testCase.args.questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.isRespondent, isRespondent, testCase.description, "isRespondent")
	}
}

func TestCheckRespondentByResponseID(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	responseID, err := respondentImpl.InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), true))
	require.NoError(t, err)

	type args struct {
		userID     string
		responseID int
	}
	type expect struct {
		isErr        bool
		err          error
		isRespondent bool
	}

	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "valid",
			args: args{
				userID:     userTwo,
				responseID: responseID,
			},
			expect: expect{
				isRespondent: true,
			},
		},
		{
			description: "not respondents",
			args: args{
				userID:     userThree,
				responseID: responseID,
			},
			expect: expect{
				isRespondent: false,
			},
		},
		{
			description: "questionnaireID does not exist",
			args: args{
				userID:     userTwo,
				responseID: -1,
			},
			expect: expect{
				isRespondent: false,
			},
		},
	}

	for _, testCase := range testCases {
		isRespondent, err := respondentImpl.CheckRespondentByResponseID(testCase.args.userID, testCase.args.responseID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else if testCase.expect.isErr {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.isRespondent, isRespondent, testCase.description, "isRespondent")
	}
}
