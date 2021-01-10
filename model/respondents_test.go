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

		_, _, responseID := insertTestRespondents(t)

		err := UpdateSubmittedAt(responseID)
		assert.NoError(err)
	})
}

func TestDeleteRespondent(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		_, _, responseID := insertTestRespondents(t)

		err := DeleteRespondent(userOne, responseID)
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

		questionnaireID, questionIDs, responseID := insertTestRespondents(t)

		respondentDetail, err := GetRespondentDetail(responseID)
		assert.NoError(err)
		assert.Equal(questionnaireID, respondentDetail.QuestionnaireID)

		assert.Equal(2, len(respondentDetail.Responses))

		questionID := questionIDs[0]
		responseBody := respondentDetail.Responses[0]
		assert.Equal(questionID, responseBody.QuestionID)
		assert.Equal("Text", responseBody.QuestionType)
		assert.Equal("リマインダーBOTを作った話", responseBody.Body.String)

		questionID = questionIDs[1]
		responseBody = respondentDetail.Responses[1]
		assert.Equal(1, len(responseBody.OptionResponse))
		optionResponse := responseBody.OptionResponse[0]
		assert.Equal(questionID, responseBody.QuestionID)
		assert.Equal("MultipleChoice", responseBody.QuestionType)
		assert.Equal("選択肢1", optionResponse)

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

		questionnaireID, questionIDs, responseID := insertTestRespondents(t)

		respondentDetails, err := GetRespondentDetails(questionnaireID, "traqid")
		assert.NoError(err)
		assert.Equal(1, len(respondentDetails))
		respondentDetail := respondentDetails[0]
		assert.Equal(responseID, respondentDetail.ResponseID)

		assert.Equal(2, len(respondentDetail.Responses))

		questionID := questionIDs[0]
		responseBody := respondentDetail.Responses[0]
		assert.Equal(questionID, responseBody.QuestionID)
		assert.Equal("Text", responseBody.QuestionType)
		assert.Equal("リマインダーBOTを作った話", responseBody.Body.String)

		questionID = questionIDs[1]
		responseBody = respondentDetail.Responses[1]
		assert.Equal(1, len(responseBody.OptionResponse))
		optionResponse := responseBody.OptionResponse[0]
		assert.Equal(questionID, responseBody.QuestionID)
		assert.Equal("MultipleChoice", responseBody.QuestionType)
		assert.Equal("選択肢1", optionResponse)
	})
	t.Run("success_sort", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
		require.NoError(t, err)

		err = InsertAdministrators(questionnaireID, []string{userOne})
		require.NoError(t, err)

		questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
		require.NoError(t, err)

		responseID1, err := InsertRespondent(userOne, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)

		err = InsertResponses(responseID1, []*ResponseMeta{
			{QuestionID: questionID, Data: "リマインダーBOTを作った話1"},
		})
		time.Sleep(time.Millisecond * 1000)

		responseID2, err := InsertRespondent(userTwo, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)

		err = InsertResponses(responseID2, []*ResponseMeta{
			{QuestionID: questionID, Data: "リマインダーBOTを作った話2"},
		})
		time.Sleep(time.Millisecond * 1000)

		responseID3, err := InsertRespondent(userThree, questionnaireID, null.NewTime(time.Now(), true))
		require.NoError(t, err)

		err = InsertResponses(responseID3, []*ResponseMeta{
			{QuestionID: questionID, Data: "リマインダーBOTを作った話3"},
		})

		respondentDetails, err := GetRespondentDetails(questionnaireID, "traqid")
		assert.NoError(err)
		assert.Equal(3, len(respondentDetails))
		respondentDetail := respondentDetails[0]
		assert.Equal(responseID1, respondentDetail.ResponseID)
		assert.Equal(userOne, respondentDetail.TraqID)
		respondentDetail = respondentDetails[1]
		assert.Equal(responseID2, respondentDetail.ResponseID)
		assert.Equal(userTwo, respondentDetail.TraqID)
		respondentDetail = respondentDetails[2]
		assert.Equal(responseID3, respondentDetail.ResponseID)
		assert.Equal(userThree, respondentDetail.TraqID)

		respondentDetails, err = GetRespondentDetails(questionnaireID, "-traqid")
		assert.NoError(err)
		assert.Equal(3, len(respondentDetails))
		respondentDetail = respondentDetails[0]
		assert.Equal(responseID3, respondentDetail.ResponseID)
		assert.Equal(userThree, respondentDetail.TraqID)
		respondentDetail = respondentDetails[1]
		assert.Equal(responseID2, respondentDetail.ResponseID)
		assert.Equal(userTwo, respondentDetail.TraqID)
		respondentDetail = respondentDetails[2]
		assert.Equal(responseID1, respondentDetail.ResponseID)
		assert.Equal(userOne, respondentDetail.TraqID)

		respondentDetails, err = GetRespondentDetails(questionnaireID, "submitted_at")
		assert.NoError(err)
		assert.Equal(3, len(respondentDetails))
		respondentDetail = respondentDetails[0]
		assert.Equal(responseID1, respondentDetail.ResponseID)
		assert.Equal(userOne, respondentDetail.TraqID)
		respondentDetail = respondentDetails[1]
		assert.Equal(responseID2, respondentDetail.ResponseID)
		assert.Equal(userTwo, respondentDetail.TraqID)
		respondentDetail = respondentDetails[2]
		assert.Equal(responseID3, respondentDetail.ResponseID)
		assert.Equal(userThree, respondentDetail.TraqID)

		respondentDetails, err = GetRespondentDetails(questionnaireID, "-submitted_at")
		assert.NoError(err)
		assert.Equal(3, len(respondentDetails))
		respondentDetail = respondentDetails[0]
		assert.Equal(responseID3, respondentDetail.ResponseID)
		assert.Equal(userThree, respondentDetail.TraqID)
		respondentDetail = respondentDetails[1]
		assert.Equal(responseID2, respondentDetail.ResponseID)
		assert.Equal(userTwo, respondentDetail.TraqID)
		respondentDetail = respondentDetails[2]
		assert.Equal(responseID1, respondentDetail.ResponseID)
		assert.Equal(userOne, respondentDetail.TraqID)
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

func insertTestRespondents(t *testing.T) (int, []int, int) {
	questionnaireID, err := InsertQuestionnaire("第1回集会らん☆ぷろ募集アンケート", "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！", null.NewTime(time.Now(), false), "public")
	require.NoError(t, err)

	err = InsertAdministrators(questionnaireID, []string{userOne})
	require.NoError(t, err)

	questionIDs := make([]int, 0, 2)

	questionID, err := InsertQuestion(questionnaireID, 1, 1, "Text", "質問文", true)
	require.NoError(t, err)
	questionIDs = append(questionIDs, questionID)

	questionID, err = InsertQuestion(questionnaireID, 1, 3, "MultipleChoice", "radio", true)
	require.NoError(t, err)
	questionIDs = append(questionIDs, questionID)

	responseID, err := InsertRespondent(userOne, questionnaireID, null.NewTime(time.Now(), true))
	require.NoError(t, err)

	err = InsertResponses(responseID, []*ResponseMeta{
		{QuestionID: questionIDs[0], Data: "リマインダーBOTを作った話"},
		{QuestionID: questionIDs[1], Data: "選択肢1"},
	})
	require.NoError(t, err)
	return questionnaireID, questionIDs, responseID
}
