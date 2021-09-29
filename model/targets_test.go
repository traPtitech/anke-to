package model

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInsertTargets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	invalidQuestionnaireID := -1

	type test struct {
		description            string
		invalidQuestionnaireID bool
		beforeTargets          []string
		afterTargets           []string
		argTargets             []string
		isErr                  bool
		err                    error
	}

	testCases := []test{
		{
			description:   "元のtargetが1つで追加できる",
			beforeTargets: []string{"a"},
			afterTargets:  []string{"a", "b"},
			argTargets:    []string{"b"},
		},
		{
			description:   "元のtargetが複数でもエラーなし",
			beforeTargets: []string{"a", "b"},
			afterTargets:  []string{"a", "b", "c"},
			argTargets:    []string{"c"},
		},
		{
			description:   "追加するターゲットがなくてもエラーなし",
			beforeTargets: []string{"a"},
			afterTargets:  []string{"a"},
			argTargets:    []string{},
		},
		{
			description:   "元のtargetがなくてもエラーなし",
			beforeTargets: []string{},
			afterTargets:  []string{"a"},
			argTargets:    []string{"a"},
		},
		{
			description:            "questionnaireIDが誤っていたらエラー",
			invalidQuestionnaireID: true,
			argTargets:             []string{"b"},
			isErr:                  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var questionnaireID int
			if !testCase.invalidQuestionnaireID {
				targets := make([]Targets, 0, len(testCase.beforeTargets))
				for _, target := range testCase.beforeTargets {
					targets = append(targets, Targets{
						UserTraqid: target,
					})
				}
				questionnaire := Questionnaires{
					Targets: targets,
				}
				err := db.
					Session(&gorm.Session{}).
					Create(&questionnaire).Error
				if err != nil {
					t.Errorf("failed to create questionnaire: %v", err)
				}

				questionnaireID = questionnaire.ID
			} else {
				questionnaireID = invalidQuestionnaireID
			}

			err := targetImpl.InsertTargets(ctx, questionnaireID, testCase.argTargets)

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

			var targets []string
			err = db.
				Session(&gorm.Session{}).
				Model(&Targets{}).
				Where("questionnaire_id = ?", questionnaireID).
				Pluck("user_traqid", &targets).Error
			if err != nil {
				t.Errorf("failed to get targets: %v", err)
			}

			assert.ElementsMatchf(t, testCase.afterTargets, targets, "targets")
		})
	}
}
