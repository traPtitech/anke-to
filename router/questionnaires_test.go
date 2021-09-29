package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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
	"gopkg.in/guregu/null.v3"
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

			t.Log(rec.Body.String())
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
