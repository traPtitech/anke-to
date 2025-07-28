// todo: set up the mock server for user group and add testCases for user group

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
	"gopkg.in/guregu/null.v4"
)

const (
	userOne   = "userOne"
	userTwo   = "userTwo"
	userThree = "userThree"
	userFour  = "userFour"
	userFive  = "userFive"
)

var (
	// groupOne   = uuid.MustParse("3d123d8d-509b-4221-bc1d-82e94deac563") // userOne, userTwo
	// groupTwo   = uuid.MustParse("c9e3766a-a307-4100-9df1-24943de026c2") // userThree, userFour
	// groupThree = uuid.MustParse("313e4211-715c-4ae6-a247-79a204c50382")
	// groupFour  = uuid.MustParse("8564a706-f9d4-46f5-852c-3d44a474902a")
	// groupFive  = uuid.MustParse("db7f3c13-eb4f-4890-9773-286dde024d4c")

	sampleAdmin                          = openapi.UsersAndGroups{}
	sampleTarget                         = openapi.UsersAndGroups{}
	sampleQuestionSettingsText           = openapi.NewQuestion{}
	sampleQuestionSettingsTextLong       = openapi.NewQuestion{}
	sampleQuestionSettingsNumber         = openapi.NewQuestion{}
	sampleQuestionSettingsSingleChoice   = openapi.NewQuestion{}
	sampleQuestionSettingsMultipleChoice = openapi.NewQuestion{}
	sampleQeustionsettingsScale          = openapi.NewQuestion{}
	sampleQuestionnaire                  = openapi.PostQuestionnaireJSONRequestBody{}
)

func setupSampleQuestionnaire() {
	if sampleQuestionnaire.Title != "" {
		return
	}
	sampleQuestionSettingsText = openapi.NewQuestion{
		Title:       "質問（テキスト）",
		IsRequired: true,
	}
	sampleQuestionSettingsTextMaxLength := 100
	sampleQuestionSettingsText.FromQuestionSettingsText(openapi.QuestionSettingsText{
		MaxLength:    &sampleQuestionSettingsTextMaxLength,
		QuestionType: openapi.QuestionSettingsTextQuestionTypeText,
	})
	sampleQuestionSettingsTextLong = openapi.NewQuestion{
		Title:       "質問（ロングテキスト）",
		IsRequired: true,
	}
	sampleQuestionSettingsTextLongMaxLength := 500
	sampleQuestionSettingsTextLong.FromQuestionSettingsTextLong(openapi.QuestionSettingsTextLong{
		MaxLength:    &sampleQuestionSettingsTextLongMaxLength,
		QuestionType: openapi.QuestionSettingsTextLongQuestionTypeTextLong,
	})
	sampleQuestionSettingsNumber = openapi.NewQuestion{
		Title:       "質問（数値）",
		IsRequired: true,
	}
	sampleQuestionSettingsNumberMaxValue := 100
	sampleQuestionSettingsNumberMinValue := 0
	sampleQuestionSettingsNumber.FromQuestionSettingsNumber(openapi.QuestionSettingsNumber{
		MaxValue:     &sampleQuestionSettingsNumberMaxValue,
		MinValue:     &sampleQuestionSettingsNumberMinValue,
		QuestionType: openapi.QuestionSettingsNumberQuestionTypeNumber,
	})
	sampleQuestionSettingsSingleChoice = openapi.NewQuestion{
		Title:       "質問（単一選択）",
		IsRequired: true,
	}
	sampleQuestionSettingsSingleChoice.FromQuestionSettingsSingleChoice(openapi.QuestionSettingsSingleChoice{
		Options:      []string{"選択肢A", "選択肢B", "選択肢C", "選択肢D"},
		QuestionType: openapi.QuestionSettingsSingleChoiceQuestionTypeSingleChoice,
	})
	sampleQuestionSettingsMultipleChoice = openapi.NewQuestion{
		Title:       "質問（複数選択）",
		IsRequired: true,
	}
	sampleQuestionSettingsMultipleChoice.FromQuestionSettingsMultipleChoice(openapi.QuestionSettingsMultipleChoice{
		Options:      []string{"選択肢A", "選択肢B", "選択肢C", "選択肢D"},
		QuestionType: openapi.QuestionSettingsMultipleChoiceQuestionTypeMultipleChoice,
	})
	sampleQeustionsettingsScale = openapi.NewQuestion{
		Title:       "質問（スケール）",
		IsRequired: true,
	}
	sampleQeustionsettingsScaleMaxLabel := "最大値"
	sampleQeustionsettingsScaleMinLabel := "最小値"
	sampleQeustionsettingsScale.FromQuestionSettingsScale(openapi.QuestionSettingsScale{
		MaxLabel:     &sampleQeustionsettingsScaleMaxLabel,
		MaxValue:     10,
		MinLabel:     &sampleQeustionsettingsScaleMinLabel,
		MinValue:     1,
		QuestionType: openapi.QuestionSettingsScaleQuestionTypeScale,
	})

	sampleAdmin = openapi.UsersAndGroups{
		Users:  []string{userOne},
		Groups: []uuid.UUID{},
	}

	sampleTarget = openapi.UsersAndGroups{
		Users:  []string{userThree},
		Groups: []uuid.UUID{},
	}

	sampleQuestionnaire = openapi.PostQuestionnaireJSONRequestBody{
		Admin:                    sampleAdmin,
		Description:              "第1回集会らん☆ぷろ参加者募集",
		IsDuplicateAnswerAllowed: true,
		IsAnonymous:              false,
		IsPublished:              true,
		Questions: []openapi.NewQuestion{
			sampleQuestionSettingsText,
			sampleQuestionSettingsTextLong,
			sampleQuestionSettingsNumber,
			sampleQuestionSettingsSingleChoice,
			sampleQuestionSettingsMultipleChoice,
			sampleQeustionsettingsScale,
		},
		ResponseDueDateTime: nil,
		ResponseViewableBy:  "anyone",
		Target:              sampleTarget,
		Title:               "第1回集会らん☆ぷろ募集アンケート",
	}
}

func newQuestion2Question(questionId *int, createdAt *time.Time, newQuestion openapi.NewQuestion) (openapi.Question, error) {
	b, err := newQuestion.MarshalJSON()
	if err != nil {
		return openapi.Question{}, err
	}
	var questionParsed map[string]interface{}
	err = json.Unmarshal([]byte(b), &questionParsed)
	if err != nil {
		return openapi.Question{}, err
	}
	if questionId != nil {
		questionParsed["question_id"] = questionId
	}
	if createdAt != nil {
		questionParsed["created_at"] = createdAt
	}
	b, err = json.Marshal(questionParsed)
	if err != nil {
		return openapi.Question{}, err
	}
	var question openapi.Question
	err = question.UnmarshalJSON(b)
	if err != nil {
		return openapi.Question{}, err
	}
	return question, nil
}

func postQuestionnaireParams2EditQuestionnaireParams(questionnaireId int, questions []openapi.Question, postQuestionnaireParams openapi.PostQuestionnaireJSONRequestBody) openapi.EditQuestionnaireJSONRequestBody {
	editQuestionnaireParams := openapi.EditQuestionnaireJSONRequestBody{
		Admin:                    &postQuestionnaireParams.Admin,
		Description:              postQuestionnaireParams.Description,
		IsDuplicateAnswerAllowed: postQuestionnaireParams.IsDuplicateAnswerAllowed,
		IsAnonymous:              postQuestionnaireParams.IsAnonymous,
		IsPublished:              postQuestionnaireParams.IsPublished,
		QuestionnaireId:          questionnaireId,
		Questions:                questions,
		ResponseDueDateTime:      postQuestionnaireParams.ResponseDueDateTime,
		ResponseViewableBy:       postQuestionnaireParams.ResponseViewableBy,
		Target:                   &postQuestionnaireParams.Target,
		Title:                    postQuestionnaireParams.Title,
	}
	return editQuestionnaireParams
}

