package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestInsertRespondent(t *testing.T) {
	t.Parallel()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, _ := insertTestRespondents(t)

		_, err := InsertRespondent("TooooooooooooooooooooLongUserName", questionnaireID, null.NewTime(time.Now(), false))
		assert.Error(err)
	})
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, _ := insertTestRespondents(t)

		_, err := InsertRespondent(userOne, questionnaireID, null.NewTime(time.Now(), true))
		assert.NoError(err)
	})
}

func TestUpdateSubmittedAt(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, responseID := insertTestRespondents(t)

		responseID, err := InsertRespondent(userOne, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)

		err = UpdateSubmittedAt(responseID)
		assert.NoError(err)
	})
}

func TestDeleteRespondent(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, responseID := insertTestRespondents(t)

		responseID, err := InsertRespondent(userOne, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)

		err = DeleteRespondent(userOne, responseID)
		assert.NoError(err)
	})
}

func TestGetRespondentInfos(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		questionnaireID, _, responseID := insertTestRespondents(t)

		respondentInfos, err := GetRespondentInfos(userOne, questionnaireID)
		assert.NoError(err)
		assert.Equal(1, len(respondentInfos))
		respondentInfo := respondentInfos[0]
		assert.Equal(responseID, respondentInfo.ResponseID)
		assert.Equal("", respondentInfo.ResTimeLimit)
		assert.Equal(questionnaireID, respondentInfo.QuestionnaireID)
	})
}

func TestGetRespondentDetail(t *testing.T) {
	t.Parallel()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		_, err := GetRespondentDetail(0)
		assert.Error(err)
	})
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, questionID, responseID := insertTestRespondents(t)

		respondentDetail, err := GetRespondentDetail(responseID)
		assert.NoError(err)
		assert.Equal(questionnaireID, respondentDetail.QuestionnaireID)

		assert.Equal(1, len(respondentDetail.Responses))
		responseBody := respondentDetail.Responses[0]
		assert.Equal(questionID, responseBody.QuestionID)
		assert.Equal("Text", responseBody.QuestionType)
		assert.Equal("リマインダーBOTを作った話", responseBody.Body.String)
	})
}

func TestGetRespondentDetails(t *testing.T) {
	t.Parallel()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		_, err := GetRespondentDetails(0, "invalid_sort")
		assert.Error(err)
	})
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, questionID, responseID := insertTestRespondents(t)

		respondentDetails, err := GetRespondentDetails(questionnaireID, "traqid")
		assert.NoError(err)
		assert.Equal(1, len(respondentDetails))
		respondentDetail := respondentDetails[0]
		assert.Equal(responseID, respondentDetail.ResponseID)

		assert.Equal(1, len(respondentDetail.Responses))
		responseBody := respondentDetail.Responses[0]
		assert.Equal(questionID, responseBody.QuestionID)
		assert.Equal("Text", responseBody.QuestionType)
		assert.Equal("リマインダーBOTを作った話", responseBody.Body.String)
	})
}

func TestGetRespondentsUserIDs(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, _ := insertTestRespondents(t)

		respondents, err := GetRespondentsUserIDs([]int{questionnaireID})
		assert.NoError(err)
		assert.Equal(1, len(respondents))
		respondent := respondents[0]
		assert.Equal(questionnaireID, respondent.QuestionnaireID)
		assert.Equal(userOne, respondent.UserTraqid)
	})
}

func TestCheckRespondent(t *testing.T) {
	t.Parallel()
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, _ := insertTestRespondents(t)

		isRespondent, err := CheckRespondent(userTwo, questionnaireID)
		assert.NoError(err)
		assert.Equal(false, isRespondent)
	})
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, _, _ := insertTestRespondents(t)

		isRespondent, err := CheckRespondent(userOne, questionnaireID)
		assert.NoError(err)
		assert.Equal(true, isRespondent)
	})
}

func insertTestRespondents(t *testing.T) (int, int, int) {
	questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
	require.NoError(t, err)

	responseID, err := InsertRespondent(userOne, questionnaireID, null.NewTime(time.Now(), true))
	require.NoError(t, err)

	err = InsertResponses(responseID, []*ResponseMeta{{QuestionID: questionID, Data: "リマインダーBOTを作った話"}})
	require.NoError(t, err)

	return questionnaireID, questionID, responseID
}
