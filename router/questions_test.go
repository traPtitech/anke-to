package router

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type questionRequestBody struct {
	questionnaireID int
	questionType    string
	questionNum     int
	pageNum         int
	body            string
	isRequired      bool
	options         []string
	scaleLabelRight string
	scaleLabelLeft  string
	scaleMin        int
	scaleMax        int
	regexPattern    string
	minBound        string
	maxBound        string
}

type questionResponseBody struct {
	questionID int
	questionRequestBody
}

func TestPostQuestion(t *testing.T) {
	testList := []struct {
		description string
		questionID  int
		request     questionRequestBody
		expectCode  int
	}{
		{
			description: "success",
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "valid regexpattern",
			questionID:  -1,
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "^\\d*\\.\\d*$",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "invalid regexpattern",
			questionID:  -1,
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "*",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusBadRequest,
		},
		{
			description: "valid bound",
			questionID:  -1,
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "0",
				maxBound:        "10",
			},
			expectCode: http.StatusCreated,
		},
		{
			description: "invalid bound",
			questionID:  -1,
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "10",
				maxBound:        "0",
			},
			expectCode: http.StatusBadRequest,
		},
		{
			description: "not found",
			questionID:  -1,
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusBadRequest,
		},
	}
	// + bad request body

	// e = echo.New()
	// for _, testData := range testList {
	// 	fmt.Println("aa")
	// 	rec := createRecorder(userOne, methodPost, makePath("/questions"), typeJSON, testData.request.createQuestionsRequestBody())
	// 	assert.Equal(t, testData.expectCode, rec.Code, testData.description)
	// 	if rec.Code != http.StatusCreated {
	// 		continue
	// 	}
	// }
	// req := httptest.NewRequest(http.MethodPost, makePath("/questions"), strings.NewReader(""))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)
	// validation := model.NewValidation()
	// question := model.NewQuestion()
	// option := model.NewOption()
	// scaleLabel := model.NewScaleLabel()
	// routerQuestion := NewQuestion(validation, question, option, scaleLabel)
	// // Assertions
	// if assert.NoError(t, routerQuestion.PostQuestion(c)) {
	// 	assert.Equal(t, http.StatusCreated, rec.Code)
	// 	assert.Equal(t, echo.MIMEApplicationJSON, rec.Body.String())
	// }
	fmt.Println(testList)
}

func TestEditQuestion(t *testing.T) {
	testList := []struct {
		description string
		request     questionRequestBody
		expectCode  int
	}{
		{
			description: "valid",
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusOK,
		},
		{
			description: "valid regexpattern",
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "^\\d*\\.\\d*$",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusOK,
		},
		{
			description: "invalid regexpattern",
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "*",
				minBound:        "",
				maxBound:        "",
			},
			expectCode: http.StatusBadRequest,
		},
		{
			description: "valid bound",
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "0",
				maxBound:        "10",
			},
			expectCode: http.StatusOK,
		},
		{
			description: "invalid bound",
			request: questionRequestBody{
				questionnaireID: 1,
				questionNum:     1,
				pageNum:         1,
				questionType:    "Text",
				body:            "質問文",
				isRequired:      true,
				options:         []string{"選択肢1"},
				scaleLabelRight: "そう思わない",
				scaleLabelLeft:  "そう思う",
				scaleMin:        1,
				scaleMax:        5,
				regexPattern:    "",
				minBound:        "10",
				maxBound:        "0",
			},
			expectCode: http.StatusBadRequest,
		},
	}
	fmt.Println(testList)
}

func TestDeleteQuestion(t *testing.T) {
	testList := []struct {
		description string
		questionID  int
		expectCode  int
	}{
		{
			description: "valid",
			questionID:  -1,
			expectCode:  http.StatusCreated,
		},
		{
			description: "not found",
			questionID:  -1,
			expectCode:  http.StatusNotFound,
		},
	}
	fmt.Println(testList)
}

func (p *questionRequestBody) createQuestionsRequestBody() string {
	options := make([]string, 0, len(p.options))
	for _, option := range p.options {
		options = append(options, fmt.Sprintf("\"%s\"", option))
	}

	return fmt.Sprintf(
		`{
  "questionnaireID": %v,
  "page_num": %v,
  "question_num": %v,
  "question_type": "%s",
  "body": "%s",
  "is_required": %v,
  "options": [%s],
  "scale_label_right": "%s",
  "scale_label_left": "%s",
  "scale_min": %v,
  "scale_max": %v,
  "regex_pattern": "%s",
  "min_bound": "%s",
  "max_bound": "%s"
}`, p.questionnaireID, p.pageNum, p.questionNum, p.questionType, p.body, p.isRequired, strings.Join(options, ",\n    "), p.scaleLabelRight, p.scaleLabelLeft, p.scaleMin, p.scaleMax, p.regexPattern, p.minBound, p.maxBound)
}
