package router

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"github.com/traPtitech/anke-to/model"
)

// GetQuestions GET /questions
func GetQuestions(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	allquestions, err := model.GetQuestions(questionnaireID)
	if err != nil {
		if errors.Is(err, &model.QuestionNotFoundError{}) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	type questionInfo struct {
		QuestionID      int      `json:"questionID"`
		PageNum         int      `json:"page_num"`
		QuestionNum     int      `json:"question_num"`
		QuestionType    string   `json:"question_type"`
		Body            string   `json:"body"`
		IsRequired      bool     `json:"is_required"`
		CreatedAt       string   `json:"created_at"`
		Options         []string `json:"options"`
		ScaleLabelRight string   `json:"scale_label_right"`
		ScaleLabelLeft  string   `json:"scale_label_left"`
		ScaleMin        int      `json:"scale_min"`
		ScaleMax        int      `json:"scale_max"`
		RegexPattern    string   `json:"regex_pattern"`
		MinBound        string   `json:"min_bound"`
		MaxBound        string   `json:"max_bound"`
	}
	var ret []questionInfo

	for _, v := range allquestions {
		options := []string{}
		scalelabel := model.ScaleLabels{}
		validation := model.Validations{}
		var err error
		switch v.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			options, err = model.GetOptions(c, v.ID)
		case "LinearScale":
			scalelabel, err = model.GetScaleLabels(c, v.ID)
		case "Text", "Number":
			validation, err = model.GetValidations(c, v.ID)
		}
		if err != nil {
			return err
		}

		ret = append(ret,
			questionInfo{
				QuestionID:      v.ID,
				PageNum:         v.PageNum,
				QuestionNum:     v.QuestionNum,
				QuestionType:    v.Type,
				Body:            v.Body,
				IsRequired:      v.IsRequired,
				CreatedAt:       v.CreatedAt.Format(time.RFC3339),
				Options:         options,
				ScaleLabelRight: scalelabel.ScaleLabelRight,
				ScaleLabelLeft:  scalelabel.ScaleLabelLeft,
				ScaleMin:        scalelabel.ScaleMin,
				ScaleMax:        scalelabel.ScaleMax,
				RegexPattern:    validation.RegexPattern,
				MinBound:        validation.MinBound,
				MaxBound:        validation.MaxBound,
			})
	}

	return c.JSON(http.StatusOK, ret)
}

// PostQuestion POST /questions
func PostQuestion(c echo.Context) error {
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
		if err := model.CheckNumberValid(req.MinBound, req.MaxBound); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	lastID, err := model.InsertQuestion(
		req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body, req.IsRequired)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		for i, v := range req.Options {
			if err := model.InsertOption(c, lastID, i+1, v); err != nil {
				return err
			}
		}
	case "LinearScale":
		if err := model.InsertScaleLabels(c, lastID,
			model.ScaleLabels{
				ScaleLabelLeft:  req.ScaleLabelLeft,
				ScaleLabelRight: req.ScaleLabelRight,
				ScaleMax:        req.ScaleMax,
				ScaleMin:        req.ScaleMin,
			}); err != nil {
			return err
		}
	case "Text", "Number":
		if err := model.InsertValidations(c, lastID,
			model.Validations{
				RegexPattern: req.RegexPattern,
				MinBound:     req.MinBound,
				MaxBound:     req.MaxBound,
			}); err != nil {
			return err
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
func EditQuestion(c echo.Context) error {
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
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
		if err := model.CheckNumberValid(req.MinBound, req.MaxBound); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	if err := model.UpdateQuestion(
		req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body,
		req.IsRequired, questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		if err := model.UpdateOptions(c, req.Options, questionID); err != nil {
			return err
		}
	case "LinearScale":
		if err := model.UpdateScaleLabels(c, questionID,
			model.ScaleLabels{
				ScaleLabelLeft:  req.ScaleLabelLeft,
				ScaleLabelRight: req.ScaleLabelRight,
				ScaleMax:        req.ScaleMax,
				ScaleMin:        req.ScaleMin,
			}); err != nil {
			return err
		}
	case "Text", "Number":
		if err := model.UpdateValidations(c, questionID,
			model.Validations{
				RegexPattern: req.RegexPattern,
				MinBound:     req.MinBound,
				MaxBound:     req.MaxBound,
			}); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestion DELETE /questions/:id
func DeleteQuestion(c echo.Context) error {
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := model.DeleteQuestion(questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := model.DeleteOptions(c, questionID); err != nil {
		return err
	}

	if err := model.DeleteScaleLabels(c, questionID); err != nil {
		return err
	}

	if err := model.DeleteValidations(c, questionID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
