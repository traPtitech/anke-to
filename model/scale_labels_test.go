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

	questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = InsertAdministrators(questionnaireID, []string{userOne})
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
			description: "too long scalelabel",
			args: args{
				validID:         true,
				ScaleLabelRight: "Right",
				ScaleLabelLeft:  "Left",
				ScaleMin:        0,
				ScaleMax:        int(math.Pow10(11)),
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
