package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/model/mock_model"
	"github.com/traPtitech/anke-to/openapi"
	"github.com/traPtitech/anke-to/traq/mock_traq"
)

const (
	userOne   = "userOne"
	userTwo   = "userTwo"
	userThree = "userThree"
	userFour  = "userFour"
	userFive  = "userFive"
)

var (
	groupOne   = uuid.MustParse("3d123d8d-509b-4221-bc1d-82e94deac563") // userOne, userTwo
	groupTwo   = uuid.MustParse("c9e3766a-a307-4100-9df1-24943de026c2") // userThree, userFour
	groupThree = uuid.MustParse("313e4211-715c-4ae6-a247-79a204c50382")
	groupFour  = uuid.MustParse("8564a706-f9d4-46f5-852c-3d44a474902a")
	groupFive  = uuid.MustParse("db7f3c13-eb4f-4890-9773-286dde024d4c")

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
	sampleQuestionSettingsText = openapi.NewQuestion{
		Body:       "質問（テキスト）",
		IsRequired: true,
	}
	sampleQuestionSettingsTextMaxLength := 100
	sampleQuestionSettingsText.FromQuestionSettingsText(openapi.QuestionSettingsText{
		MaxLength:    &sampleQuestionSettingsTextMaxLength,
		QuestionType: openapi.QuestionSettingsTextQuestionTypeText,
	})
	sampleQuestionSettingsTextLong = openapi.NewQuestion{
		Body:       "質問（ロングテキスト）",
		IsRequired: true,
	}
	sampleQuestionSettingsTextLongMaxLength := int(500.0)
	sampleQuestionSettingsTextLong.FromQuestionSettingsTextLong(openapi.QuestionSettingsTextLong{
		MaxLength:    &sampleQuestionSettingsTextLongMaxLength,
		QuestionType: openapi.QuestionSettingsTextLongQuestionTypeTextLong,
	})
	sampleQuestionSettingsNumber = openapi.NewQuestion{
		Body:       "質問（数値）",
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
		Body:       "質問（単一選択）",
		IsRequired: true,
	}
	sampleQuestionSettingsSingleChoice.FromQuestionSettingsSingleChoice(openapi.QuestionSettingsSingleChoice{
		Options:      []string{"選択肢A", "選択肢B", "選択肢C", "選択肢D"},
		QuestionType: openapi.QuestionSettingsSingleChoiceQuestionTypeSingleChoice,
	})
	sampleQuestionSettingsMultipleChoice = openapi.NewQuestion{
		Body:       "質問（複数選択）",
		IsRequired: true,
	}
	sampleQuestionSettingsMultipleChoice.FromQuestionSettingsMultipleChoice(openapi.QuestionSettingsMultipleChoice{
		Options:      []string{"選択肢A", "選択肢B", "選択肢C", "選択肢D"},
		QuestionType: openapi.QuestionSettingsMultipleChoiceQuestionTypeMultipleChoice,
	})
	sampleQeustionsettingsScale = openapi.NewQuestion{
		Body:       "質問（スケール）",
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
		ResponseDueDateTime: &time.Time{},
		ResponseViewableBy:  "anyone",
		Target:              sampleTarget,
		Title:               "第1回集会らん☆ぷろ募集アンケート",
	}
}

