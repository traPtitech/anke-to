// validations test
package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

type ValidationTest struct {
	validation Validations
	body       string
	result     string
}

func TestInsertValidation(t *testing.T) {
	t.Parallel()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), true), "public")
		require.NoError(t, err)
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)

		err = InsertValidation(-1, Validations{MinBound: "0", MaxBound: "10"})
		assert.Error(err)

		_ = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		err = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		assert.Error(err)
	})
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), true), "public")
		require.NoError(t, err)
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)

		err = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		assert.NoError(err)
	})
}

func TestUpdateValidation(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), true), "public")
		require.NoError(t, err)
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)
		err = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		require.NoError(t, err)

		err = UpdateValidation(questionID, Validations{MinBound: "10", MaxBound: "100"})
		assert.NoError(err)
	})
}

func TestDeleteValidation(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), true), "public")
		require.NoError(t, err)
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)
		err = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		require.NoError(t, err)

		err = DeleteValidation(questionID)
		assert.NoError(err)
	})
}

func TestGetValidation(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), true), "public")
		require.NoError(t, err)
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)
		err = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		require.NoError(t, err)

		validation, err := GetValidation(questionID)
		assert.NoError(err)
		assert.Equal(validation, Validations{QuestionID: questionID, MinBound: "0", MaxBound: "10"})
	})
}

func TestGetValidations(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), true), "public")
		require.NoError(t, err)
		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)
		err = InsertValidation(questionID, Validations{MinBound: "0", MaxBound: "10"})
		require.NoError(t, err)
		validations, err := GetValidations([]int{questionID})
		assert.NoError(err)
		assert.Equal(validations, []Validations{Validations{QuestionID: questionID, MinBound: "0", MaxBound: "10"}})
	})
}
func TestCheckNumberValidation(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{MinBound: "0", MaxBound: "10"},
			body:       "5",
			result:     "",
		},
		{
			validation: Validations{MinBound: "0", MaxBound: "10"},
			body:       "",
			result:     "",
		},
		{
			validation: Validations{MinBound: "0", MaxBound: "10"},
			body:       "a",
			result:     "strconv.Atoi: parsing \"a\": invalid syntax",
		},
		{
			validation: Validations{MinBound: "0", MaxBound: "10"},
			body:       "20",
			result:     "failed to meet the boundary value. the number must be less than MaxBound (number: 20, MaxBound: 10)",
		},
		{
			validation: Validations{MinBound: "10", MaxBound: ""},
			body:       "5",
			result:     "failed to meet the boundary value. the number must be greater than MinBound (number: 5, MinBound: 10)",
		},
		{
			validation: Validations{MinBound: "20", MaxBound: "10"},
			body:       "20",
			result:     "failed to check the boundary value. MinBound must be less than MaxBound (MinBound: 20, MaxBound: 10)",
		},
	}

	for _, validationTest := range validationTests {
		validation := validationTest.validation
		body := validationTest.body
		result := validationTest.result
		if err := CheckNumberValidation(validation, body); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}

func TestCheckTextValidation(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{RegexPattern: "^\\d*\\.\\d*$"},
			body:       "4.55",
			result:     "",
		},
		{
			validation: Validations{RegexPattern: "^\\d*\\.\\d*$"},
			body:       "a4.55",
			result:     "failed to match the pattern (Response: a4.55, RegexPattern: ^\\d*\\.\\d*$)",
		},
		{
			validation: Validations{RegexPattern: "*"},
			body:       "hoge",
			result:     "error parsing regexp: missing argument to repetition operator: `*`",
		},
	}
	for _, validationTest := range validationTests {
		validation := validationTest.validation
		body := validationTest.body
		result := validationTest.result
		if err := CheckTextValidation(validation, body); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}

func TestCheckNumberValid(t *testing.T) {
	validationTests := []ValidationTest{
		{
			validation: Validations{MinBound: "0", MaxBound: "10"},
			result:     "",
		},
		{
			validation: Validations{MinBound: "1a", MaxBound: "10"},
			result:     "failed to check the boundary value. MinBound is not a numerical value: strconv.Atoi: parsing \"1a\": invalid syntax",
		},
		{
			validation: Validations{MinBound: "10", MaxBound: "1a"},
			result:     "failed to check the boundary value. MaxBound is not a numerical value: strconv.Atoi: parsing \"1a\": invalid syntax",
		},
		{
			validation: Validations{MinBound: "10", MaxBound: "0"},
			result:     "failed to check the boundary value. MinBound must be less than MaxBound (MinBound: 10, MaxBound: 0)",
		},
	}

	for _, validationTest := range validationTests {
		validation := validationTest.validation
		result := validationTest.result
		if err := CheckNumberValid(validation.MinBound, validation.MaxBound); err != nil {
			assert.EqualError(t, err, result)
		} else {
			assert.NoError(t, nil)
		}
	}
}
