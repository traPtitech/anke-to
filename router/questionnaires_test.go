package router

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
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
			request:     &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortが-created_atでもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "-created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortがtitleでもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "title",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortが-titleでもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "-title",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortがmodified_atでもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "modified_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Sortが-modified_atでもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "-modified_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Nontargetedをfalseにしてもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "false",
			},
		},
		{
			description: "Sortを空文字にしてもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "",
				Search:      "a",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Searchを空文字にしてもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "",
				Page:        "2",
				Nontargeted: "true",
			},
		},
		{
			description: "Pageを空文字にしてもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "",
				Nontargeted: "true",
			},
		},
		{
			description: "Nontargetedを空文字にしてもエラーなし",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "2",
				Nontargeted: "",
			},
		},
		{
			description: "Pageが数字ではないのでエラー",
			request:     &GetQuestionnairesQueryParam{
				Sort:        "created_at",
				Search:      "a",
				Page:        "xx",
				Nontargeted: "true",
			},
			isErr: true,
		},
		{
			description: "Nontargetedがbool値ではないのでエラー",
			request:     &GetQuestionnairesQueryParam{
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
			}else {
				assert.NoError(t, err)
			}
		})
	}
}