package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

type resSharedTo string

const (
	sharePublic         resSharedTo = "public"
	shareRespondents    resSharedTo = "respondents"
	shareAdministrators resSharedTo = "administrators"
)

type postQuestionnaireRequestBody struct {
	user           users
	title          string
	description    string
	administrators []string
	targets        []string
	resTimeLimit   null.Time
	resSharedTo    string
}

type postQuestionnaireResponseBody struct {
	QuestionnaireID int       `json:"questionnaireID"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ResTimeLimit    null.Time `json:"res_time_limit"`
	ResSharedTo     string    `json:"res_shared_to"`
	Targets         []string  `json:"targets"`
	Administrators  []string  `json:"administrators"`
}

func TestPostQuestionnaires(t *testing.T) {
	testList := []struct {
		description string
		request     postQuestionnaireRequestBody
		expectCode  int
	}{
		{
			description: "normal request(public)",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "unauthorized user",
			request: postQuestionnaireRequestBody{
				user:           userUnAuthorized,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{userUnAuthorized},
				targets:        []string{},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			description: "res_shared_to: respondents",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{},
				resSharedTo:    string(shareRespondents),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "res_shared_to: administrators",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{},
				resSharedTo:    string(shareAdministrators),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "the administrator is a deferent user",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userTwo)},
				targets:        []string{},
				resSharedTo:    string(shareAdministrators),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "invalid administrator",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userUnAuthorized)},
				targets:        []string{},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated, //いずれ400を返すようにする
		},
		{
			description: "target(user)",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{string(userTwo)},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "target(traP)",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{"traP"},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "invalid target",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{userUnAuthorized},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated, //いずれ400を返すようにする
		},
		{
			description: "empty title",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "",
				description:    "第1回メンバー集会でのらん☆ぷろで発表したい人を募集します らん☆ぷろで発表したい人あつまれー！",
				administrators: []string{string(userOne)},
				targets:        []string{},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "empty description",
			request: postQuestionnaireRequestBody{
				user:           userOne,
				title:          "第1回集会らん☆ぷろ募集アンケート",
				description:    "",
				administrators: []string{string(userOne)},
				targets:        []string{},
				resSharedTo:    string(sharePublic),
			},
			expectCode: http.StatusCreated,
		},
	}

	for _, testData := range testList {
		rec := createRecorder(testData.request.user, methodPost, makePath("/questionnaires"), typeJSON, testData.request.createPostQuestionnairesRequestBody())

		assert.Equal(t, testData.expectCode, rec.Code, testData.description)
		if rec.Code != http.StatusOK {
			continue
		}

		body := postQuestionnaireResponseBody{}
		err := json.NewDecoder(rec.Body).Decode(&body)
		if err != nil {
			t.Fatalf("description: %s\njson decode error: %s", testData.description, err.Error())
		}

		assert.Equal(t, testData.request.title, body.Title, testData.description)
		assert.Equal(t, testData.request.description, body.Description, testData.description)
		assert.Equal(t, testData.request.resTimeLimit, body.ResTimeLimit, testData.description)
		assert.Equal(t, testData.request.resSharedTo, body.ResSharedTo, testData.description)
		assertion := assert.New(t)
		assertion.ElementsMatch(testData.request.targets, body.Targets, testData.description)
		assertion.ElementsMatch(testData.request.administrators, body.Administrators, testData.description)
	}
}

func (p *postQuestionnaireRequestBody) createPostQuestionnairesRequestBody() string {
	targets := make([]string, 0, len(p.targets))
	for _, target := range p.targets {
		targets = append(targets, fmt.Sprintf("\"%s\"", target))
	}

	administrators := make([]string, 0, len(p.administrators))
	for _, administrator := range p.administrators {
		administrators = append(administrators, fmt.Sprintf("\"%s\"", administrator))
	}

	if !p.resTimeLimit.Valid {
		return fmt.Sprintf(
			`{
  "title": "%s",
  "description": "%s",
  "res_shared_to": "%s",
  "targets": [%s],
  "administrators": [%s]
}`, p.title, p.description, p.resSharedTo, strings.Join(targets, ",\n    "), strings.Join(administrators, ",\n    "))
	}

	return fmt.Sprintf(
		`{
  "title": "%s",
  "description": "%s",
  "res_time_limit": "%s",
  "res_shared_to": "%s",
  "targets": [%s],
  "administrators": [%s]
}`, p.title, p.description, p.resTimeLimit.Time.Format(time.RFC3339), p.resSharedTo, strings.Join(targets, ",\n    "), strings.Join(administrators, ",\n    "))
}
