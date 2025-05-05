package model

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
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

func TestDeleteTargets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type test struct {
		description   string
		beforeTargets []string
		isErr         bool
		err           error
	}

	testCases := []test{
		{
			description:   "targetが1人でもエラーなし",
			beforeTargets: []string{"a"},
		},
		{
			description:   "targetが複数でもエラーなし",
			beforeTargets: []string{"a", "b"},
		},
		{
			description:   "targetがなくてもエラーなし",
			beforeTargets: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
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

			err = targetImpl.DeleteTargets(ctx, questionnaire.ID)

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

			err = db.
				Session(&gorm.Session{}).
				Where("questionnaire_id = ?", questionnaire.ID).
				Take(&Targets{}).Error
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("invalid error(%s): expected: %+v, actual: %+v", testCase.description, gorm.ErrRecordNotFound, err)
			}
		})
	}
}

func TestGetTargets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type test struct {
		description            string
		questionnaires         []Questionnaires
		questionnaireIDIndexes []int
		isErr                  bool
		err                    error
	}

	testCases := []test{
		{
			description: "questionnaireが1つ、targetが1人でもエラーなし",
			questionnaires: []Questionnaires{
				{
					Targets: []Targets{
						{
							UserTraqid: "a",
						},
					},
				},
			},
			questionnaireIDIndexes: []int{0},
		},
		{
			description: "questionnaireが2つ、targetがそれぞれ1人でもエラーなし",
			questionnaires: []Questionnaires{
				{
					Targets: []Targets{
						{
							UserTraqid: "a",
						},
					},
				},
				{
					Targets: []Targets{
						{
							UserTraqid: "a",
						},
					},
				},
			},
			questionnaireIDIndexes: []int{0, 1},
		},
		{
			description: "一部のquestionnaireの取得でもエラーなし",
			questionnaires: []Questionnaires{
				{
					Targets: []Targets{
						{
							UserTraqid: "a",
						},
					},
				},
				{
					Targets: []Targets{
						{
							UserTraqid: "a",
						},
					},
				},
			},
			questionnaireIDIndexes: []int{0},
		},
		{
			description: "targetがなくてもエラーなし",
			questionnaires: []Questionnaires{
				{
					Targets: []Targets{},
				},
			},
			questionnaireIDIndexes: []int{0},
		},
		{
			description: "questionnaireIDがなくてもエラーなし",
			questionnaires: []Questionnaires{
				{
					Targets: []Targets{
						{
							UserTraqid: "a",
						},
					},
				},
			},
			questionnaireIDIndexes: []int{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			err := db.
				Session(&gorm.Session{}).
				Create(&testCase.questionnaires).Error
			if err != nil {
				t.Errorf("failed to create questionnaire: %v", err)
			}

			questionnaireIDs := make([]int, 0, len(testCase.questionnaireIDIndexes))
			for _, index := range testCase.questionnaireIDIndexes {
				questionnaireIDs = append(questionnaireIDs, testCase.questionnaires[index].ID)
			}

			targets, err := targetImpl.GetTargets(ctx, questionnaireIDs)

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

			expectTargets := make([]Targets, 0, len(testCase.questionnaireIDIndexes))
			for _, index := range testCase.questionnaireIDIndexes {
				expectTargets = append(expectTargets, testCase.questionnaires[index].Targets...)
			}

			assert.ElementsMatchf(t, expectTargets, targets, testCase.description, "targets")
		})
	}
}

func TestIsTargetingMe(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "private", true, false, true)
	require.NoError(t, err)

	err = targetImpl.InsertTargets(ctx, questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		userID string
	}
	type expect struct {
		isErr      bool
		err        error
		isTargeted bool
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "is targeted",
			args: args{
				userID: userOne,
			},
			expect: expect{
				isTargeted: true,
			},
		},
		{
			description: "not targeted",
			args: args{
				userID: userTwo,
			},
			expect: expect{
				isTargeted: false,
			},
		},
	}

	for _, testCase := range testCases {
		isTargeted, err := targetImpl.IsTargetingMe(ctx, questionnaireID, testCase.args.userID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}

		assertion.Equal(testCase.expect.isTargeted, isTargeted, testCase.description, "isTargeted")
	}
}