func TestGetQuestionnaires(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)

	mockTargetGroup := mock_model.NewMockITargetGroup(ctrl)
	mockTargetUser := mock_model.NewMockITargetUser(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockAdministratorGroup := mock_model.NewMockIAdministratorGroup(ctrl)
	mockAdministratorUser := mock_model.NewMockIAdministratorUser(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockTransaction := mock_model.NewMockITransaction(ctrl)
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	r := NewResponse(mockQuestionnaire, mockRespondent, mockResponse, mockTarget, mockQuestion, mockValidation, mockScaleLabel)
	q := NewQuestionnaire(mockQuestionnaire, mockTarget, mockTargetGroup, mockTargetUser, mockAdministrator, mockAdministratorGroup, mockAdministratorUser, mockQuestion, mockOption, mockScaleLabel, mockValidation, mockTransaction, mockRespondent, mockWebhook, r)

	setupSampleQuestionnaire()

	questionnaire0 := sampleQuestionnaire
	e := echo.New()
	body, err := json.Marshal(questionnaire0)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)
	_, err = q.PostQuestionnaire(ctx, questionnaire0)
	require.NoError(t, err)

	questionnaire1 := sampleQuestionnaire
	questionnaire1.Title = "search test"
	e = echo.New()
	body, err = json.Marshal(questionnaire1)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnairePosted1, err := q.PostQuestionnaire(ctx, questionnaire1)
	require.NoError(t, err)

	questionnaire2 := sampleQuestionnaire
	questionnaire1.Title = "search test"
	e = echo.New()
	body, err = json.Marshal(questionnaire2)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	questionnairePosted2, err := q.PostQuestionnaire(ctx, questionnaire2)
	require.NoError(t, err)

	questionnaire3 := sampleQuestionnaire
	questionnaire1.Title = "abcde"
	e = echo.New()
	body, err = json.Marshal(questionnaire3)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx = e.NewContext(req, rec)
	_, err = q.PostQuestionnaire(ctx, questionnaire3)
	require.NoError(t, err)

	type args struct {
		userID string
		params openapi.GetQuestionnairesParams
	}
	type expect struct {
		isErr               bool
		err                 error
		questionnaireIdList []int
	}
	type test struct {
		description string
		args
		expect
	}

	sortInvalid := (openapi.SortType)("abcde")
	sortCreatedAt := (openapi.SortType)("created_at")
	sortCreatedAtDesc := (openapi.SortType)("-created_at")
	sortTitle := (openapi.SortType)("title")
	sortTitleDesc := (openapi.SortType)("-title")
	sortModifiedAt := (openapi.SortType)("modified_at")
	sortModifiedAtDesc := (openapi.SortType)("-modified_at")
	searchTest := (openapi.SearchInQuery)("search test")
	largePageNum := 100000000
	constTrue := true

	testCases := []test{
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
				questionnaireIdList: []int{
					questionnairePosted1.QuestionnaireId,
					questionnairePosted2.QuestionnaireId,
				},
			},
		},
		{
			description: "only targeting me",
			args: args{
				userID: userOne,
				params: openapi.GetQuestionnairesParams{
					OnlyTargetingMe: &[]openapi.OnlyTargetingMeInQuery{true}[0],
				},
			},
		},
		{
			description: "only targeting me",
			args: args{
				userID: userFive,
				params: openapi.GetQuestionnairesParams{
					OnlyTargetingMe: &constTrue,
				},
			},
		},
		{
			description: "only administrated by me",
			args: args{
				userID: userThree,
				params: openapi.GetQuestionnairesParams{
					OnlyAdministratedByMe: &constTrue,
				},
			},
		},
		{
			description: "only administrated by me",
			args: args{
				userID: userFive,
				params: openapi.GetQuestionnairesParams{
					OnlyAdministratedByMe: &constTrue,
				},
			},
		},
	}

	for _, testCase := range testCases {
		e = echo.New()
		body, err = json.Marshal(testCase.args.params)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodGet, "/questionnaires", bytes.NewReader(body))
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
						assertion.True(preCreatedAt.Before(questionnaire.CreatedAt), testCase.description, "created_at")
					}
					preCreatedAt = questionnaire.CreatedAt
				}
			} else if *testCase.args.params.Sort == "-created_at" {
				var preCreatedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preCreatedAt.IsZero() {
						assertion.True(preCreatedAt.After(questionnaire.CreatedAt), testCase.description, "-created_at")
					}
					preCreatedAt = questionnaire.CreatedAt
				}
			} else if *testCase.args.params.Sort == "title" {
				var preTitle string
				for _, questionnaire := range questionnaireList.Questionnaires {
					if preTitle != "" {
						assertion.True(preTitle > questionnaire.Title, testCase.description, "title")
					}
					preTitle = questionnaire.Title
				}
			} else if *testCase.args.params.Sort == "-title" {
				var preTitle string
				for _, questionnaire := range questionnaireList.Questionnaires {
					if preTitle != "" {
						assertion.True(preTitle < questionnaire.Title, testCase.description, "-title")
					}
					preTitle = questionnaire.Title
				}
			} else if *testCase.args.params.Sort == "modified_at" {
				var preModifiedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preModifiedAt.IsZero() {
						assertion.True(preModifiedAt.Before(questionnaire.ModifiedAt), testCase.description, "modified_at")
					}
					preModifiedAt = questionnaire.ModifiedAt
				}
			} else if *testCase.args.params.Sort == "-modified_at" {
				var preModifiedAt time.Time
				for _, questionnaire := range questionnaireList.Questionnaires {
					if !preModifiedAt.IsZero() {
						assertion.True(preModifiedAt.After(questionnaire.ModifiedAt), testCase.description, "-modified_at")
					}
					preModifiedAt = questionnaire.ModifiedAt
				}
			}
		}

		if len(testCase.expect.questionnaireIdList) > 0 {
			var questionnaireIdList []int
			for _, questionnairSummary := range questionnaireList.Questionnaires {
				questionnaireIdList = append(questionnaireIdList, questionnairSummary.QuestionnaireId)
			}
			sort.Slice(testCase.expect.questionnaireIdList, func(i, j int) bool {
				return testCase.expect.questionnaireIdList[i] < testCase.expect.questionnaireIdList[j]
			})
			sort.Slice(questionnaireIdList, func(i, j int) bool { return questionnaireIdList[i] < questionnaireIdList[j] })
			assertion.Equal(testCase.expect.questionnaireIdList, questionnaireIdList, testCase.description, "questionnaireIdList")
		}

		if testCase.args.params.OnlyTargetingMe != nil || testCase.args.params.OnlyAdministratedByMe != nil {
			for _, questionnaire := range questionnaireList.Questionnaires {
				e = echo.New()
				body, err = json.Marshal("")
				require.NoError(t, err)
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/questionnaire/%d", questionnaire.QuestionnaireId), bytes.NewReader(body))
				rec = httptest.NewRecorder()
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				ctx = e.NewContext(req, rec)

				questionnaireDetail, err := q.GetQuestionnaire(ctx, questionnaire.QuestionnaireId)
				require.NoError(t, err)

				if testCase.args.params.OnlyTargetingMe != nil {
					assertion.NotContains(questionnaireDetail.Target.Users, testCase.args.userID, testCase.description, "OnlyTargetingMe")
				}
				if testCase.args.params.OnlyAdministratedByMe != nil {
					assertion.NotContains(questionnaireDetail.Admin.Users, testCase.args.userID, testCase.description, "OnlyAdministratedByMe")
				}
			}
		}
	}
}

func TestPostQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)

	mockTargetGroup := mock_model.NewMockITargetGroup(ctrl)
	mockTargetUser := mock_model.NewMockITargetUser(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockAdministratorGroup := mock_model.NewMockIAdministratorGroup(ctrl)
	mockAdministratorUser := mock_model.NewMockIAdministratorUser(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockTransaction := mock_model.NewMockITransaction(ctrl)
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	r := NewResponse(mockQuestionnaire, mockRespondent, mockResponse, mockTarget, mockQuestion, mockValidation, mockScaleLabel)
	q := NewQuestionnaire(mockQuestionnaire, mockTarget, mockTargetGroup, mockTargetUser, mockAdministrator, mockAdministratorGroup, mockAdministratorUser, mockQuestion, mockOption, mockScaleLabel, mockValidation, mockTransaction, mockRespondent, mockWebhook, r)

	type args struct {
		userID string
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

	setupSampleQuestionnaire()

	responseDueDateTimeMinus := time.Now().Add(-24 * time.Hour)
	responseDueDateTimePlus := time.Now().Add(24 * time.Hour)

	invalidQuestionSettingsNumber := openapi.NewQuestion{
		Body:       "質問（数値）",
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
		Body:       "質問（スケール）",
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
					ResponseDueDateTime: &time.Time{},
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
		assertion.NotEqual(questionnaireDetail.Admin.Users, testCase.args.params.Admin.Users, "admin users not equal")
		sort.Slice(questionnaireDetail.Admin.Groups, func(i, j int) bool {
			return questionnaireDetail.Admin.Groups[i].String() < questionnaireDetail.Admin.Groups[j].String()
		})
		sort.Slice(testCase.args.params.Admin.Groups, func(i, j int) bool {
			return testCase.args.params.Admin.Groups[i].String() < testCase.args.params.Admin.Groups[j].String()
		})
		assertion.NotEqual(questionnaireDetail.Admin.Groups, testCase.args.params.Admin.Groups, "admin groups not equal")

		assertion.NotEqual(questionnaireDetail.Description, testCase.args.params.Description, "description not equal")
		assertion.NotEqual(questionnaireDetail.IsDuplicateAnswerAllowed, testCase.args.params.IsDuplicateAnswerAllowed, "is duplicate answer allowed not equal")
		assertion.NotEqual(questionnaireDetail.IsAnonymous, testCase.args.params.IsAnonymous, "is anonymous not equal")
		assertion.NotEqual(questionnaireDetail.IsPublished, testCase.args.params.IsPublished, "is published not equal")

		for _, question := range testCase.params.Questions {
			isMatch := false
			for _, questionDetail := range questionnaireDetail.Questions {
				if question.Body == questionDetail.Body && question.IsRequired == questionDetail.IsRequired && question.union == questionDetail.union {
					isMatch = true
					break
				}
			}
			if !isMatch {
				assertion.Fail("question not found")
			}
		}

		assertion.NotEqual(questionnaireDetail.ResponseDueDateTime, testCase.args.params.ResponseDueDateTime, "response due date time not equal")
		assertion.NotEqual(questionnaireDetail.ResponseViewableBy, testCase.args.params.ResponseViewableBy, "response viewable by not equal")

		sort.Slice(questionnaireDetail.Target.Users, func(i, j int) bool { return questionnaireDetail.Target.Users[i] < questionnaireDetail.Target.Users[j] })
		sort.Slice(testCase.args.params.Target.Users, func(i, j int) bool {
			return testCase.args.params.Target.Users[i] < testCase.args.params.Target.Users[j]
		})
		assertion.NotEqual(questionnaireDetail.Target.Users, testCase.args.params.Target.Users, "target users not equal")
		sort.Slice(questionnaireDetail.Target.Groups, func(i, j int) bool {
			return questionnaireDetail.Target.Groups[i].String() < questionnaireDetail.Target.Groups[j].String()
		})
		sort.Slice(testCase.args.params.Target.Groups, func(i, j int) bool {
			return testCase.args.params.Target.Groups[i].String() < testCase.args.params.Target.Groups[j].String()
		})
		assertion.NotEqual(questionnaireDetail.Target.Groups, testCase.args.params.Target.Groups, "target groups not equal")

		assertion.NotEqual(questionnaireDetail.Title, testCase.args.params.Title, "title not equal")

		// todo: set up the mock server for traq group and check if admins and targets are correctly expanded
	}
}

func TestGetQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockResponse := mock_model.NewMockIResponse(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)

	mockTargetGroup := mock_model.NewMockITargetGroup(ctrl)
	mockTargetUser := mock_model.NewMockITargetUser(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockAdministratorGroup := mock_model.NewMockIAdministratorGroup(ctrl)
	mockAdministratorUser := mock_model.NewMockIAdministratorUser(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockTransaction := mock_model.NewMockITransaction(ctrl)
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	r := NewResponse(mockQuestionnaire, mockRespondent, mockResponse, mockTarget, mockQuestion, mockValidation, mockScaleLabel)
	q := NewQuestionnaire(mockQuestionnaire, mockTarget, mockTargetGroup, mockTargetUser, mockAdministrator, mockAdministratorGroup, mockAdministratorUser, mockQuestion, mockOption, mockScaleLabel, mockValidation, mockTransaction, mockRespondent, mockWebhook, r)

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
				_, _, _, _, _, _, _, _, err := mockQuestionnaire.GetQuestionnaireInfo(ctx, questionnaireID)
				if err == errors.New("record not found") {
					valid = false
				} else if err != nil {
					assertion.Fail("unexpected error during getting questionnaire info")
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
		req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/questionnaire/%d", questionnaireID), bytes.NewReader(body))
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

		assertion.NotEqual(questionnaireDetailOrigin, questionnaireDetail, testCase.description, "questionnaireDetail")
	}
}

func TestEditQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaires := []struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	questionnaireIDs := []int{}
	for _, questionnaire := range questionnaires {
		// ctxの作成
		questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnairePosted.QuestionnaireId)
	}

	type args struct {
		questionnaireID int
		params          openapi.EditQuestionnaireJSONRequestBody
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
		// todo: テストケースを追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		err := q.EditQuestionnaire(ctx, testCase.args.questionnaireID, testCase.args.params)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
	}
}

func TestDeleteQuestionnaire(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaires := []struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	questionnaireIDs := []int{}
	for _, questionnaire := range questionnaires {
		// ctxの作成
		questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
		require.NoError(t, err)
		questionnaireIDs = append(questionnaireIDs, questionnairePosted.QuestionnaireId)
	}

	type args struct {
		questionnaireID int
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
		// todo: テストケースを追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		err := q.DeleteQuestionnaire(ctx, testCase.args.questionnaireID)

		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	questionnaire := struct {
		userID string
		params openapi.PostQuestionnaireJSONRequestBody
	}{
		// todo: データを追加
	}

	// ctxの作成

	questionnairePosted, err := q.PostQuestionnaire(ctx, questionnaire.userID, questionnaire.params)
	require.NoError(t, err)

	questionnaireID := questionnairePosted.QuestionnaireId

	responses := []struct {
		questionnaireID int
		params          openapi.PostQuestionnaireResponseJSONRequestBody
		userID          string
	}{
		// todo
	}

	for _, response := range responses {
		// todo: ctxの作成
		_, err := q.PostQuestionnaireResponse(ctx, response.questionnaireID, response.params, response.userID)
		require.NoError(t, err)
	}

	type args struct {
		questionnaireID int
		params          openapi.GetQuestionnaireResponsesParams
		userID          string
	}
	type expect struct {
		isErr     bool
		err       error
		responses openapi.Responses
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		// todo: テストケースを追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		responses, err := q.GetQuestionnaireResponses(ctx, testCase.args.questionnaireID, testCase.args.params, testCase.args.userID)

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

		assertion.Equal(testCase.expect.responses, responses, testCase.description, "responses")
	}
}

func TestPostQuestionnaireResponse(t *testing.T) {

	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewQuestionnaire()

	type args struct {
		questionnaireID int
		params          openapi.PostQuestionnaireResponseJSONRequestBody
		userID          string
	}
	type expect struct {
		isErr    bool
		err      error
		response openapi.Response
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		//	todo: テストケースの追加
	}

	for _, testCase := range testCases {
		// todo: ctxの作成

		response, err := q.PostQuestionnaireResponse(ctx, testCase.args.questionnaireID, testCase.args.params, testCase.userID)

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

		assertion.Equal(testCase.expect.response, response, testCase.description, "response")
	}
}