func TestGetQuestionnaires(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaire := sampleQuestionnaire
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	_, err = q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.Title = "search test"
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnairePosted1, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.Title = "search test"
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnairePosted2, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.Title = "abcde"
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.Target.Users = []string{"specialTargetUser"}
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnairePosted4, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	type args struct {
		userID string
		params openapi.GetQuestionnairesParams
	}
	type expect struct {
		isErr               bool
		err                 error
		questionnaireIdList *[]int
	}
	type test struct {
		description string
		args
		expect
	}

	sortInvalid := (openapi.SortInQuery)("abcde")
	sortCreatedAt := (openapi.SortInQuery)("created_at")
	sortCreatedAtDesc := (openapi.SortInQuery)("-created_at")
	sortTitle := (openapi.SortInQuery)("title")
	sortTitleDesc := (openapi.SortInQuery)("-title")
	sortModifiedAt := (openapi.SortInQuery)("modified_at")
	sortModifiedAtDesc := (openapi.SortInQuery)("-modified_at")
	searchTest := (openapi.SearchInQuery)("search test")
	largePageNum := 100000000
	constTrue := true

	testCases := []test{
		{
			description: "valid",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{},
			},
		},
		{
			description: "invalid param sort",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortInvalid,
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "too large page num",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Page: &largePageNum,
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "sort created_at",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortCreatedAt,
				},
			},
		},
		{
			description: "sort -created_at",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortCreatedAtDesc,
				},
			},
		},
		{
			description: "sort title",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortTitle,
				},
			},
		},
		{
			description: "sort -title",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortTitleDesc,
				},
			},
		},
		{
			description: "sort modified_at",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortModifiedAt,
				},
			},
		},
		{
			description: "sort -modified_at",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort: &sortModifiedAtDesc,
				},
			},
		},
		{
			description: "search test",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Search: &searchTest,
				},
			},
			expect: expect{
				questionnaireIdList: &[]int{
					questionnairePosted1.QuestionnaireId,
					questionnairePosted2.QuestionnaireId,
				},
			},
		},
		{
			description: "search test sort created_at",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					Sort:   &sortCreatedAt,
					Search: &searchTest,
				},
			},
			expect: expect{
				questionnaireIdList: &[]int{
					questionnairePosted1.QuestionnaireId,
					questionnairePosted2.QuestionnaireId,
				},
			},
		},
		{
			description: "only targeting me user one",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					OnlyTargetingMe: &constTrue,
				},
			},
		},
		{
			description: "only targeting me user five",
			args: args{
				userID: userFive,
				params: openapi.GetQuestionnairesParams{
					OnlyTargetingMe: &constTrue,
				},
			},
		},
		{
			description: "only targeting by me special target user",
			args: args{
				userID: "specialTargetUser",
				params: openapi.GetQuestionnairesParams{
					OnlyTargetingMe: &constTrue,
				},
			},
			expect: expect{
				questionnaireIdList: &[]int{
					questionnairePosted4.QuestionnaireId,
				},
			},
		},
		{
			description: "only administrated by me user three",
			args: args{
				userID: userThree,
				params: openapi.GetQuestionnairesParams{
					OnlyAdministratedByMe: &constTrue,
				},
			},
		},
		{
			description: "only administrated by me userfive",
			args: args{
				userID: userFive,
				params: openapi.GetQuestionnairesParams{
					OnlyAdministratedByMe: &constTrue,
				},
			},
		},
	}

	for _, testCase := range testCases {
		params := url.Values{}
		if testCase.args.params.Sort != nil {
			params.Add("sort", string(*testCase.args.params.Sort))
		}
		if testCase.args.params.Search != nil {
			params.Add("search", string(*testCase.args.params.Search))
		}
		if testCase.args.params.Page != nil {
			params.Add("page", fmt.Sprint(*testCase.args.params.Page))
		}
		if testCase.args.params.OnlyTargetingMe != nil {
			params.Add("onlyTargetingMe", fmt.Sprint(*testCase.args.params.OnlyTargetingMe))
		}
		if testCase.args.params.OnlyAdministratedByMe != nil {
			params.Add("onlyAdministratedByMe", fmt.Sprint(*testCase.args.params.OnlyAdministratedByMe))
		}
		e = echo.New()
		req = httptest.NewRequest(http.MethodGet, "/questionnaires"+params.Encode(), nil)
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)

		questionnaireList, err := q.GetQuestionnaires(ctx, testCase.args.userID, testCase.args.params)

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

		if testCase.args.params.Sort != nil {
			if *testCase.args.params.Sort == "created_at" {
				var preCreatedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preCreatedAt.IsZero() {
						assertion.False(preCreatedAt.After(questionnaire.CreatedAt), testCase.description, "created_at")
					}
					preCreatedAt = questionnaire.CreatedAt
				}
			} else if *testCase.args.params.Sort == "-created_at" {
				var preCreatedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preCreatedAt.IsZero() {
						assertion.False(preCreatedAt.Before(questionnaire.CreatedAt), testCase.description, "-created_at")
					}
					preCreatedAt = questionnaire.CreatedAt
				}
			} else if *testCase.args.params.Sort == "title" {
				var preTitle string
				for _, questionnaire := range questionnaireList.Questionnaires {
					if preTitle != "" {
						assertion.False(preTitle > questionnaire.Title, testCase.description, "title")
					}
					preTitle = questionnaire.Title
				}
			} else if *testCase.args.params.Sort == "-title" {
				var preTitle string
				for _, questionnaire := range questionnaireList.Questionnaires {
					if preTitle != "" {
						assertion.False(preTitle < questionnaire.Title, testCase.description, "-title")
					}
					preTitle = questionnaire.Title
				}
			} else if *testCase.args.params.Sort == "modified_at" {
				var preModifiedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preModifiedAt.IsZero() {
						assertion.False(preModifiedAt.After(questionnaire.ModifiedAt), testCase.description, "modified_at")
					}
					preModifiedAt = questionnaire.ModifiedAt
				}
			} else if *testCase.args.params.Sort == "-modified_at" {
				var preModifiedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preModifiedAt.IsZero() {
						assertion.False(preModifiedAt.Before(questionnaire.ModifiedAt), testCase.description, "-modified_at")
					}
					preModifiedAt = questionnaire.ModifiedAt
				}
			}
		}

		if testCase.expect.questionnaireIdList != nil {
			var questionnaireIdList []int
			for _, questionnairSummary := range questionnaireList.Questionnaires {
				questionnaireIdList = append(questionnaireIdList, questionnairSummary.QuestionnaireId)
			}
			sort.Slice(*testCase.expect.questionnaireIdList, func(i, j int) bool {
				return (*testCase.expect.questionnaireIdList)[i] < (*testCase.expect.questionnaireIdList)[j]
			})
			sort.Slice(questionnaireIdList, func(i, j int) bool { return questionnaireIdList[i] < questionnaireIdList[j] })
			assertion.Equal(*testCase.expect.questionnaireIdList, questionnaireIdList, testCase.description, "questionnaireIdList")
		}

		if testCase.args.params.OnlyTargetingMe != nil || testCase.args.params.OnlyAdministratedByMe != nil {
			for _, questionnaire := range questionnaireList.Questionnaires {
				e = echo.New()
				require.NoError(t, err)
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/questionnaire/%d", questionnaire.QuestionnaireId), nil)
				rec = httptest.NewRecorder()
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				ctx = e.NewContext(req, rec)

				questionnaireDetail, err := q.GetQuestionnaire(ctx, questionnaire.QuestionnaireId)
				require.NoError(t, err)

				if testCase.args.params.OnlyTargetingMe != nil {
					assertion.Contains(questionnaireDetail.Targets, testCase.args.userID, testCase.description, "OnlyTargetingMe")
				}
				if testCase.args.params.OnlyAdministratedByMe != nil {
					assertion.Contains(questionnaireDetail.Admins, testCase.args.userID, testCase.description, "OnlyAdministratedByMe")
				}
			}
		}
	}
}

func TestPostQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		params openapi.PostQuestionnaireJSONRequestBody
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

	responseDueDateTimeMinus := time.Now().Add(-24 * time.Hour)
	responseDueDateTimePlus := time.Now().Add(24 * time.Hour)

	invalidQuestionSettingsNumber := openapi.NewQuestion{
		Title:       "質問（数値）",
		IsRequired: true,
	}
	invalidQuestionSettingsNumberMaxValue := 0
	invalidQuestionSettingsNumberMinValue := 100
	invalidQuestionSettingsNumber.FromQuestionSettingsNumber(openapi.QuestionSettingsNumber{
		MaxValue:     &invalidQuestionSettingsNumberMaxValue,
		MinValue:     &invalidQuestionSettingsNumberMinValue,
		QuestionType: openapi.QuestionSettingsNumberQuestionTypeNumber,
	})
	invalidQeustionsettingsScale := openapi.NewQuestion{
		Title:       "質問（スケール）",
		IsRequired: true,
	}
	invalidQeustionsettingsScaleMaxLabel := "最大値"
	invalidQeustionsettingsScaleMinLabel := "最小値"
	invalidQeustionsettingsScale.FromQuestionSettingsScale(openapi.QuestionSettingsScale{
		MaxLabel:     &invalidQeustionsettingsScaleMaxLabel,
		MaxValue:     1,
		MinLabel:     &invalidQeustionsettingsScaleMinLabel,
		MinValue:     10,
		QuestionType: openapi.QuestionSettingsScaleQuestionTypeScale,
	})

	testCases := []test{
		{
			description: "valid",
			args: args{
				params: sampleQuestionnaire,
			},
		},
		{
			description: "valid response due date time",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: &responseDueDateTimePlus,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "invalid response due date time",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: &responseDueDateTimeMinus,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "no title",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "",
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "too long title",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "response viewable by admins",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "admins",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "response viewable by respondents",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "respondents",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "response viewable by invalid",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "invalid",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "no admin",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    openapi.UsersAndGroups{},
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "invalid",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "duplicate answer not allowed",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: false,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "is anonymous",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              true,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "not published",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              false,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "invalid question settings number",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						invalidQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid question settings scale",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						invalidQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		e := echo.New()
		body, err := json.Marshal(testCase.args.params)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		questionnaireDetail, err := q.PostQuestionnaire(ctx, testCase.args.params)

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

		sort.Slice(questionnaireDetail.Admin.Users, func(i, j int) bool { return questionnaireDetail.Admin.Users[i] < questionnaireDetail.Admin.Users[j] })
		sort.Slice(testCase.args.params.Admin.Users, func(i, j int) bool { return testCase.args.params.Admin.Users[i] < testCase.args.params.Admin.Users[j] })
		assertion.Equal(questionnaireDetail.Admin.Users, testCase.args.params.Admin.Users, testCase.description, "admin users not equal")
		sort.Slice(questionnaireDetail.Admin.Groups, func(i, j int) bool {
			return questionnaireDetail.Admin.Groups[i].String() < questionnaireDetail.Admin.Groups[j].String()
		})
		sort.Slice(testCase.args.params.Admin.Groups, func(i, j int) bool {
			return testCase.args.params.Admin.Groups[i].String() < testCase.args.params.Admin.Groups[j].String()
		})
		assertion.Equal(questionnaireDetail.Admin.Groups, testCase.args.params.Admin.Groups, testCase.description, "admin groups not equal")

		assertion.Equal(questionnaireDetail.Description, testCase.args.params.Description, "description not equal")
		assertion.Equal(questionnaireDetail.IsDuplicateAnswerAllowed, testCase.args.params.IsDuplicateAnswerAllowed, "is duplicate answer allowed not equal")
		assertion.Equal(questionnaireDetail.IsAnonymous, testCase.args.params.IsAnonymous, "is anonymous not equal")
		assertion.Equal(questionnaireDetail.IsPublished, testCase.args.params.IsPublished, "is published not equal")

		for _, question := range testCase.args.params.Questions {
			isMatch := false
			for _, questionDetail := range questionnaireDetail.Questions {
				b, err := question.MarshalJSON()
				require.NoError(t, err)
				var questionParsed map[string]interface{}
				err = json.Unmarshal([]byte(b), &questionParsed)
				require.NoError(t, err)

				b, err = questionDetail.MarshalJSON()
				require.NoError(t, err)
				var questionDetailParsed map[string]interface{}
				err = json.Unmarshal([]byte(b), &questionDetailParsed)
				require.NoError(t, err)

				if questionParsed["body"] == questionDetailParsed["body"] &&
					questionParsed["is_required"] == questionDetailParsed["is_required"] &&
					questionParsed["question_type"] == questionDetailParsed["question_type"] {
					isMatch = true
					break
				}
			}
			if !isMatch {
				assertion.Fail("question not found", testCase.description)
			}
		}

		if testCase.args.params.ResponseDueDateTime != nil {
			assertion.WithinDuration(testCase.args.params.ResponseDueDateTime.UTC().Truncate(time.Second), questionnaireDetail.ResponseDueDateTime.UTC(), time.Second, testCase.description, "response due date time not equal")
		} else {
			assertion.Nil(questionnaireDetail.ResponseDueDateTime, testCase.description, "response due date time not equal")
		}
		if testCase.args.params.ResponseViewableBy != "invalid" {
			assertion.Equal(testCase.args.params.ResponseViewableBy, questionnaireDetail.ResponseViewableBy, "response viewable by not equal")
		} else {
			assertion.Equal((openapi.ResShareType)("admins"), questionnaireDetail.ResponseViewableBy, "response viewable by not equal")
		}
		sort.Slice(questionnaireDetail.Target.Users, func(i, j int) bool { return questionnaireDetail.Target.Users[i] < questionnaireDetail.Target.Users[j] })
		sort.Slice(testCase.args.params.Target.Users, func(i, j int) bool {
			return testCase.args.params.Target.Users[i] < testCase.args.params.Target.Users[j]
		})
		assertion.Equal(testCase.args.params.Target.Users, questionnaireDetail.Target.Users, "target users not equal")
		sort.Slice(questionnaireDetail.Target.Groups, func(i, j int) bool {
			return questionnaireDetail.Target.Groups[i].String() < questionnaireDetail.Target.Groups[j].String()
		})
		sort.Slice(testCase.args.params.Target.Groups, func(i, j int) bool {
			return testCase.args.params.Target.Groups[i].String() < testCase.args.params.Target.Groups[j].String()
		})
		assertion.Equal(testCase.args.params.Target.Groups, questionnaireDetail.Target.Groups, "target groups not equal")

		assertion.Equal(testCase.args.params.Title, questionnaireDetail.Title, "title not equal")
	}
}

func TestGetQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaire := sampleQuestionnaire
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	questionnaireDetailOrigin, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	type args struct {
		questionnaireID        int
		invalidQuestionnaireID bool
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
				questionnaireID: questionnaireDetailOrigin.QuestionnaireId,
			},
		},
		{
			description: "invalid questionnaire id",
			args: args{
				invalidQuestionnaireID: true,
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		var questionnaireID int
		if testCase.args.invalidQuestionnaireID {
			questionnaireID = 10000
			valid := true
			for valid {
				ctx := context.Background()
				_, _, _, _, _, _, _, _, err := IQuestionnaire.GetQuestionnaireInfo(ctx, questionnaireID)
				if errors.Is(err, model.ErrRecordNotFound) {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting questionnaire info", err)
				} else {
					questionnaireID *= 10
				}
			}
		} else {
			questionnaireID = testCase.args.questionnaireID
		}
		e := echo.New()
		body, err := json.Marshal(questionnaire)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/questionnaire/%d", questionnaireID), bytes.NewReader(body))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		questionnaireDetail, err := q.GetQuestionnaire(ctx, questionnaireID)

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

		assertion.Equal(questionnaireDetailOrigin, questionnaireDetail, testCase.description, "questionnaireDetail")
	}
}

func TestEditQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		isAnonymousToNotAnonymous bool
		invalidQuestionnaireID    bool
		invalidQuestionID         bool
		params                    openapi.PostQuestionnaireJSONRequestBody
		isNewQuestion             []bool
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

	responseDueDateTimeMinus := time.Now().Add(-24 * time.Hour)
	responseDueDateTimePlus := time.Now().Add(24 * time.Hour)

	invalidQuestionSettingsNumber := openapi.NewQuestion{
		Title:       "質問（数値）",
		IsRequired: true,
	}
	invalidQuestionSettingsNumberMaxValue := 0
	invalidQuestionSettingsNumberMinValue := 100
	invalidQuestionSettingsNumber.FromQuestionSettingsNumber(openapi.QuestionSettingsNumber{
		MaxValue:     &invalidQuestionSettingsNumberMaxValue,
		MinValue:     &invalidQuestionSettingsNumberMinValue,
		QuestionType: openapi.QuestionSettingsNumberQuestionTypeNumber,
	})
	invalidQeustionsettingsScale := openapi.NewQuestion{
		Title:       "質問（スケール）",
		IsRequired: true,
	}
	invalidQeustionsettingsScaleMaxLabel := "最大値"
	invalidQeustionsettingsScaleMinLabel := "最小値"
	invalidQeustionsettingsScale.FromQuestionSettingsScale(openapi.QuestionSettingsScale{
		MaxLabel:     &invalidQeustionsettingsScaleMaxLabel,
		MaxValue:     1,
		MinLabel:     &invalidQeustionsettingsScaleMinLabel,
		MinValue:     10,
		QuestionType: openapi.QuestionSettingsScaleQuestionTypeScale,
	})

	testCases := []test{
		{
			description: "valid",
			args: args{
				params:        sampleQuestionnaire,
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "valid new question",
			args: args{
				params:        sampleQuestionnaire,
				isNewQuestion: []bool{true, true, true, true, true, true},
			},
		},
		{
			description: "valid some new question",
			args: args{
				params:        sampleQuestionnaire,
				isNewQuestion: []bool{true, false, true, false, true, false},
			},
		},
		{
			description: "valid no question",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions:                []openapi.NewQuestion{},
					ResponseDueDateTime:      nil,
					ResponseViewableBy:       "anyone",
					Target:                   sampleTarget,
					Title:                    "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
		{
			description: "valid response due date time",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: &responseDueDateTimePlus,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "invalid question id",
			args: args{
				invalidQuestionID: true,
				params:            sampleQuestionnaire,
				isNewQuestion:     []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid anonymous to not anonymous",
			args: args{
				isAnonymousToNotAnonymous: true,
				params:                    sampleQuestionnaire,
				isNewQuestion:             []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response due date time",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: &responseDueDateTimeMinus,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "no title",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "too long title",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "response viewable by admins",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "admins",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "response viewable by respondents",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "respondents",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "response viewable by invalid",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "invalid",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "no admin",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    openapi.UsersAndGroups{},
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "invalid",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "duplicate answer not allowed",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: false,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "is anonymous",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              true,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "not published",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              false,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
		},
		{
			description: "invalid question settings number",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						invalidQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						sampleQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid question settings scale",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions: []openapi.NewQuestion{
						sampleQuestionSettingsText,
						sampleQuestionSettingsTextLong,
						sampleQuestionSettingsNumber,
						sampleQuestionSettingsSingleChoice,
						sampleQuestionSettingsMultipleChoice,
						invalidQeustionsettingsScale,
					},
					ResponseDueDateTime: nil,
					ResponseViewableBy:  "anyone",
					Target:              sampleTarget,
					Title:               "第1回集会らん☆ぷろ募集アンケート",
				},
				isNewQuestion: []bool{false, false, false, false, false, false},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "delete questions",
			args: args{
				params: openapi.PostQuestionnaireJSONRequestBody{
					Admin:                    sampleAdmin,
					Description:              "第1回集会らん☆ぷろ参加者募集",
					IsDuplicateAnswerAllowed: true,
					IsAnonymous:              false,
					IsPublished:              true,
					Questions:                []openapi.NewQuestion{},
					ResponseDueDateTime:      nil,
					ResponseViewableBy:       "anyone",
					Target:                   sampleTarget,
					Title:                    "第1回集会らん☆ぷろ募集アンケート",
				},
			},
		},
	}

	for _, testCase := range testCases {
		questionnaire := sampleQuestionnaire
		if testCase.args.isAnonymousToNotAnonymous {
			questionnaire.IsAnonymous = true
		}
		e := echo.New()
		body, err := json.Marshal(questionnaire)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
		require.NoError(t, err)

		var oldQuestionIDs, deletedQuestionIDs []int
		for i, isNewQuestion := range testCase.args.isNewQuestion {
			if !isNewQuestion {
				oldQuestionIDs = append(oldQuestionIDs, *questionnaireDetail.Questions[i].QuestionId)
			} else {
				deletedQuestionIDs = append(deletedQuestionIDs, *questionnaireDetail.Questions[i].QuestionId)
			}
		}

		var questionnaireID int
		if testCase.args.invalidQuestionnaireID {
			questionnaireID = 10000
			valid := true
			for valid {
				ctx := context.Background()
				_, _, _, _, _, _, _, _, err := IQuestionnaire.GetQuestionnaireInfo(ctx, questionnaireID)
				if errors.Is(err, model.ErrRecordNotFound) {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting questionnaire info")
				} else {
					questionnaireID *= 10
				}
			}
		} else {
			questionnaireID = questionnaireDetail.QuestionnaireId
		}

		var questions []openapi.Question
		for i, newQuestion := range testCase.args.params.Questions {
			if testCase.args.isNewQuestion[i] {
				question, err := newQuestion2Question(nil, nil, newQuestion)
				require.NoError(t, err)
				questions = append(questions, question)
			} else {
				questionID := *questionnaireDetail.Questions[i].QuestionId
				if testCase.args.invalidQuestionID {
					questionID = 10000
					valid := true
					for valid {
						ctx := context.Background()
						valid, err = IQuestion.CheckQuestionNum(ctx, questionnaireID, questionID)
						require.NoError(t, err)
						if valid {
							questionID *= 10
						}
					}
				}
				question, err := newQuestion2Question(&questionID, questionnaireDetail.Questions[i].CreatedAt, newQuestion)
				require.NoError(t, err)
				questions = append(questions, question)
			}
		}
		params := postQuestionnaireParams2EditQuestionnaireParams(questionnaireDetail.QuestionnaireId, questions, testCase.args.params)

		e = echo.New()
		body, err = json.Marshal(params)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/questionnaire/%d", questionnaireID), bytes.NewReader(body))
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		err = q.EditQuestionnaire(ctx, questionnaireID, params)

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

		questionnaireDetailEdited, err := q.GetQuestionnaire(ctx, questionnaireID)
		require.NoError(t, err)

		assertion.Equal(questionnaireDetail.QuestionnaireId, questionnaireDetailEdited.QuestionnaireId, testCase.description, "questionnaireId")

		for _, oldQuestionID := range oldQuestionIDs {
			exist := false
			for _, question := range questionnaireDetailEdited.Questions {
				if *question.QuestionId == oldQuestionID {
					exist = true
					break
				}
			}
			assertion.True(exist, testCase.description, "question was incorrectly deleted")
		}
		for _, dedeletedQuestionID := range deletedQuestionIDs {
			for _, question := range questionnaireDetailEdited.Questions {
				if *question.QuestionId == dedeletedQuestionID {
					assertion.Fail("question was not correctly deleted")
				}
			}
		}

		e = echo.New()
		body, err = json.Marshal(questionnaire)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		questionnaireDetailExpected, err := q.PostQuestionnaire(ctx, testCase.args.params)
		require.NoError(t, err)

		questionnaireDetailExpected.QuestionnaireId = questionnaireDetailEdited.QuestionnaireId
		questionnaireDetailExpected.CreatedAt = questionnaireDetailEdited.CreatedAt
		questionnaireDetailExpected.ModifiedAt = questionnaireDetailEdited.ModifiedAt

		assertion.Equal(len(questionnaireDetailExpected.Questions), len(questionnaireDetailEdited.Questions), testCase.description, "question length")
		for i := range questionnaireDetailExpected.Questions {
			questionnaireDetailExpected.Questions[i].QuestionId = questionnaireDetailEdited.Questions[i].QuestionId
			questionnaireDetailExpected.Questions[i].CreatedAt = questionnaireDetailEdited.Questions[i].CreatedAt
		}
		assertion.Equal(questionnaireDetailExpected, questionnaireDetailEdited, testCase.description, "questionnaireDetail")
	}
}

func TestDeleteQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		invalidQuestionnaireID bool
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
				invalidQuestionnaireID: false,
			},
		},
		{
			description: "invalid",
			args: args{
				invalidQuestionnaireID: true,
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		var questionnaireID int
		if !testCase.args.invalidQuestionnaireID {
			questionnaire := sampleQuestionnaire
			e := echo.New()
			body, err := json.Marshal(questionnaire)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
			require.NoError(t, err)
			questionnaireID = questionnaireDetail.QuestionnaireId
		} else {
			questionnaireID = 10000
			valid := true
			for valid {
				c := context.Background()
				_, _, _, _, _, _, _, _, err := IQuestionnaire.GetQuestionnaireInfo(c, questionnaireID)
				if errors.Is(err, model.ErrRecordNotFound) {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting questionnaire info")
				} else {
					questionnaireID *= 10
				}
			}
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/questionnaires/%d", questionnaireID), nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		err := q.DeleteQuestionnaire(ctx, questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}

		c := context.Background()
		_, _, _, _, _, _, _, _, err = IQuestionnaire.GetQuestionnaireInfo(c, questionnaireID)

		if err == nil {
			assertion.Fail("questionnaire not deleted")
		} else if !errors.Is(err, model.ErrRecordNotFound) {
			assertion.Fail("unexpected error during getting questionnaire info")
		}
	}
}

func TestGetQuestionnaireMyRemindStatus(t *testing.T) {
	// todo
}

func TestEditQuestionnaireMyRemindStatus(t *testing.T) {
	// todo
}

func TestGetQuestionnaireResponses(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	questionnaire := sampleQuestionnaire
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	AddQuestionID2SampleResponseMutex.Lock()

	AddQuestionID2SampleResponse(questionnaireDetail.QuestionnaireId)

	newResponse := sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response00, err := q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response01, err := q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	newResponse.IsDraft = true
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response03, err := q.PostQuestionnaireResponse(ctx, questionnaireDetail.QuestionnaireId, newResponse, userTwo)
	require.NoError(t, err)

	questionnaireAnonymous := sampleQuestionnaire
	questionnaireAnonymous.IsAnonymous = true
	e = echo.New()
	body, err = json.Marshal(questionnaireAnonymous)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireAnonymousDetail, err := q.PostQuestionnaire(ctx, questionnaireAnonymous)
	require.NoError(t, err)

	AddQuestionID2SampleResponse(questionnaireAnonymousDetail.QuestionnaireId)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireAnonymousDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response10, err := q.PostQuestionnaireResponse(ctx, questionnaireAnonymousDetail.QuestionnaireId, newResponse, userOne)
	require.NoError(t, err)

	newResponse = sampleResponse
	e = echo.New()
	body, err = json.Marshal(newResponse)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireAnonymousDetail.QuestionnaireId), bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	response11, err := q.PostQuestionnaireResponse(ctx, questionnaireAnonymousDetail.QuestionnaireId, newResponse, userTwo)
	require.NoError(t, err)

	AddQuestionID2SampleResponseMutex.Unlock()

	type args struct {
		isAnonymousQuestionnaire bool
		userID                   string
		questionnaireID          int
		params                   openapi.GetQuestionnaireResponsesParams
	}
	type expect struct {
		isErr          bool
		err            error
		responseIdList *[]int
	}
	type test struct {
		description string
		args
		expect
	}

	sortInvalid := (openapi.ResponseSortInQuery)("abcde")
	sortTraqID := (openapi.ResponseSortInQuery)("traqid")
	sortTraqIDDesc := (openapi.ResponseSortInQuery)("-traqid")
	sortSubmittedAt := (openapi.ResponseSortInQuery)("submitted_at")
	sortSubmittedAtDesc := (openapi.ResponseSortInQuery)("-submitted_at")
	sortModifiedAt := (openapi.ResponseSortInQuery)("modified_at")
	sortModifiedAtDesc := (openapi.ResponseSortInQuery)("-modified_at")
	constTrue := true

	testCases := []test{
		{
			description: "valid",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params:          openapi.GetQuestionnaireResponsesParams{},
			},
			expect: expect{
				responseIdList: &[]int{
					response00.ResponseId,
					response01.ResponseId,
					response03.ResponseId,
				},
			},
		},
		{
			description: "invalid param sort",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortInvalid,
				},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "sort submitted_at",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortSubmittedAt,
				},
			},
		},
		{
			description: "sort -submitted_at",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortSubmittedAtDesc,
				},
			},
		},
		{
			description: "sort traqid",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortTraqID,
				},
			},
		},
		{
			description: "sort -traqid",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortTraqIDDesc,
				},
			},
		},
		{
			description: "sort modified_at",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortModifiedAt,
				},
			},
		},
		{
			description: "sort -modified_at",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					Sort: &sortModifiedAtDesc,
				},
			},
		},
		{
			description: "only my response",
			args: args{
				userID:          userOne,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					OnlyMyResponse: &constTrue,
				},
			},
			expect: expect{
				responseIdList: &[]int{
					response00.ResponseId,
					response01.ResponseId,
				},
			},
		},
		{
			description: "only my response no response",
			args: args{
				userID:          userThree,
				questionnaireID: questionnaireDetail.QuestionnaireId,
				params: openapi.GetQuestionnaireResponsesParams{
					OnlyMyResponse: &constTrue,
				},
			},
			expect: expect{
				responseIdList: &[]int{},
			},
		},
		{
			description: "anonymous questionnaire",
			args: args{
				isAnonymousQuestionnaire: true,
				userID:                   userOne,
				questionnaireID:          questionnaireAnonymousDetail.QuestionnaireId,
				params:                   openapi.GetQuestionnaireResponsesParams{},
			},
			expect: expect{
				responseIdList: &[]int{
					response10.ResponseId,
					response11.ResponseId,
				},
			},
		},
	}

	for _, testCase := range testCases {
		params := url.Values{}
		if testCase.args.params.Sort != nil {
			params.Add("sort", string(*testCase.args.params.Sort))
		}
		if testCase.args.params.OnlyMyResponse != nil {
			params.Add("onlyMyResponse", fmt.Sprint(*testCase.args.params.OnlyMyResponse))
		}
		e = echo.New()
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/questionnaires/%d/responses", testCase.args.questionnaireID)+params.Encode(), nil)
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)

		responseList, err := q.GetQuestionnaireResponses(ctx, testCase.args.questionnaireID, testCase.args.params, testCase.args.userID)

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

		if testCase.args.params.Sort != nil {
			if *testCase.args.params.Sort == "submitted_at" {
				var preCreatedAt time.Time
				for _, response := range responseList {
					if !preCreatedAt.IsZero() {
						assertion.False(preCreatedAt.After(response.SubmittedAt), testCase.description, "submitted_at")
					}
					preCreatedAt = response.SubmittedAt
				}
			} else if *testCase.args.params.Sort == "-submitted_at" {
				var preCreatedAt time.Time
				for _, response := range responseList {
					if !preCreatedAt.IsZero() {
						assertion.False(preCreatedAt.Before(response.SubmittedAt), testCase.description, "-submitted_at")
					}
					preCreatedAt = response.SubmittedAt
				}
			} else if *testCase.args.params.Sort == "traqid" {
				var preTraqID string
				for _, response := range responseList {
					if preTraqID != "" {
						assertion.False(preTraqID > *response.Respondent, testCase.description, "traqid")
					}
					preTraqID = *response.Respondent
				}
			} else if *testCase.args.params.Sort == "-traqid" {
				var preTraqID string
				for _, response := range responseList {
					if preTraqID != "" {
						assertion.False(preTraqID < *response.Respondent, testCase.description, "-traqid")
					}
					preTraqID = *response.Respondent
				}
			} else if *testCase.args.params.Sort == "modified_at" {
				var preModifiedAt time.Time
				for _, response := range responseList {
					if !preModifiedAt.IsZero() {
						assertion.False(preModifiedAt.After(response.ModifiedAt), testCase.description, "modified_at")
					}
					preModifiedAt = response.ModifiedAt
				}
			} else if *testCase.args.params.Sort == "-modified_at" {
				var preModifiedAt time.Time
				for _, response := range responseList {
					if !preModifiedAt.IsZero() {
						assertion.False(preModifiedAt.Before(response.ModifiedAt), testCase.description, "-modified_at")
					}
					preModifiedAt = response.ModifiedAt
				}
			}
		}

		if testCase.expect.responseIdList != nil {
			responseIdList := []int{}
			for _, response := range responseList {
				responseIdList = append(responseIdList, response.ResponseId)
			}
			sort.Slice(*testCase.expect.responseIdList, func(i, j int) bool {
				return (*testCase.expect.responseIdList)[i] < (*testCase.expect.responseIdList)[j]
			})
			sort.Slice(responseIdList, func(i, j int) bool { return responseIdList[i] < responseIdList[j] })
			assertion.Equal(*testCase.expect.responseIdList, responseIdList, testCase.description, "responseIdList")
		}

		if testCase.args.params.OnlyMyResponse != nil {
			for _, response := range responseList {
				assertion.Equal(*response.Respondent, testCase.args.userID, testCase.description, "OnlyMyResponse")
			}
		}

		if testCase.args.isAnonymousQuestionnaire {
			for _, response := range responseList {
				assertion.Nil(response.Respondent, testCase.description, "anonymous questionnaire with respondent")
			}
		} else {
			for _, response := range responseList {
				assertion.NotEqual(response.Respondent, nil, testCase.description, "not anonymous questionnaire with no respondent")
			}
		}
	}
}

