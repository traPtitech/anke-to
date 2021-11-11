package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/model/mock_model"
	"github.com/traPtitech/anke-to/traq/mock_traq"
	"gopkg.in/guregu/null.v4"
)

func TestPostAndEditQuestionnaireValidate(t *testing.T) {
	tests := []struct {
		description string
		request     *PostAndEditQuestionnaireRequest
		isErr       bool
	}{
		{
			description: "旧クライアントの一般的なリクエストなのでエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "タイトルが空なのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			isErr: true,
		},
		{
			description: "タイトルが50文字なのでエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "アイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオ",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "タイトルが50文字を超えるのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "アイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオアイウエオア",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			isErr: true,
		},
		{
			description: "descriptionが空でもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "resTimeLimitが設定されていてもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Now(), true),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "resSharedToがadministratorsでもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "administrators",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "resSharedToがrespondentsでもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "respondents",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "resSharedToがadministrators、respondents、publicのいずれでもないのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "test",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			isErr: true,
		},
		{
			description: "targetがnullでもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        nil,
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "targetが32文字でもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{"01234567890123456789012345678901"},
				Administrators: []string{"mazrean"},
			},
		},
		{
			description: "targetが32文字を超えるのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{"012345678901234567890123456789012"},
				Administrators: []string{"mazrean"},
			},
			isErr: true,
		},
		{
			description: "administratorsがいないのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{"01234567890123456789012345678901"},
				Administrators: []string{},
			},
			isErr: true,
		},
		{
			description: "administratorsがnullなのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: nil,
			},
			isErr: true,
		},
		{
			description: "administratorsが32文字でもエラーなし",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"01234567890123456789012345678901"},
			},
		},
		{
			description: "administratorsが32文字を超えるのでエラー",
			request: &PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"012345678901234567890123456789012"},
			},
			isErr: true,
		},
	}

	for _, test := range tests {
		validate := validator.New()

		t.Run(test.description, func(t *testing.T) {
			err := validate.Struct(test.request)

			if test.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetQuestionnaireValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		request     *GetQuestionnairesQueryParam
		isErr       bool
	}{
		{
			description: "一般的なQueryParameterなのでエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortが-created_atでもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "-created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortがtitleでもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "title",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortが-titleでもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "-title",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortがmodified_atでもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "modified_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortが-modified_atでもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "-modified_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Nontargetedをfalseにしてもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "false",
			},
		},
		{
			description: "Sortを空文字にしてもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Searchを空文字にしてもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Pageを空文字にしてもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "",
				Nontargeted: "true",
			},
		},
		{
			description: "Nontargetedを空文字にしてもエラーなし",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "",
			},
		},
		{
			description: "Pageが数字ではないのでエラー",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "xx",
				Nontargeted: "true",
			},
			isErr: true,
		},
		{
			description: "Nontargetedがbool値ではないのでエラー",
			request: &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "arupaka",
			},
			isErr: true,
		},
	}
	for _, test := range tests {
		validate := validator.New()

		t.Run(test.description, func(t *testing.T) {
			err := validate.Struct(test.request)
			if test.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostQuestionnaire(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockTransaction := &model.MockTransaction{}
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	questionnaire := NewQuestionnaire(
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
		mockQuestion,
		mockOption,
		mockScaleLabel,
		mockValidation,
		mockTransaction,
		mockWebhook,
	)

	type expect struct {
		statusCode int
	}
	type test struct {
		description               string
		invalidRequest            bool
		request                   PostAndEditQuestionnaireRequest
		ExecutesCreation          bool
		questionnaireID           int
		InsertQuestionnaireError  error
		InsertTargetsError        error
		InsertAdministratorsError error
		PostMessageError          error
		expect
	}

	testCases := []test{
		{
			description:    "リクエストの形式が誤っているので400",
			invalidRequest: true,
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "validationで落ちるので400",
			request:     PostAndEditQuestionnaireRequest{},
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "resTimeLimitが誤っているので400",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Now().Add(-24*time.Hour), true),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "InsertQuestionnaireがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:         true,
			InsertQuestionnaireError: errors.New("InsertQuestionnaireError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "InsertTargetsがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:   true,
			questionnaireID:    1,
			InsertTargetsError: errors.New("InsertTargetsError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "InsertAdministratorsがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:          true,
			questionnaireID:           1,
			InsertAdministratorsError: errors.New("InsertAdministratorsError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "PostMessageがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  1,
			PostMessageError: errors.New("PostMessageError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "一般的なリクエストなので201",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "questionnaireIDが0でも201",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  0,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "回答期限が設定されていてもでも201",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Now().Add(24*time.Hour), true),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var request io.Reader
			if testCase.invalidRequest {
				request = strings.NewReader("test")
			} else {
				buf := bytes.NewBuffer(nil)
				err := json.NewEncoder(buf).Encode(testCase.request)
				if err != nil {
					t.Errorf("failed to encode request: %v", err)
				}

				request = buf
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/questionnaires", request)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)

			c.Set(validatorKay, validator.New())

			if testCase.ExecutesCreation {
				// 時刻は完全一致しないためその対応
				var mockTimeLimit interface{}
				if testCase.request.ResTimeLimit.Valid {
					mockTimeLimit = gomock.Any()
				} else {
					mockTimeLimit = testCase.request.ResTimeLimit
				}

				mockQuestionnaire.
					EXPECT().
					InsertQuestionnaire(
						c.Request().Context(),
						testCase.request.Title,
						testCase.request.Description,
						mockTimeLimit,
						testCase.request.ResSharedTo,
					).
					Return(testCase.questionnaireID, testCase.InsertQuestionnaireError)

				if testCase.InsertQuestionnaireError == nil {
					mockTarget.
						EXPECT().
						InsertTargets(
							c.Request().Context(),
							testCase.questionnaireID,
							testCase.request.Targets,
						).
						Return(testCase.InsertTargetsError)

					if testCase.InsertTargetsError == nil {
						mockAdministrator.
							EXPECT().
							InsertAdministrators(
								c.Request().Context(),
								testCase.questionnaireID,
								testCase.request.Administrators,
							).
							Return(testCase.InsertAdministratorsError)

						if testCase.InsertAdministratorsError == nil {
							mockWebhook.
								EXPECT().
								PostMessage(gomock.Any()).
								Return(testCase.PostMessageError)
						}
					}
				}
			}

			e.HTTPErrorHandler(questionnaire.PostQuestionnaire(c), c)

			assert.Equal(t, testCase.expect.statusCode, rec.Code, "status code")

			if testCase.expect.statusCode == http.StatusCreated {
				var questionnaire map[string]interface{}
				err := json.NewDecoder(rec.Body).Decode(&questionnaire)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}

				assert.Equal(t, float64(testCase.questionnaireID), questionnaire["questionnaireID"], "questionnaireID")
				assert.Equal(t, testCase.request.Title, questionnaire["title"], "title")
				assert.Equal(t, testCase.request.Description, questionnaire["description"], "description")
				if testCase.request.ResTimeLimit.Valid {
					strResTimeLimit, ok := questionnaire["res_time_limit"].(string)
					assert.True(t, ok, "res_time_limit convert")
					resTimeLimit, err := time.Parse(time.RFC3339, strResTimeLimit)
					assert.NoError(t, err, "res_time_limit parse")

					assert.WithinDuration(t, testCase.request.ResTimeLimit.Time, resTimeLimit, 2*time.Second, "resTimeLimit")
				} else {
					assert.Nil(t, questionnaire["res_time_limit"], "resTimeLimit nil")
				}
				assert.Equal(t, testCase.request.ResSharedTo, questionnaire["res_shared_to"], "resSharedTo")

				strCreatedAt, ok := questionnaire["created_at"].(string)
				assert.True(t, ok, "created_at convert")
				createdAt, err := time.Parse(time.RFC3339, strCreatedAt)
				assert.NoError(t, err, "created_at parse")
				assert.WithinDuration(t, time.Now(), createdAt, time.Second, "created_at")

				strModifiedAt, ok := questionnaire["modified_at"].(string)
				assert.True(t, ok, "modified_at convert")
				modifiedAt, err := time.Parse(time.RFC3339, strModifiedAt)
				assert.NoError(t, err, "modified_at parse")
				assert.WithinDuration(t, time.Now(), modifiedAt, time.Second, "modified_at")

				assert.ElementsMatch(t, testCase.request.Targets, questionnaire["targets"], "targets")
				assert.ElementsMatch(t, testCase.request.Administrators, questionnaire["administrators"], "administrators")
			}
		})
	}
}

func TestPostQuestionByQuestionnaireID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockTransaction := &model.MockTransaction{}
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	questionnaire := NewQuestionnaire(
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
		mockQuestion,
		mockOption,
		mockScaleLabel,
		mockValidation,
		mockTransaction,
		mockWebhook,
	)

	type expect struct {
		statusCode int
	}
	type test struct {
		description           string
		invalidRequest        bool
		request               PostAndEditQuestionRequest
		ExecutesCreation      bool
		questionID            int
		questionnaireID       string
		validator             string
		InsertQuestionError   error
		InsertOptionError     error
		InsertValidationError error
		InsertScaleLabelError error
		CheckNumberValid      error
		expect
	}
	testCases := []test{
		{
			description:    "一般的なリクエストなので201",
			invalidRequest: false,
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Text",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			ExecutesCreation: true,
			questionID:       1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "questionIDが0でも201",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Text",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			ExecutesCreation: true,
			questionID:       0,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "questionnaireIDがstringでも201",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Text",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			questionnaireID: "1",
			ExecutesCreation: true,
			questionID:       1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "QuestionTypeがMultipleChoiceでも201",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "MultipleChoice",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			ExecutesCreation: true,
			questionID:       1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "QuestionTypeがLinearScaleでも201",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "LinearScale",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			ExecutesCreation: true,
			questionID:       1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "QuestionTypeがNumberでも201",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Number",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			ExecutesCreation: true,
			questionID:       1,
			expect: expect{
				statusCode: http.StatusCreated,
			},
		},
		{
			description: "QuestionTypeが存在しないものは400",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "aaa",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			InsertQuestionError: errors.New("InsertQuestionError"),
			ExecutesCreation:    false,
			questionID:          1,
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "InsertValidationがエラーで500",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Text",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			InsertValidationError: errors.New("InsertValidationError"),
			ExecutesCreation:      true,
			questionID:            1,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "CheckNumberValidがエラーで500",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Number",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			CheckNumberValid: errors.New("CheckNumberValidError"),
			ExecutesCreation: false,
			questionID:       1,
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "InsertQuestionがエラーで500",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "Text",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			InsertQuestionError: errors.New("InsertQuestionError"),
			ExecutesCreation:    true,
			questionID:          1,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "InsertScaleLabelErrorがエラーで500",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "LinearScale",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka", "mazrean"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			InsertScaleLabelError: errors.New("InsertScaleLabelError"),
			ExecutesCreation:      true,
			questionID:            1,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "InsertOptionErrorがエラーで500",
			request: PostAndEditQuestionRequest{
				QuestionnaireID: 1,
				QuestionType:    "MultipleChoice",
				QuestionNum:     1,
				PageNum:         1,
				Body:            "発表タイトル",
				IsRequired:      true,
				Options:         []string{"arupaka"},
				ScaleLabelRight: "arupaka",
				ScaleLabelLeft:  "xxarupakaxx",
				ScaleMin:        1,
				ScaleMax:        2,
				RegexPattern:    "^\\d*\\.\\d*$",
				MinBound:        "0",
				MaxBound:        "10",
			},
			InsertOptionError: errors.New("InsertOptionError"),
			ExecutesCreation:  true,
			questionID:        1,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "validatorが\"validator\"ではないので500",
			request: PostAndEditQuestionRequest{},
			validator: "arupaka",
			ExecutesCreation: false,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "正規表現が間違っているので400",
			request: PostAndEditQuestionRequest{
				RegexPattern: `^\/\/(.*?)`,
			},
			ExecutesCreation: false,
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description:    "リクエストの形式が異なっているので400",
			invalidRequest: true,
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "validationで落ちるので400",
			request:     PostAndEditQuestionRequest{},
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			var request io.Reader
			if test.invalidRequest {
				request = strings.NewReader("test")
			} else {
				buf := bytes.NewBuffer(nil)
				err := json.NewEncoder(buf).Encode(test.request)
				if err != nil {
					t.Errorf("failed to encode request: %w", err)
				}

				request = buf
			}

			e := echo.New()
			var req *http.Request
			if test.questionnaireID != "" {
				req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%s/questions", test.questionnaireID), request)
			}else {
				req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d/questions", test.request.QuestionnaireID), request)
			}

			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetParamNames("questionnaireID")
			c.SetParamValues(strconv.Itoa(test.request.QuestionnaireID))

			c.Set(questionnaireIDKey, test.request.QuestionnaireID)
			if test.validator != ""{
				c.Set(test.validator,validator.New())
			}else {
				c.Set(validatorKay,validator.New())
			}

			if test.ExecutesCreation {
				mockQuestion.
					EXPECT().
					InsertQuestion(c.Request().Context(), test.request.QuestionnaireID, test.request.PageNum, test.request.QuestionNum, test.request.QuestionType, test.request.Body, test.request.IsRequired).
					Return(test.questionID, test.InsertQuestionError)
			}
			if test.InsertQuestionError == nil && test.request.QuestionType == "LinearScale" {
				mockScaleLabel.
					EXPECT().
					InsertScaleLabel(c.Request().Context(), test.questionID, model.ScaleLabels{
						ScaleLabelRight: test.request.ScaleLabelRight,
						ScaleLabelLeft:  test.request.ScaleLabelLeft,
						ScaleMin:        test.request.ScaleMin,
						ScaleMax:        test.request.ScaleMax,
					}).Return(test.InsertScaleLabelError)
			}
			if test.InsertQuestionError == nil && (test.request.QuestionType == "MultipleChoice" || test.request.QuestionType == "Checkbox" || test.request.QuestionType == "Dropdown") {
				for i, option := range test.request.Options {
					mockOption.
						EXPECT().
						InsertOption(c.Request().Context(), test.questionID, i+1, option).Return(test.InsertOptionError)
				}
			}
			if test.request.QuestionType == "Number" {
				mockValidation.
					EXPECT().
					CheckNumberValid(test.request.MinBound, test.request.MaxBound).
					Return(test.CheckNumberValid)
			}
			if test.InsertQuestionError == nil && test.CheckNumberValid == nil && (test.request.QuestionType == "Text" || test.request.QuestionType == "Number") {
				mockValidation.
					EXPECT().
					InsertValidation(c.Request().Context(), test.questionID, model.Validations{
						RegexPattern: test.request.RegexPattern,
						MinBound:     test.request.MinBound,
						MaxBound:     test.request.MaxBound,
					}).
					Return(test.InsertValidationError)
			}

			e.HTTPErrorHandler(questionnaire.PostQuestionByQuestionnaireID(c), c)

			assert.Equal(t, test.expect.statusCode, rec.Code, "status code")

			if test.expect.statusCode == http.StatusCreated {
				var question map[string]interface{}
				err := json.NewDecoder(rec.Body).Decode(&question)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}

				assert.Equal(t, float64(test.questionID), question["questionID"], "questionID")
				assert.Equal(t, test.request.QuestionType, question["question_type"], "question_type")
				assert.Equal(t, float64(test.request.QuestionNum), question["question_num"], "question_num")
				assert.Equal(t, float64(test.request.PageNum), question["page_num"], "page_num")
				assert.Equal(t, test.request.Body, question["body"], "body")
				assert.Equal(t, test.request.IsRequired, question["is_required"], "is_required")
				assert.ElementsMatch(t, test.request.Options, question["options"], "options")
				assert.Equal(t, test.request.ScaleLabelRight, question["scale_label_right"], "scale_label_right")
				assert.Equal(t, test.request.ScaleLabelLeft, question["scale_label_left"], "scale_label_left")
				assert.Equal(t, float64(test.request.ScaleMax), question["scale_max"], "scale_max")
				assert.Equal(t, float64(test.request.ScaleMin), question["scale_min"], "scale_min")
				assert.Equal(t, test.request.RegexPattern, question["regex_pattern"], "regex_pattern")
				assert.Equal(t, test.request.MinBound, question["min_bound"], "min_bound")
				assert.Equal(t, test.request.MaxBound, question["max_bound"], "max_bound")

			}
		})
	}
}

func TestEditQuestionnaire(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockTransaction := &model.MockTransaction{}
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	questionnaire := NewQuestionnaire(
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
		mockQuestion,
		mockOption,
		mockScaleLabel,
		mockValidation,
		mockTransaction,
		mockWebhook,
	)

	type expect struct {
		statusCode int
	}
	type test struct {
		description               string
		invalidRequest            bool
		request                   PostAndEditQuestionnaireRequest
		ExecutesCreation          bool
		questionnaireID           int
		InsertQuestionnaireError  error
		DeleteTargetsError        error
		InsertTargetsError        error
		DeleteAdministratorsError error
		InsertAdministratorsError error
		PostMessageError          error
		expect
	}

	testCases := []test{
		{
			description:    "リクエストの形式が誤っているので400",
			invalidRequest: true,
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "validationで落ちるので400",
			request:     PostAndEditQuestionnaireRequest{},
			expect: expect{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "InsertQuestionnaireがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:         true,
			InsertQuestionnaireError: errors.New("InsertQuestionnaireError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "DeleteTargetsがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:   true,
			questionnaireID:    1,
			DeleteTargetsError: errors.New("DeleteTargetsError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "InsertTargetsがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:   true,
			questionnaireID:    1,
			InsertTargetsError: errors.New("InsertTargetsError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "DeleteAdministratorsがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:          true,
			questionnaireID:           1,
			DeleteAdministratorsError: errors.New("DeleteAdministratorsError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "InsertAdministratorsがエラーなので500",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation:          true,
			questionnaireID:           1,
			InsertAdministratorsError: errors.New("InsertAdministratorsError"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "一般的なリクエストなので200",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Time{}, false),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  1,
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			description: "resTimeLimitが現在時刻より前でも200",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Now().Add(-24*time.Hour), true),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  1,
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			description: "回答期限が設定されていてもでも200",
			request: PostAndEditQuestionnaireRequest{
				Title:          "第1回集会らん☆ぷろ募集アンケート",
				Description:    "第1回集会らん☆ぷろ参加者募集",
				ResTimeLimit:   null.NewTime(time.Now().Add(24*time.Hour), true),
				ResSharedTo:    "public",
				Targets:        []string{},
				Administrators: []string{"mazrean"},
			},
			ExecutesCreation: true,
			questionnaireID:  1,
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var request io.Reader
			if testCase.invalidRequest {
				request = strings.NewReader("test")
			} else {
				buf := bytes.NewBuffer(nil)
				err := json.NewEncoder(buf).Encode(testCase.request)
				if err != nil {
					t.Errorf("failed to encode request: %v", err)
				}

				request = buf
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questionnaires/%d", testCase.questionnaireID), request)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetParamNames("questionnaireID")
			c.SetParamValues(strconv.Itoa(testCase.questionnaireID))

			c.Set(questionnaireIDKey, testCase.questionnaireID)
			c.Set(validatorKay, validator.New())

			if testCase.ExecutesCreation {
				// 時刻は完全一致しないためその対応
				var mockTimeLimit interface{}
				if testCase.request.ResTimeLimit.Valid {
					mockTimeLimit = gomock.Any()
				} else {
					mockTimeLimit = testCase.request.ResTimeLimit
				}

				mockQuestionnaire.
					EXPECT().
					UpdateQuestionnaire(
						c.Request().Context(),
						testCase.request.Title,
						testCase.request.Description,
						mockTimeLimit,
						testCase.request.ResSharedTo,
						testCase.questionnaireID,
					).
					Return(testCase.InsertQuestionnaireError)

				if testCase.InsertQuestionnaireError == nil {
					mockTarget.
						EXPECT().
						DeleteTargets(
							c.Request().Context(),
							testCase.questionnaireID,
						).
						Return(testCase.DeleteTargetsError)

					if testCase.DeleteTargetsError == nil {
						mockTarget.
							EXPECT().
							InsertTargets(
								c.Request().Context(),
								testCase.questionnaireID,
								testCase.request.Targets,
							).
							Return(testCase.InsertTargetsError)

						if testCase.InsertTargetsError == nil {
							mockAdministrator.
								EXPECT().
								DeleteAdministrators(
									c.Request().Context(),
									testCase.questionnaireID,
								).
								Return(testCase.DeleteAdministratorsError)

							if testCase.DeleteAdministratorsError == nil {
								mockAdministrator.
									EXPECT().
									InsertAdministrators(
										c.Request().Context(),
										testCase.questionnaireID,
										testCase.request.Administrators,
									).
									Return(testCase.InsertAdministratorsError)
							}
						}
					}
				}
			}

			e.HTTPErrorHandler(questionnaire.EditQuestionnaire(c), c)

			assert.Equal(t, testCase.expect.statusCode, rec.Code, "status code")
		})
	}
}

func TestDeleteQuestionnaire(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)
	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockOption := mock_model.NewMockIOption(ctrl)
	mockScaleLabel := mock_model.NewMockIScaleLabel(ctrl)
	mockValidation := mock_model.NewMockIValidation(ctrl)
	mockTransaction := &model.MockTransaction{}
	mockWebhook := mock_traq.NewMockIWebhook(ctrl)

	questionnaire := NewQuestionnaire(
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
		mockQuestion,
		mockOption,
		mockScaleLabel,
		mockValidation,
		mockTransaction,
		mockWebhook,
	)

	type expect struct {
		statusCode int
	}
	type test struct {
		description               string
		questionnaireID           int
		DeleteQuestionnaireError  error
		DeleteTargetsError        error
		DeleteAdministratorsError error
		expect
	}

	testCases := []test{
		{
			description:               "エラーなしなので200",
			questionnaireID:           1,
			DeleteQuestionnaireError:  nil,
			DeleteTargetsError:        nil,
			DeleteAdministratorsError: nil,
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			description:               "questionnaireIDが0でも200",
			questionnaireID:           0,
			DeleteQuestionnaireError:  nil,
			DeleteTargetsError:        nil,
			DeleteAdministratorsError: nil,
			expect: expect{
				statusCode: http.StatusOK,
			},
		},
		{
			description:               "DeleteQuestionnaireがエラーなので500",
			questionnaireID:           1,
			DeleteQuestionnaireError:  errors.New("error"),
			DeleteTargetsError:        nil,
			DeleteAdministratorsError: nil,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description:               "DeleteTargetsがエラーなので500",
			questionnaireID:           1,
			DeleteQuestionnaireError:  nil,
			DeleteTargetsError:        errors.New("error"),
			DeleteAdministratorsError: nil,
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description:               "DeleteAdministratorsがエラーなので500",
			questionnaireID:           1,
			DeleteQuestionnaireError:  nil,
			DeleteTargetsError:        nil,
			DeleteAdministratorsError: errors.New("error"),
			expect: expect{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/questionnaire/%d", testCase.questionnaireID), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/questionnaires/:questionnaire_id")
			c.SetParamNames("questionnaire_id")
			c.SetParamValues(strconv.Itoa(testCase.questionnaireID))

			c.Set(questionnaireIDKey, testCase.questionnaireID)

			mockQuestionnaire.
				EXPECT().
				DeleteQuestionnaire(
					c.Request().Context(),
					testCase.questionnaireID,
				).
				Return(testCase.DeleteQuestionnaireError)

			if testCase.DeleteQuestionnaireError == nil {
				mockTarget.
					EXPECT().
					DeleteTargets(
						c.Request().Context(),
						testCase.questionnaireID,
					).
					Return(testCase.DeleteTargetsError)

				if testCase.DeleteTargetsError == nil {
					mockAdministrator.
						EXPECT().
						DeleteAdministrators(
							c.Request().Context(),
							testCase.questionnaireID,
						).
						Return(testCase.DeleteAdministratorsError)
				}
			}

			e.HTTPErrorHandler(questionnaire.DeleteQuestionnaire(c), c)

			assert.Equal(t, testCase.expect.statusCode, rec.Code, "status code")
		})
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