func TestGetTargetRemindStatus(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctx := context.Background()

	type test struct {
		description    string
		validTargets   []string // is_canceled = false
		invalidTargets []string // is_canceled = true
		argTarget      string
		argTargetValid bool
		isErr          bool
		err            error
	}

	testCases := []test{
		{
			description:    "all valid targets",
			validTargets:   []string{"a", "b"},
			invalidTargets: []string{},
			argTarget:      "a",
			argTargetValid: true,
		},
		{
			description:    "all invalid targets",
			validTargets:   []string{},
			invalidTargets: []string{"a", "b"},
			argTarget:      "b",
			argTargetValid: false,
		},
		{
			description:    "both valid and invalid targets",
			validTargets:   []string{"a", "b"},
			invalidTargets: []string{"c", "d"},
			argTarget:      "c",
			argTargetValid: false,
		},
		{
			description:    "both valid and invalid targets changed order",
			validTargets:   []string{"a", "b"},
			invalidTargets: []string{"c", "d"},
			argTarget:      "b",
			argTargetValid: true,
		},
		{
			description:    "both valid and invalid targets changed order duplicated arg targets",
			validTargets:   []string{"a", "b"},
			invalidTargets: []string{"c", "d"},
			argTarget:      "c",
			argTargetValid: false,
		},
		{
			description:    "argTargets with target not in db",
			validTargets:   []string{"a", "b"},
			invalidTargets: []string{"c", "d"},
			argTarget:      "e",
			isErr:          true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			targets := make([]Targets, 0, len(testCase.validTargets)+len(testCase.invalidTargets))
			for _, target := range testCase.validTargets {
				targets = append(targets, Targets{
					UserTraqid: target,
					IsCanceled: false,
				})
			}
			for _, target := range testCase.invalidTargets {
				targets = append(targets, Targets{
					UserTraqid: target,
					IsCanceled: true,
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

			remindStatus, err := targetImpl.GetTargetRemindStatus(ctx, questionnaire.ID, testCase.argTarget)
			if !testCase.isErr {
				assertion.NoError(err, testCase.description, "no error")
			} else if testCase.err != nil {
				assertion.Equal(true, errors.Is(err, testCase.err), testCase.description, "errorIs")
			} else {
				assertion.Error(err, testCase.description, "any error")
			}
			if err != nil {
				return
			}

			assert.Equal(t, remindStatus, testCase.argTargetValid, testCase.description, "remindStatus")
		})
	}
}

func TestUpdateTargetsRemindStatus(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctx := context.Background()

	type test struct {
		description          string
		beforeValidTargets   []string
		beforeInvalidTargets []string
		afterValidTargets    []string
		afterInvalidTargets  []string
		argUpdateTargets     []string
		argRemindStatus      bool
		isErr                bool
		err                  error
	}

	testCases := []test{
		{
			description:          "キャンセルするtargetが1人でエラーなし",
			beforeValidTargets:   []string{"a"},
			beforeInvalidTargets: []string{},
			afterValidTargets:    []string{},
			afterInvalidTargets:  []string{"a"},
			argUpdateTargets:     []string{"a"},
			argRemindStatus:      false,
		},
		{
			description:          "キャンセルするtargetが複数でエラーなし",
			beforeValidTargets:   []string{"a", "b"},
			beforeInvalidTargets: []string{},
			afterValidTargets:    []string{},
			afterInvalidTargets:  []string{"a", "b"},
			argUpdateTargets:     []string{"a", "b"},
			argRemindStatus:      false,
		},
		{
			description:          "キャンセルするtargetがないときエラーなし",
			beforeValidTargets:   []string{"a"},
			beforeInvalidTargets: []string{},
			afterValidTargets:    []string{"a"},
			afterInvalidTargets:  []string{},
			argUpdateTargets:     []string{},
			argRemindStatus:      false,
		},
		{
			description:          "キャンセルするtargetが見つからないときエラーなし",
			beforeValidTargets:   []string{"a"},
			beforeInvalidTargets: []string{},
			afterValidTargets:    []string{"a"},
			afterInvalidTargets:  []string{},
			argUpdateTargets:     []string{"b"},
			argRemindStatus:      false,
		},
		{
			description:          "再開するtargetが1人でエラーなし",
			beforeValidTargets:   []string{"a"},
			beforeInvalidTargets: []string{"b"},
			afterValidTargets:    []string{"a", "b"},
			afterInvalidTargets:  []string{},
			argUpdateTargets:     []string{"b"},
			argRemindStatus:      true,
		},
		{
			description:          "再開するtargetが複数でエラーなし",
			beforeValidTargets:   []string{},
			beforeInvalidTargets: []string{"a", "b"},
			afterValidTargets:    []string{"a", "b"},
			afterInvalidTargets:  []string{},
			argUpdateTargets:     []string{"a", "b"},
			argRemindStatus:      true,
		},
		{
			description:          "再開するtargetがないときエラーなし",
			beforeValidTargets:   []string{"a"},
			beforeInvalidTargets: []string{},
			afterValidTargets:    []string{"a"},
			afterInvalidTargets:  []string{},
			argUpdateTargets:     []string{},
			argRemindStatus:      true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			targets := make([]Targets, 0, len(testCase.beforeValidTargets)+len(testCase.beforeInvalidTargets))
			for _, target := range testCase.beforeValidTargets {
				targets = append(targets, Targets{
					UserTraqid: target,
					IsCanceled: false,
				})
			}
			for _, target := range testCase.beforeInvalidTargets {
				targets = append(targets, Targets{
					UserTraqid: target,
					IsCanceled: true,
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

			err = targetImpl.UpdateTargetsRemindStatus(ctx, questionnaire.ID, testCase.argUpdateTargets, testCase.argRemindStatus)
			if !testCase.isErr {
				assertion.NoError(err, testCase.description, "no error")
			} else if testCase.err != nil {
				assertion.Equal(true, errors.Is(err, testCase.err), testCase.description, "errorIs")
			} else {
				assertion.Error(err, testCase.description, "any error")
			}
			if err != nil {
				return
			}

			afterTargets := make([]Targets, 0, len(testCase.afterValidTargets)+len(testCase.afterInvalidTargets))
			for _, afterTarget := range testCase.afterInvalidTargets {
				afterTargets = append(afterTargets, Targets{
					QuestionnaireID: questionnaire.ID,
					UserTraqid:      afterTarget,
					IsCanceled:      true,
				})
			}
			for _, afterTarget := range testCase.afterValidTargets {
				afterTargets = append(afterTargets, Targets{
					QuestionnaireID: questionnaire.ID,
					UserTraqid:      afterTarget,
					IsCanceled:      false,
				})
			}

			actualTargets := make([]Targets, 0, len(testCase.afterValidTargets)+len(testCase.afterInvalidTargets))
			err = db.
				Session(&gorm.Session{}).
				Model(&Targets{}).
				Where("questionnaire_id = ?", questionnaire.ID).
				Find(&actualTargets).Error
			if err != nil {
				t.Errorf("failed to get targets: %v", err)
			}

			assert.ElementsMatchf(t, afterTargets, actualTargets, "targets")
		})
	}
}
