package model

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestInsertScaleLabel(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validID         bool
		ScaleLabelRight string
		ScaleLabelLeft  string
		ScaleMin        int
		ScaleMax        int
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
			description: "valid",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        5,
			},
		},
		{
			description: "long scalelabel",
			args: args{
				validID:         true,
				ScaleLabelRight: strings.Repeat("a", 2000),
				ScaleLabelLeft:  strings.Repeat("a", 2000),
				ScaleMin:        0,
				ScaleMax:        5,
			},
		},
		{
			description: "too long scalelabel",
			args: args{
				validID:         true,
				ScaleLabelRight: strings.Repeat("a", 200000),
				ScaleLabelLeft:  strings.Repeat("a", 200000),
				ScaleMin:        0,
				ScaleMax:        5,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "large scaleMax/Min",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        int(math.Pow10(5)),
			},
		},
		{
			description: "too large scaleMax/Min",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        int(math.Pow10(10)),
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "question does not exist",
			args: args{
				validID:         false,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        5,
			},
			expect: expect{
				isErr: true,
			},
		},
	}
	for _, testCase := range testCases {
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "LinearScale", "Linear", true)
		require.NoError(t, err)
		if !testCase.args.validID {
			questionID = -1
		}

		label := ScaleLabels{
			QuestionID:      questionID,
			ScaleLabelRight: testCase.args.ScaleLabelRight,
			ScaleLabelLeft:  testCase.args.ScaleLabelLeft,
			ScaleMin:        testCase.args.ScaleMin,
			ScaleMax:        testCase.args.ScaleMax,
		}

		err = InsertScaleLabel(questionID, label)
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

		label = ScaleLabels{}
		err = db.Where("question_id = ?", questionID).First(&label).Error
		assertion.NoError(err, testCase.description, "get scalelabels")

		assertion.Equal(questionID, label.QuestionID, testCase.description, "questionID")
		assertion.Equal(testCase.args.ScaleLabelRight, label.ScaleLabelRight, testCase.description, "ScaleLabelRight")
		assertion.Equal(testCase.args.ScaleLabelLeft, label.ScaleLabelLeft, testCase.description, "ScaleLabelLeft")
		assertion.Equal(testCase.args.ScaleMin, label.ScaleMin, testCase.description, "ScaleMin")
		assertion.Equal(testCase.args.ScaleMax, label.ScaleMax, testCase.description, "ScaleMax")
	}
}

func TestUpdateScaleLabel(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validID         bool
		ScaleLabelRight string
		ScaleLabelLeft  string
		ScaleMin        int
		ScaleMax        int
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
			description: "valid",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right_updated",
				ScaleLabelLeft:  "Left_updated",
				ScaleMin:        1,
				ScaleMax:        6,
			},
		},
		{
			description: "no update",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        5,
			},
		},
		{
			description: "question does not exist",
			args: args{
				validID:         false,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        5,
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordUpdated,
			},
		},
	}
	for _, testCase := range testCases {
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "LinearScale", "Linear", true)
		require.NoError(t, err)

		label := ScaleLabels{
			QuestionID:      questionID,
			ScaleLabelRight: "Right",
			ScaleLabelLeft:  "Left",
			ScaleMin:        0,
			ScaleMax:        5,
		}

		err = InsertScaleLabel(questionID, label)
		require.NoError(t, err)

		if !testCase.args.validID {
			questionID = -1
		}

		label = ScaleLabels{
			QuestionID:      questionID,
			ScaleLabelRight: testCase.args.ScaleLabelRight,
			ScaleLabelLeft:  testCase.args.ScaleLabelLeft,
			ScaleMin:        testCase.args.ScaleMin,
			ScaleMax:        testCase.args.ScaleMax,
		}

		err = UpdateScaleLabel(questionID, label)

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

		label = ScaleLabels{}
		err = db.Where("question_id = ?", questionID).First(&label).Error
		assertion.NoError(err, testCase.description, "get scalelabels")

		assertion.Equal(questionID, label.QuestionID, testCase.description, "questionID")
		assertion.Equal(testCase.args.ScaleLabelRight, label.ScaleLabelRight, testCase.description, "ScaleLabelRight")
		assertion.Equal(testCase.args.ScaleLabelLeft, label.ScaleLabelLeft, testCase.description, "ScaleLabelLeft")
		assertion.Equal(testCase.args.ScaleMin, label.ScaleMin, testCase.description, "ScaleMin")
		assertion.Equal(testCase.args.ScaleMax, label.ScaleMax, testCase.description, "ScaleMax")
	}
}

