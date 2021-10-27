// validations test
package model

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func TestInsertValidation(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(ctx, questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validID      bool
		QuestionType string
		RegexPattern string
		MinBound     string
		MaxBound     string
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
			description: "valid Text(pattern not empty)",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: "^\\d*\\.\\d*$",
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "valid Text(pattern empty)",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: "",
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "long regexpattern",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: strings.Repeat("a", 2000),
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "too long regexpattern",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: strings.Repeat("a", 200000),
				MinBound:     "",
				MaxBound:     "",
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "valid Number(no bound empty)",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "0",
				MaxBound:     "10",
			},
		},
		{
			description: "valid Number(min bound empty)",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "",
				MaxBound:     "10",
			},
		},
		{
			description: "valid Number(max bound empty)",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "0",
				MaxBound:     "",
			},
		},
		{
			description: "valid Number(all bound empty)",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "long bounds",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     strings.Repeat("1", 2000),
				MaxBound:     strings.Repeat("1", 2000),
			},
		},
		{
			description: "too long bounds",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     strings.Repeat("1", 200000),
				MaxBound:     strings.Repeat("1", 200000),
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "question does not exist",
			args: args{
				validID:      false,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "",
				MaxBound:     "",
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		questionID, err := questionImpl.InsertQuestion(ctx, questionnaireID, 1, 1, testCase.QuestionType, testCase.QuestionType, true)
		require.NoError(t, err)
		if !testCase.args.validID {
			questionID = -1
		}

		validation := Validations{
			QuestionID:   questionID,
			RegexPattern: testCase.args.RegexPattern,
			MinBound:     testCase.args.MinBound,
			MaxBound:     testCase.args.MaxBound,
		}

		err = validationImpl.InsertValidation(ctx, questionID, validation)
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

		validation = Validations{}
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("question_id = ?", questionID).
			First(&validation).Error
		assertion.NoError(err, testCase.description, "get validations")

		assertion.Equal(questionID, validation.QuestionID, testCase.description, "questionID")
		assertion.Equal(testCase.args.RegexPattern, validation.RegexPattern, testCase.description, "RegexPattern")
		assertion.Equal(testCase.args.MinBound, validation.MinBound, testCase.description, "MinBound")
		assertion.Equal(testCase.args.MaxBound, validation.MaxBound, testCase.description, "MaxBound")
	}
}

func TestUpdateValidation(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(ctx, questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validID      bool
		QuestionType string
		RegexPattern string
		MinBound     string
		MaxBound     string
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
			description: "valid updated",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: "^(updated)+$",
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "Text no updated",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: "^(original)+$",
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "Number updated",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "-100",
				MaxBound:     "100",
			},
		},
		{
			description: "Number no updated",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "",
				MinBound:     "0",
				MaxBound:     "10",
			},
		},
		{
			description: "question does not exist",
			args: args{
				validID:      false,
				QuestionType: "Text",
				RegexPattern: "^(updated)+$",
				MinBound:     "",
				MaxBound:     "",
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordUpdated,
			},
		},
	}
	for _, testCase := range testCases {
		questionID, err := questionImpl.InsertQuestion(ctx, questionnaireID, 1, 1, testCase.args.QuestionType, testCase.args.QuestionType, true)
		require.NoError(t, err)

		validation := Validations{}

		if testCase.args.QuestionType == "Text" {
			validation = Validations{
				QuestionID:   questionID,
				RegexPattern: "^(original)+$",
				MinBound:     "",
				MaxBound:     "",
			}
		} else {
			validation = Validations{
				QuestionID:   questionID,
				RegexPattern: "",
				MinBound:     "0",
				MaxBound:     "10",
			}
		}

		err = validationImpl.InsertValidation(ctx, questionID, validation)
		require.NoError(t, err)

		if !testCase.args.validID {
			questionID = -1
		}

		validation = Validations{
			QuestionID:   questionID,
			RegexPattern: testCase.args.RegexPattern,
			MinBound:     testCase.args.MinBound,
			MaxBound:     testCase.args.MaxBound,
		}

		err = validationImpl.UpdateValidation(ctx, questionID, validation)

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

		validation = Validations{}
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("question_id = ?", questionID).
			First(&validation).Error
		assertion.NoError(err, testCase.description, "get validations")

		assertion.Equal(questionID, validation.QuestionID, testCase.description, "questionID")
		assertion.Equal(testCase.args.RegexPattern, validation.RegexPattern, testCase.description, "RegexPattern")
		assertion.Equal(testCase.args.MinBound, validation.MinBound, testCase.description, "MinBound")
		assertion.Equal(testCase.args.MaxBound, validation.MaxBound, testCase.description, "MaxBound")
	}
}

