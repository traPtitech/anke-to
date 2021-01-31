package router

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo"

	"github.com/traPtitech/anke-to/model"
)

// Question Questionの構造体
type Question struct {
	model.ValidationRepository
	model.QuestionRepository
	model.OptionRepository
	model.ScaleLabelRepository
}

// PostQuestion POST /questions
func (q *Question) PostQuestion(c echo.Context) error {
	req := struct {
		QuestionnaireID int      `json:"questionnaireID"`
		QuestionType    string   `json:"question_type"`
		QuestionNum     int      `json:"question_num"`
		PageNum         int      `json:"page_num"`
		Body            string   `json:"body"`
		IsRequired      bool     `json:"is_required"`
		Options         []string `json:"options"`
		ScaleLabelRight string   `json:"scale_label_right"`
		ScaleLabelLeft  string   `json:"scale_label_left"`
		ScaleMin        int      `json:"scale_min"`
		ScaleMax        int      `json:"scale_max"`
		RegexPattern    string   `json:"regex_pattern"`
		MinBound        string   `json:"min_bound"`
		MaxBound        string   `json:"max_bound"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	switch req.QuestionType {
	case "Text":
		//正規表現のチェック
		if _, err := regexp.Compile(req.RegexPattern); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	case "Number":
		//数字か，min<=maxになってるか
		if err := q.CheckNumberValid(req.MinBound, req.MaxBound); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	lastID, err := q.InsertQuestion(req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body, req.IsRequired)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		for i, v := range req.Options {
			if err := q.InsertOption(lastID, i+1, v); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	case "LinearScale":
		if err := q.InsertScaleLabel(lastID,
			model.ScaleLabels{
				ScaleLabelLeft:  req.ScaleLabelLeft,
				ScaleLabelRight: req.ScaleLabelRight,
				ScaleMax:        req.ScaleMax,
				ScaleMin:        req.ScaleMin,
			}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "Text", "Number":
		if err := q.InsertValidation(lastID,
			model.Validations{
				RegexPattern: req.RegexPattern,
				MinBound:     req.MinBound,
				MaxBound:     req.MaxBound,
			}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionID":        int(lastID),
		"questionnaireID":   req.QuestionnaireID,
		"question_type":     req.QuestionType,
		"question_num":      req.QuestionNum,
		"page_num":          req.PageNum,
		"body":              req.Body,
		"is_required":       req.IsRequired,
		"options":           req.Options,
		"scale_label_right": req.ScaleLabelRight,
		"scale_label_left":  req.ScaleLabelLeft,
		"scale_max":         req.ScaleMax,
		"scale_min":         req.ScaleMin,
		"regex_pattern":     req.RegexPattern,
		"min_bound":         req.MinBound,
		"max_bound":         req.MaxBound,
	})
}

// EditQuestion PATCH /questions/:id
func (q *Question) EditQuestion(c echo.Context) error {
	questionID, err := getQuestionID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionID: %w", err))
	}

	req := struct {
		QuestionnaireID int      `json:"questionnaireID"`
		QuestionType    string   `json:"question_type"`
		QuestionNum     int      `json:"question_num"`
		PageNum         int      `json:"page_num"`
		Body            string   `json:"body"`
		IsRequired      bool     `json:"is_required"`
		Options         []string `json:"options"`
		ScaleLabelRight string   `json:"scale_label_right"`
		ScaleLabelLeft  string   `json:"scale_label_left"`
		ScaleMax        int      `json:"scale_max"`
		ScaleMin        int      `json:"scale_min"`
		RegexPattern    string   `json:"regex_pattern"`
		MinBound        string   `json:"min_bound"`
		MaxBound        string   `json:"max_bound"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	switch req.QuestionType {
	case "Text":
		//正規表現のチェック
		if _, err := regexp.Compile(req.RegexPattern); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	case "Number":
		//数字か，min<=maxになってるか
		if err := q.CheckNumberValid(req.MinBound, req.MaxBound); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	if err := q.UpdateQuestion(req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body,
		req.IsRequired, questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		if err := q.UpdateOptions(req.Options, questionID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "LinearScale":
		if err := q.UpdateScaleLabel(questionID,
			model.ScaleLabels{
				ScaleLabelLeft:  req.ScaleLabelLeft,
				ScaleLabelRight: req.ScaleLabelRight,
				ScaleMax:        req.ScaleMax,
				ScaleMin:        req.ScaleMin,
			}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "Text", "Number":
		if err := q.UpdateValidation(questionID,
			model.Validations{
				RegexPattern: req.RegexPattern,
				MinBound:     req.MinBound,
				MaxBound:     req.MaxBound,
			}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestion DELETE /questions/:id
func (q *Question) DeleteQuestion(c echo.Context) error {
	questionID, err := getQuestionID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionID: %w", err))
	}

	if err := q.QuestionRepository.DeleteQuestion(questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteOptions(questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteScaleLabel(questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteValidation(questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