func TestDeleteScaleLabel(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validID         bool
		ScaleLabelRight string
		ScaleLabelLeft  string
		ScaleMin        int
		ScaleMax        int
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
			description: "valid",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right_updated",
				ScaleLabelLeft:  "Left_updated",
				ScaleMin:        1,
				ScaleMax:        6,
			},
		},
		{
			description: "question does not exist",
			args: args{
				validID:         false,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        5,
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordDeleted,
			},
		},
	}
	for _, testCase := range testCases {
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "LinearScale", "Linear", true)
		require.NoError(t, err)

		label := ScaleLabels{
			QuestionID:      questionID,
			ScaleLabelRight: testCase.args.ScaleLabelRight,
			ScaleLabelLeft:  testCase.args.ScaleLabelLeft,
			ScaleMin:        testCase.args.ScaleMin,
			ScaleMax:        testCase.args.ScaleMax,
		}

		err = InsertScaleLabel(questionID, label)
		require.NoError(t, err)

		if !testCase.args.validID {
			questionID = -1
		}

		err = DeleteScaleLabel(questionID)

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
	}
}

func TestGetScaleLabels(t *testing.T) {
	t.Parallel()
	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		questionIDs []int
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
	labels := []ScaleLabels{
		{
			ScaleLabelRight: "Right1",
			ScaleLabelLeft:  "Left1",
			ScaleMin:        0,
			ScaleMax:        5,
		},
		{
			ScaleLabelRight: "Right2",
			ScaleLabelLeft:  "Left2",
			ScaleMin:        0,
			ScaleMax:        5,
		},
		{
			ScaleLabelRight: "Right3",
			ScaleLabelLeft:  "Left3",
			ScaleMin:        0,
			ScaleMax:        5,
		},
	}
	questionIDs := make([]int, 0, 3)
	labelMap := make(map[int]ScaleLabels)
	for _, label := range labels {
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "LinearScale", "Linear", true)
		require.NoError(t, err)
		err = InsertScaleLabel(questionID, label)
		require.NoError(t, err)
		label.QuestionID = questionID
		questionIDs = append(questionIDs, questionID)
		labelMap[questionID] = label
	}

	testCases := []test{
		{
			description: "valid",
			args: args{
				questionIDs: questionIDs,
			},
			expect: expect{
				length: len(questionIDs),
			},
		},
		{
			description: "empty list",
			args: args{
				questionIDs: []int{},
			},
			expect: expect{
				length: 0,
			},
		},
		{
			description: "not exist",
			args: args{
				questionIDs: []int{-1},
			},
			expect: expect{
				length: 0,
			},
		},
	}

	for _, testCase := range testCases {

		labels, err := GetScaleLabels(testCase.args.questionIDs)
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
		assertion.Equal(testCase.expect.length, len(labels), testCase.description, "length")
		if len(labels) < 1 {
			continue
		}

		for _, label := range labels {
			expectLabel, ok := labelMap[label.QuestionID]
			require.Equal(t, true, ok)

			assertion.Equal(expectLabel.QuestionID, label.QuestionID, testCase.description, "questionID")
			assertion.Equal(expectLabel.ScaleLabelRight, label.ScaleLabelRight, testCase.description, "ScaleLabelRight")
			assertion.Equal(expectLabel.ScaleLabelLeft, label.ScaleLabelLeft, testCase.description, "ScaleLabelLeft")
			assertion.Equal(expectLabel.ScaleMin, label.ScaleMin, testCase.description, "ScaleMin")
			assertion.Equal(expectLabel.ScaleMax, label.ScaleMax, testCase.description, "ScaleMax")
		}
	}
}

func TestCheckScaleLabel(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	questionID, err := InsertQuestion(questionnaireID, 1, 1, "LinearScale", "Linear", true)
	require.NoError(t, err)

	label := ScaleLabels{
		QuestionID:      questionID,
		ScaleLabelRight: "Right",
		ScaleLabelLeft:  "Left",
		ScaleMin:        0,
		ScaleMax:        5,
	}

	err = InsertScaleLabel(questionID, label)
	require.NoError(t, err)

	type args struct {
		label    ScaleLabels
		response string
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
			description: "valid",
			args: args{
				label:    label,
				response: "2",
			},
		},
		{
			description: "too large response",
			args: args{
				label:    label,
				response: "10",
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "too short response",
			args: args{
				label:    label,
				response: "-1",
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response",
			args: args{
				label:    label,
				response: "a",
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "empty response",
			args: args{
				label:    label,
				response: "",
			},
		},
	}
	for _, testCase := range testCases {

		err = CheckScaleLabel(testCase.args.label, testCase.args.response)

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
	}
}