func TestPostQuestionnaireResponse(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	responseDueDateTimePlus := time.Now().Add(24 * time.Hour)

	questionnaire := sampleQuestionnaire
	questionnaire.ResponseDueDateTime = &responseDueDateTimePlus
	e := echo.New()
	body, err := json.Marshal(questionnaire)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	questionnaireDetail, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.ResponseDueDateTime = &responseDueDateTimePlus
	questionnaire.IsAnonymous = true
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireDetailAnonymous, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	questionnaire.ResponseDueDateTime = &responseDueDateTimePlus
	questionnaire.IsDuplicateAnswerAllowed = false
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireDetailNoMultipleResponse, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	questionnaire = sampleQuestionnaire
	e = echo.New()
	body, err = json.Marshal(questionnaire)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnaireDetailNoDue, err := q.PostQuestionnaire(ctx, questionnaire)
	require.NoError(t, err)

	type args struct {
		invalidQuestionnaireID bool
		questionnaireDetail    openapi.QuestionnaireDetail
		isNoMultipleResponse   bool
		isAnonymous            bool
		params                 openapi.PostQuestionnaireResponseJSONRequestBody
		userID                 string
		isTimeAfterDue         bool
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

	invalidResponseBodyText := openapi.NewResponseBody{}
	invalidResponseBodyText.QuestionId = *questionnaireDetail.Questions[0].QuestionId
	invalidResponseBodyText.FromResponseBodyText(openapi.ResponseBodyText{
		Answer:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		QuestionType: "Text",
	})
	invalidResponseBodyTextLong := openapi.NewResponseBody{}
	invalidResponseBodyTextLong.QuestionId = *questionnaireDetail.Questions[1].QuestionId
	invalidResponseBodyTextLong.FromResponseBodyTextLong(openapi.ResponseBodyTextLong{
		Answer:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		QuestionType: "TextLong",
	})
	invalidResponseBodyNumber := openapi.NewResponseBody{}
	invalidResponseBodyNumber.QuestionId = *questionnaireDetail.Questions[2].QuestionId
	invalidResponseBodyNumber.FromResponseBodyNumber(openapi.ResponseBodyNumber{
		Answer:       101,
		QuestionType: "Number",
	})
	invalidResponseBodySingleChoice := openapi.NewResponseBody{}
	invalidResponseBodySingleChoice.QuestionId = *questionnaireDetail.Questions[3].QuestionId
	invalidResponseBodySingleChoice.FromNewResponseBodySingleChoice(openapi.NewResponseBodySingleChoice{
		Answer:       5,
		QuestionType: "SingleChoice",
	})
	invalidResponseBodyMultipleChoice := openapi.NewResponseBody{}
	invalidResponseBodyMultipleChoice.QuestionId = *questionnaireDetail.Questions[4].QuestionId
	invalidResponseBodyMultipleChoice.FromNewResponseBodyMultipleChoice(openapi.NewResponseBodyMultipleChoice{
		Answer:       []int{5},
		QuestionType: "MultipleChoice",
	})
	invalidResponseBodyScale := openapi.NewResponseBody{}
	invalidResponseBodyScale.QuestionId = *questionnaireDetail.Questions[5].QuestionId
	invalidResponseBodyScale.FromResponseBodyScale(openapi.ResponseBodyScale{
		Answer:       0,
		QuestionType: "Scale",
	})

	AddQuestionID2SampleResponseMutex.Lock()

	AddQuestionID2SampleResponse(questionnaireDetail.QuestionnaireId)

	testCases := []test{
		{
			description: "valid",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params:              sampleResponse,
				userID:              userOne,
			},
		},
		{
			description: "valid draft",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: true,
				},
				userID: userOne,
			},
		},
		{
			description: "invalid questionnaire id",
			args: args{
				invalidQuestionnaireID: true,
				params:                 sampleResponse,
				userID:                 userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "no enough response",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyScale,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body text",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						invalidResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body text long",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyText,
						invalidResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body number",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						invalidResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body single choice",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						invalidResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body multiple choice",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						invalidResponseBodyMultipleChoice,
						sampleResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "invalid response body scale",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params: openapi.PostQuestionnaireResponseJSONRequestBody{
					Body: []openapi.NewResponseBody{
						sampleResponseBodyText,
						sampleResponseBodyTextLong,
						sampleResponseBodyNumber,
						sampleResponseBodySingleChoice,
						sampleResponseBodyMultipleChoice,
						invalidResponseBodyScale,
					},
					IsDraft: false,
				},
				userID: userOne,
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "is time after due (attention: need to be the last testCase)",
			args: args{
				questionnaireDetail: questionnaireDetail,
				params:              sampleResponse,
				userID:              userOne,
				isTimeAfterDue:      true,
			},
			expect: expect{
				isErr: true,
			},
		},
	}
	tmp_testCase := test{
		description: "question type not match",
		args: args{
			questionnaireDetail: questionnaireDetail,
			params: openapi.PostQuestionnaireResponseJSONRequestBody{
				Body: []openapi.NewResponseBody{
					sampleResponseBodyScale,
					sampleResponseBodyTextLong,
					sampleResponseBodyNumber,
					sampleResponseBodySingleChoice,
					sampleResponseBodyMultipleChoice,
					sampleResponseBodyText,
				},
				IsDraft: false,
			},
			userID: userOne,
		},
		expect: expect{
			isErr: true,
		},
	}
	tmp_testCase.args.params.Body[0].QuestionId = sampleResponseBodyText.QuestionId
	tmp_testCase.args.params.Body[5].QuestionId = sampleResponseBodyScale.QuestionId
	testCases = append(testCases, tmp_testCase)

	AddQuestionID2SampleResponse(questionnaireDetailAnonymous.QuestionnaireId)

	testCases = append(testCases, test{
		description: "valid anonymous",
		args: args{
			questionnaireDetail: questionnaireDetailAnonymous,
			isAnonymous:         true,
			params:              sampleResponse,
			userID:              userOne,
		},
	})

	AddQuestionID2SampleResponse(questionnaireDetailNoMultipleResponse.QuestionnaireId)

	testCases = append(testCases, test{
		description: "is no multiple response",
		args: args{
			questionnaireDetail:  questionnaireDetailNoMultipleResponse,
			isNoMultipleResponse: true,
			params:               sampleResponse,
			userID:               userOne,
		},
	})

	AddQuestionID2SampleResponse(questionnaireDetailNoDue.QuestionnaireId)

	testCases = append(testCases, test{
		description: "questionnaire no due",
		args: args{
			questionnaireDetail: questionnaireDetailNoDue,
			params:              sampleResponse,
			userID:              userOne,
		},
	})

	AddQuestionID2SampleResponseMutex.Unlock()

	for _, testCase := range testCases {
		var questionnaireID int

		if !testCase.args.invalidQuestionnaireID {
			questionnaireID = testCase.args.questionnaireDetail.QuestionnaireId
		} else {
			questionnaireID = 10000
			valid := true
			for valid {
				ctx := context.Background()
				_, _, _, _, _, _, _, _, err := IQuestionnaire.GetQuestionnaireInfo(ctx, questionnaireID)
				if errors.Is(err, model.ErrRecordNotFound) {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting questionnaire info")
				} else {
					questionnaireID *= 10
				}
			}
		}

		if testCase.args.isTimeAfterDue {
			c := context.Background()
			responseDueDateTime := null.Time{}
			responseDueDateTime.Valid = true
			responseDueDateTime.Time = time.Now().Add(-24 * time.Hour)
			err = IQuestionnaire.UpdateQuestionnaire(c, testCase.args.questionnaireDetail.Title, testCase.args.questionnaireDetail.Description, responseDueDateTime, string(testCase.args.questionnaireDetail.ResponseViewableBy), testCase.args.questionnaireDetail.QuestionnaireId, testCase.args.questionnaireDetail.IsPublished, testCase.args.questionnaireDetail.IsAnonymous, testCase.args.questionnaireDetail.IsDuplicateAnswerAllowed)
			require.NoError(t, err)
		}

		e = echo.New()
		body, err = json.Marshal(testCase.args.params)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireID), bytes.NewReader(body))
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		response, err := q.PostQuestionnaireResponse(ctx, questionnaireID, testCase.args.params, testCase.args.userID)

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

		actualResponseBody := make([]openapi.ResponseBody, len(testCase.args.params.Body))
		for i, body := range testCase.args.params.Body {
			actualResponseBody[i].QuestionId = body.QuestionId
			b, err := body.MarshalJSON()
			require.NoError(t, err)
			var responseParsed map[string]interface{}
			err = json.Unmarshal([]byte(b), &responseParsed)
			require.NoError(t, err)
			questionType := responseParsed["question_type"].(string)
			switch questionType {
			case "Text":
				actualResponseBody[i].FromResponseBodyText(openapi.ResponseBodyText{
					Answer:       responseParsed["answer"].(string),
					QuestionType: openapi.ResponseBodyTextQuestionType(questionType),
				})
			case "TextLong":
				actualResponseBody[i].FromResponseBodyTextLong(openapi.ResponseBodyTextLong{
					Answer:       responseParsed["answer"].(string),
					QuestionType: openapi.ResponseBodyTextLongQuestionType(questionType),
				})
			case "Number":
				actualResponseBody[i].FromResponseBodyNumber(openapi.ResponseBodyNumber{
					Answer:       float32(responseParsed["answer"].(float64)),
					QuestionType: openapi.ResponseBodyNumberQuestionType(questionType),
				})
			case "SingleChoice":
				options, err := q.IOption.GetOptions(context.Background(), []int{body.QuestionId})
				require.NoError(t, err)
				actualResponseBody[i].FromResponseBodySingleChoice(openapi.ResponseBodySingleChoice{
					Answer:       options[int(responseParsed["answer"].(float64))].Body,
					QuestionType: openapi.ResponseBodySingleChoiceQuestionType(questionType),
				})
			case "MultipleChoice":
				options, err := q.IOption.GetOptions(context.Background(), []int{body.QuestionId})
				require.NoError(t, err)
				answers := make([]string, len(responseParsed["answer"].([]interface{})))
				for j, answer := range responseParsed["answer"].([]interface{}) {
					answers[j] = options[int(answer.(float64))].Body
				}
				actualResponseBody[i].FromResponseBodyMultipleChoice(openapi.ResponseBodyMultipleChoice{
					Answer:       answers,
					QuestionType: openapi.ResponseBodyMultipleChoiceQuestionType(questionType),
				})
			case "Scale":
				actualResponseBody[i].FromResponseBodyScale(openapi.ResponseBodyScale{
					Answer:       int(responseParsed["answer"].(float64)),
					QuestionType: openapi.ResponseBodyScaleQuestionType(questionType),
				})
			default:
				assertion.Fail("unknown question type", "question type: %s", questionType)
			}
		}
		assertion.Equal(actualResponseBody, response.Body, testCase.description, "response body")
		assertion.Equal(testCase.args.params.IsDraft, response.IsDraft, testCase.description, "is draft")

		assertion.Equal(testCase.args.questionnaireDetail.QuestionnaireId, response.QuestionnaireId, testCase.description, "questionnaire id")
		if !testCase.args.isAnonymous {
			assertion.Equal(testCase.args.userID, *response.Respondent, testCase.description, "respondent")
		} else {
			assertion.Nil(response.Respondent, testCase.description, "respondent")
		}
		assertion.Equal(testCase.args.isAnonymous, *response.IsAnonymous, testCase.description, "is anonymous")

		if testCase.args.params.IsDraft {
			assertion.Equal(response.SubmittedAt, time.Time{}, testCase.description, "submitted at")
		} else {
			assertion.NotEqual(response.SubmittedAt, time.Time{}, testCase.description, "submitted at")
			assertion.NotEqual(response.ModifiedAt, time.Time{}, testCase.description, "modified at")
		}

		AddQuestionID2SampleResponseMutex.Lock()

		AddQuestionID2SampleResponse(questionnaireID)

		e = echo.New()
		body, err = json.Marshal(sampleResponse)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/responses", questionnaireID), bytes.NewReader(body))
		rec = httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(req, rec)
		_, err = q.PostQuestionnaireResponse(ctx, questionnaireID, sampleResponse, testCase.args.userID)

		AddQuestionID2SampleResponseMutex.Unlock()

		if !testCase.args.isNoMultipleResponse {
			assertion.NoError(err, testCase.description, "multiple response")
		} else {
			assertion.Error(err, testCase.description, "multiple response")
		}
	}
}

func TestCreateQuestionnaireMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		questionnaireID int
		title           string
		description     string
		administrators  []string
		resTimeLimit    null.Time
		targets         []string
	}
	type expect struct {
		message string
	}
	type test struct {
		description string
		args
		expect
	}

	tm, err := time.ParseInLocation("2006/01/02 15:04", "2021/10/01 09:06", time.Local)
	if err != nil {
		t.Errorf("failed to parse time: %v", err)
	}

	testCases := []test{
		{
			description: "通常の引数なので問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "questionnaireIDが0でも問題なし",
			args: args{
				questionnaireID: 0,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/0)』が作成されました
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/0`,
			},
		},
		{
			// 実際には発生しないけど念の為
			description: "titleが空文字でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "説明が空文字でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1
#### 説明

#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "administrator複数人でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1", "administrator2"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1,administrator2
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			// 実際には発生しないけど念の為
			description: "administratorがいなくても問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者

#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "回答期限なしでも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.NewTime(time.Time{}, false),
				targets:         []string{"target1"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1
#### 説明
description
#### 回答期限
なし
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "対象者が複数人でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{"target1", "target2"},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1 @target2
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "対象者がいなくても問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    null.TimeFrom(tm),
				targets:         []string{},
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』が作成されました
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
なし
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			message := createQuestionnaireMessage(
				testCase.args.questionnaireID,
				testCase.args.title,
				testCase.args.description,
				testCase.args.administrators,
				testCase.args.resTimeLimit,
				testCase.args.targets,
			)

			assert.Equal(t, testCase.expect.message, message)
		})
	}
}

func TestCreateReminderMessage(t *testing.T) {
	t.Parallel()
	type args struct {
		questionnaireID int
		title           string
		description     string
		administrators  []string
		resTimeLimit    time.Time
		targets         []string
		leftTimeText    string
	}
	type expect struct {
		message string
	}
	type test struct {
		description string
		args
		expect
	}

	tm, err := time.ParseInLocation("2006/01/02 15:04", "2021/10/01 09:06", time.Local)
	if err != nil {
		t.Errorf("failed to parse time: %v", err)
	}

	testCases := []test{
		{
			description: "通常の引数なので問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    tm,
				targets:         []string{"target1"},
				leftTimeText:    "5分",
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』の回答期限が迫っています!
==残り5分です!==
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "questionnaireIDが0でも問題なし",
			args: args{
				questionnaireID: 0,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    tm,
				targets:         []string{"target1"},
				leftTimeText:    "30分",
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/0)』の回答期限が迫っています!
==残り30分です!==
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/0`,
			},
		},
		{
			// 実際には発生しないけど念の為
			description: "titleが空文字でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    tm,
				targets:         []string{"target1"},
				leftTimeText:    "1時間",
			},
			expect: expect{
				message: `### アンケート『[](https://anke-to.trap.jp/questionnaires/1)』の回答期限が迫っています!
==残り1時間です!==
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "説明が空文字でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "",
				administrators:  []string{"administrator1"},
				resTimeLimit:    tm,
				targets:         []string{"target1"},
				leftTimeText:    "1日",
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』の回答期限が迫っています!
==残り1日です!==
#### 管理者
administrator1
#### 説明

#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "administrator複数人でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1", "administrator2"},
				resTimeLimit:    tm,
				targets:         []string{"target1"},
				leftTimeText:    "5分",
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』の回答期限が迫っています!
==残り5分です!==
#### 管理者
administrator1,administrator2
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			// 実際には発生しないけど念の為
			description: "administratorがいなくても問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{},
				resTimeLimit:    tm,
				targets:         []string{"target1"},
				leftTimeText:    "1週間",
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』の回答期限が迫っています!
==残り1週間です!==
#### 管理者

#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
		{
			description: "対象者が複数人でも問題なし",
			args: args{
				questionnaireID: 1,
				title:           "title",
				description:     "description",
				administrators:  []string{"administrator1"},
				resTimeLimit:    tm,
				targets:         []string{"target1", "target2"},
				leftTimeText:    "5分",
			},
			expect: expect{
				message: `### アンケート『[title](https://anke-to.trap.jp/questionnaires/1)』の回答期限が迫っています!
==残り5分です!==
#### 管理者
administrator1
#### 説明
description
#### 回答期限
2021/10/01 09:06
#### 対象者
@target1 @target2
#### 回答リンク
https://anke-to.trap.jp/responses/new/1`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			message := createReminderMessage(
				testCase.args.questionnaireID,
				testCase.args.title,
				testCase.args.description,
				testCase.args.administrators,
				testCase.args.resTimeLimit,
				testCase.args.targets,
				testCase.args.leftTimeText,
			)

			assert.Equal(t, testCase.expect.message, message)
		})
	}
}
