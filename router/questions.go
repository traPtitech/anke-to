package router

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/traPtitech/anke-to/model"
)

// Question Questionの構造体
type Question struct {
	model.IValidation
	model.IQuestion
	model.IOption
	model.IScaleLabel
}

// NewQuestion Questionのコンストラクタ
func NewQuestion(validation model.IValidation, question model.IQuestion, option model.IOption, scaleLabel model.IScaleLabel) *Question {
	return &Question{
		IValidation: validation,
		IQuestion:   question,
		IOption:     option,
		IScaleLabel: scaleLabel,
	}
}

type PostAndEditQuestionRequest struct {
	QuestionnaireID int      `json:"questionnaireID" validate:"min=0"`
	QuestionType    string   `json:"question_type" validate:"required,oneof=Text TextArea Number MultipleChoice Checkbox LinearScale"`
	QuestionNum     int      `json:"question_num" validate:"min=0"`
	PageNum         int      `json:"page_num" validate:"min=0"`
	Body            string   `json:"body" validate:"required"`
	IsRequired      bool     `json:"is_required"`
	Options         []string `json:"options" validate:"required_if=QuestionType Checkbox,required_if=QuestionType MultipleChoice,dive,max=1000"`
	ScaleLabelRight string   `json:"scale_label_right" validate:"max=50"`
	ScaleLabelLeft  string   `json:"scale_label_left" validate:"max=50"`
	ScaleMin        int      `json:"scale_min"`
	ScaleMax        int      `json:"scale_max" validate:"gtecsfield=ScaleMin"`
	RegexPattern    string   `json:"regex_pattern"`
	MinBound        string   `json:"min_bound" validate:"omitempty,number"`
	MaxBound        string   `json:"max_bound" validate:"omitempty,number"`
}

// EditQuestion PATCH /questions/:id
func (q *Question) EditQuestion(c echo.Context) error {
	questionID, err := getQuestionID(c)
	if err != nil {
		c.Logger().Errorf("failed to get question id: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionID: %w", err))
	}

	req := PostAndEditQuestionRequest{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Infof("failed to bind PostAndEditQuestionRequest: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	switch req.QuestionType {
	case "Text":
		//正規表現のチェック
		if _, err := regexp.Compile(req.RegexPattern); err != nil {
			c.Logger().Infof("invalid regex pattern: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	case "Number":
		//数字か，min<=maxになってるか
		if err := q.CheckNumberValid(req.MinBound, req.MaxBound); err != nil {
			c.Logger().Info("invalid number: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	err = q.UpdateQuestion(c.Request().Context(), req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body, req.IsRequired, questionID)
	if err != nil {
		c.Logger().Errorf("failed to update question: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		if err := q.UpdateOptions(c.Request().Context(), req.Options, questionID); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			c.Logger().Errorf("failed to update options: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "LinearScale":
		if err := q.UpdateScaleLabel(c.Request().Context(), questionID,
			model.ScaleLabels{
				ScaleLabelLeft:  req.ScaleLabelLeft,
				ScaleLabelRight: req.ScaleLabelRight,
				ScaleMax:        req.ScaleMax,
				ScaleMin:        req.ScaleMin,
			}); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			c.Logger().Errorf("failed to update scale label: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "Text", "Number":
		if err := q.UpdateValidation(c.Request().Context(), questionID,
			model.Validations{
				RegexPattern: req.RegexPattern,
				MinBound:     req.MinBound,
				MaxBound:     req.MaxBound,
			}); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			c.Logger().Errorf("failed to update validation: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestion DELETE /questions/:id
func (q *Question) DeleteQuestion(c echo.Context) error {
	questionID, err := getQuestionID(c)
	if err != nil {
		c.Logger().Errorf("failed to get question id: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionID: %w", err))
	}

	if err := q.IQuestion.DeleteQuestion(c.Request().Context(), questionID); err != nil {
		c.Logger().Errorf("failed to delete question: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteOptions(c.Request().Context(), questionID); err != nil {
		c.Logger().Errorf("failed to delete options: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteScaleLabel(c.Request().Context(), questionID); err != nil {
		c.Logger().Errorf("failed to delete scale label: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteValidation(c.Request().Context(), questionID); err != nil {
		c.Logger().Errorf("failed to delete validation: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