func TestDeleteValidation(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(ctx, questionnaireID, []string{userOne})
	require.NoError(t, err)

	type args struct {
		validID      bool
		QuestionType string
		RegexPattern string
		MinBound     string
		MaxBound     string
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
			description: "Text delete",
			args: args{
				validID:      true,
				QuestionType: "Text",
				RegexPattern: "^\\d*\\.\\d*$",
				MinBound:     "",
				MaxBound:     "",
			},
		},
		{
			description: "Number delete",
			args: args{
				validID:      true,
				QuestionType: "Number",
				RegexPattern: "^\\d*\\.\\d*$",
				MinBound:     "0",
				MaxBound:     "10",
			},
		},
		{
			description: "question does not exist",
			args: args{
				validID:      false,
				QuestionType: "Text",
				RegexPattern: "^\\d*\\.\\d*$",
				MinBound:     "",
				MaxBound:     "",
			},
			expect: expect{
				isErr: true,
				err:   ErrNoRecordDeleted,
			},
		},
	}
	for _, testCase := range testCases {
		questionID, err := questionImpl.InsertQuestion(ctx, questionnaireID, 1, 1, testCase.args.QuestionType, testCase.args.QuestionType, true)
		require.NoError(t, err)

		validation := Validations{
			QuestionID:   questionID,
			RegexPattern: testCase.args.RegexPattern,
			MinBound:     testCase.args.MinBound,
			MaxBound:     testCase.args.MaxBound,
		}

		err = validationImpl.InsertValidation(ctx, questionID, validation)
		require.NoError(t, err)

		if !testCase.args.validID {
			questionID = -1
		}

		err = validationImpl.DeleteValidation(ctx, questionID)

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

func TestGetValidations(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)
	ctx := context.Background()

	questionnaireID, err := questionnaireImpl.InsertQuestionnaire(ctx, "第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = administratorImpl.InsertAdministrators(ctx, questionnaireID, []string{userOne})
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
	validations := []Validations{
		{
			RegexPattern: "valid1",
			MinBound:     "",
			MaxBound:     "",
		},
		{
			RegexPattern: "valid2",
			MinBound:     "",
			MaxBound:     "",
		},
		{
			RegexPattern: "valid2",
			MinBound:     "",
			MaxBound:     "",
		},
	}

	questionIDs := make([]int, 0, 3)
	validationMap := make(map[int]Validations)
	for _, validation := range validations {
		questionID, err := questionImpl.InsertQuestion(ctx, questionnaireID, 1, 1, "Text", "Text", true)
		require.NoError(t, err)
		err = validationImpl.InsertValidation(ctx, questionID, validation)
		require.NoError(t, err)
		validation.QuestionID = questionID
		questionIDs = append(questionIDs, questionID)
		validationMap[questionID] = validation
	}

	testCases := []test{
		{
			description: "Text",
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

		validations, err := validationImpl.GetValidations(ctx, testCase.args.questionIDs)
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
		assertion.Equal(testCase.expect.length, len(validations), testCase.description, "length")
		if len(validations) < 1 {
			continue
		}

		for _, validation := range validations {
			expectValidation, ok := validationMap[validation.QuestionID]
			require.Equal(t, true, ok)

			assertion.Equal(expectValidation.QuestionID, validation.QuestionID, testCase.description, "questionID")
			assertion.Equal(expectValidation.RegexPattern, validation.RegexPattern, testCase.description, "RegexPattern")
			assertion.Equal(expectValidation.MinBound, validation.MinBound, testCase.description, "MinBound")
			assertion.Equal(expectValidation.MaxBound, validation.MaxBound, testCase.description, "MaxBound")
		}
	}
}

func TestCheckNumberValidation(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		validation Validations
		response   string
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
				validation: Validations{
					MinBound: "0",
					MaxBound: "10",
				},
				response: "5",
			},
		},
		{
			description: "no response",
			args: args{
				validation: Validations{
					MinBound: "0",
					MaxBound: "10",
				},
				response: "",
			},
		},
		{
			description: "invalid response",
			args: args{
				validation: Validations{
					MinBound: "0",
					MaxBound: "10",
				},
				response: "a",
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidNumber,
			},
		},
		{
			description: "maxBound over",
			args: args{
				validation: Validations{
					MinBound: "0",
					MaxBound: "10",
				},
				response: "20",
			},
			expect: expect{
				isErr: true,
				err:   ErrNumberBoundary,
			},
		},
		{
			description: "minBound over",
			args: args{
				validation: Validations{
					MinBound: "0",
					MaxBound: "10",
				},
				response: "-10",
			},
			expect: expect{
				isErr: true,
				err:   ErrNumberBoundary,
			},
		},
		{
			description: "invalid bounds",
			args: args{
				validation: Validations{
					MinBound: "20",
					MaxBound: "10",
				},
				response: "",
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidNumber,
			},
		},
	}
	for _, testCase := range testCases {

		err := validationImpl.CheckNumberValidation(testCase.args.validation, testCase.args.response)

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

func TestCheckTextValidation(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		validation Validations
		response   string
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
				validation: Validations{
					RegexPattern: "^\\d*\\.\\d*$",
				},
				response: "4.55",
			},
		},
		{
			description: "not match",
			args: args{
				validation: Validations{
					RegexPattern: "^\\d*\\.\\d*$",
				},
				response: "a4.55",
			},
			expect: expect{
				isErr: true,
				err:   ErrTextMatching,
			},
		},
		{
			description: "not match",
			args: args{
				validation: Validations{
					RegexPattern: "*",
				},
				response: "",
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidRegex,
			},
		},
		{
			description: "empty pattern",
			args: args{
				validation: Validations{
					RegexPattern: "",
				},
				response: "4.55",
			},
		},
	}
	for _, testCase := range testCases {

		err := validationImpl.CheckTextValidation(testCase.args.validation, testCase.args.response)

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

func TestCheckNumberValid(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		validation Validations
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
				validation: Validations{
					MinBound: "0",
					MaxBound: "10",
				},
			},
		},
		{
			description: "invalid min bound",
			args: args{
				validation: Validations{
					MinBound: "1a",
					MaxBound: "10",
				},
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidNumber,
			},
		},
		{
			description: "invalid max bound",
			args: args{
				validation: Validations{
					MinBound: "10",
					MaxBound: "1a",
				},
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidNumber,
			},
		},
		{
			description: "min exceeds max",
			args: args{
				validation: Validations{
					MinBound: "10",
					MaxBound: "0",
				},
			},
			expect: expect{
				isErr: true,
				err:   ErrInvalidNumber,
			},
		},
	}
	for _, testCase := range testCases {

		err := validationImpl.CheckNumberValid(testCase.args.validation.MinBound, testCase.args.validation.MaxBound)

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
