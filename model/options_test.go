package model

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateOptions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	questionnaire := Questionnaires{}
	err := db.
		Session(&gorm.Session{}).
		Create(&questionnaire).Error
	if err != nil {
		t.Errorf("failed to create Questionnaire: %v", err)
	}

	type test struct {
		description   string
		beforeOptions []Options
		afterOptions  []Options
		argOptions    []string
		isErr         bool
		err           error
	}

	/*Note:
	OptinonNumの重複はInsertの仕組み上できないのでテストから除外
	万が一入ってしまった場合現在の実装だと壊れる
	*/
	testCases := []test{
		{
			description: "変更なしなのでエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			argOptions: []string{"a"},
		},
		{
			description: "追加が1つあってもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
				{
					OptionNum: 2,
					Body:      "b",
				},
			},
			argOptions: []string{"a", "b"},
		},
		{
			description: "追加が複数あってもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
				{
					OptionNum: 2,
					Body:      "b",
				},
				{
					OptionNum: 3,
					Body:      "c",
				},
			},
			argOptions: []string{"a", "b", "c"},
		},
		{
			description: "ラベル変更が1つあってもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "b",
				},
			},
			argOptions: []string{"b"},
		},
		{
			description: "ラベル変更が複数あってもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
				{
					OptionNum: 2,
					Body:      "b",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "b",
				},
				{
					OptionNum: 2,
					Body:      "c",
				},
			},
			argOptions: []string{"b", "c"},
		},
		{
			description: "option_numが飛んでいて、飛んだoptionがいらなくてもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
				{
					OptionNum: 3,
					Body:      "c",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			argOptions: []string{"a"},
		},
		{
			description: "option_numが飛んでいて、飛んだoptionに変更が入ってもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
				{
					OptionNum: 3,
					Body:      "c",
				},
			},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
				{
					OptionNum: 2,
					Body:      "b",
				},
				{
					OptionNum: 3,
					Body:      "c",
				},
			},
			argOptions: []string{"a", "b", "c"},
		},
		{
			description:   "元々の選択肢がなくてもエラーなし",
			beforeOptions: []Options{},
			afterOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			argOptions: []string{"a"},
		},
		{
			description: "新たな選択肢がなくてもエラーなし",
			beforeOptions: []Options{
				{
					OptionNum: 1,
					Body:      "a",
				},
			},
			afterOptions: []Options{},
			argOptions:   []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			question := &Questions{
				QuestionnaireID: questionnaire.ID,
				Type:            "Checkbox",
			}
			err := db.
				Session(&gorm.Session{}).
				Create(&question).Error
			if err != nil {
				t.Errorf("failed to create question: %v", err)
			}

			if len(testCase.beforeOptions) > 0 {
				for i := range testCase.beforeOptions {
					testCase.beforeOptions[i].QuestionID = question.ID
				}
				err = db.
					Session(&gorm.Session{}).
					Create(&testCase.beforeOptions).Error
				if err != nil {
					t.Errorf("failed to create options: %v", err)
				}
			}

			err = optionImpl.UpdateOptions(ctx, testCase.argOptions, question.ID)

			if !testCase.isErr {
				assert.NoErrorf(t, err, testCase.description, "no error")
			} else if testCase.err != nil {
				if !errors.Is(err, testCase.err) {
					t.Errorf("invalid error(%s): expected: %+v, actual: %+v", testCase.description, testCase.err, err)
				}
			}
			if err != nil {
				return
			}

			var actualOptions []Options
			err = db.
				Session(&gorm.Session{}).
				Where("question_id = ?", question.ID).
				Select("QuestionID", "OptionNum", "Body").
				Order("option_num").
				Find(&actualOptions).Error
			if err != nil {
				t.Errorf("failed to get options: %v", err)
			}

			assert.Equalf(t, len(testCase.afterOptions), len(actualOptions), testCase.description, "length")
			for i, option := range testCase.afterOptions {
				assert.Equalf(t, question.ID, actualOptions[i].QuestionID, testCase.description, "questionID")
				assert.Equalf(t, option.OptionNum, actualOptions[i].OptionNum, testCase.description, "option num")
				assert.Equalf(t, option.Body, actualOptions[i].Body, testCase.description, "body")
			}
		})
	}
}
