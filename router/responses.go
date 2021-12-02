package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"gopkg.in/guregu/null.v4"

	"github.com/traPtitech/anke-to/model"
)

// Response Responseの構造体
type Response struct {
	model.IQuestionnaire
	model.IValidation
	model.IScaleLabel
	model.IRespondent
	model.IResponse
}

// NewResponse Responseのコンストラクタ
func NewResponse(questionnaire model.IQuestionnaire, validation model.IValidation, scaleLabel model.IScaleLabel, respondent model.IRespondent, response model.IResponse) *Response {
	return &Response{
		IQuestionnaire: questionnaire,
		IValidation:    validation,
		IScaleLabel:    scaleLabel,
		IRespondent:    respondent,
		IResponse:      response,
	}
}

// Responses 質問に対する回答一覧の構造体
type Responses struct {
	ID          int                  `json:"questionnaireID" validate:"min=0"`
	Temporarily bool                 `json:"temporarily"`
	Body        []model.ResponseBody `json:"body" validate:"required,dive"`
}

// PostResponse POST /responses
func (r *Response) PostResponse(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		c.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	req := Responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Infof("failed to bind Responses: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	validate, err := getValidator(c)
	if err != nil {
		c.Logger().Errorf("failed to get validator: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(), req)
	if err != nil {
		c.Logger().Infof("failed to validate: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	limit, err := r.GetQuestionnaireLimit(c.Request().Context(), req.ID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info("questionnaire not found")
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		c.Logger().Errorf("failed to get questionnaire limit: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 回答期限を過ぎた回答は許可しない
	if limit.Valid && limit.Time.Before(time.Now()) {
		c.Logger().Info("expired questionnaire")
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	// validationsのパターンマッチ
	questionIDs := make([]int, 0, len(req.Body))
	QuestionTypes := make(map[int]model.ResponseBody, len(req.Body))

	for _, body := range req.Body {
		questionIDs = append(questionIDs, body.QuestionID)
		QuestionTypes[body.QuestionID] = body
	}

	validations, err := r.GetValidations(c.Request().Context(), questionIDs)
	if err != nil {
		c.Logger().Errorf("failed to get validations: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// パターンマッチしてエラーなら返す
	for _, validation := range validations {
		body := QuestionTypes[validation.QuestionID]
		switch body.QuestionType {
		case "Number":
			if err := r.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, model.ErrInvalidNumber) {
					c.Logger().Errorf("invalid number: %+v", err)
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				c.Logger().Infof("invalid number: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Text":
			if err := r.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, model.ErrTextMatching) {
					c.Logger().Infof("invalid text: %+v", err)
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				c.Logger().Errorf("invalid text: %+v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}

	scaleLabelIDs := []int{}
	for _, body := range req.Body {
		switch body.QuestionType {
		case "LinearScale":
			scaleLabelIDs = append(scaleLabelIDs, body.QuestionID)
		}
	}

	scaleLabels, err := r.GetScaleLabels(c.Request().Context(), scaleLabelIDs)
	if err != nil {
		c.Logger().Errorf("failed to get scale labels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	scaleLabelMap := make(map[int]*model.ScaleLabels, len(scaleLabels))
	for _, label := range scaleLabels {
		scaleLabelMap[label.QuestionID] = &label
	}

	// LinearScaleのパターンマッチ
	for _, body := range req.Body {
		switch body.QuestionType {
		case "LinearScale":
			label, ok := scaleLabelMap[body.QuestionID]
			if !ok {
				label = &model.ScaleLabels{}
			}
			if err := r.CheckScaleLabel(*label, body.Body.ValueOrZero()); err != nil {
				c.Logger().Infof("failed to check scale label: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}
	}

	var submittedAt time.Time
	//一時保存のときはnull
	if req.Temporarily {
		submittedAt = time.Time{}
	} else {
		submittedAt = time.Now()
	}

	responseID, err := r.InsertRespondent(c.Request().Context(), userID, req.ID, null.NewTime(submittedAt, !req.Temporarily))
	if err != nil {
		c.Logger().Errorf("failed to insert respondent: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	responseMetas := make([]*model.ResponseMeta, 0, len(req.Body))
	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				responseMetas = append(responseMetas, &model.ResponseMeta{
					QuestionID: body.QuestionID,
					Data:       option,
				})
			}
		default:
			responseMetas = append(responseMetas, &model.ResponseMeta{
				QuestionID: body.QuestionID,
				Data:       body.Body.ValueOrZero(),
			})
		}
	}

	err = r.InsertResponses(c.Request().Context(), responseID, responseMetas)
	if err != nil {
		c.Logger().Errorf("failed to insert responses: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert responses: %w", err))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"responseID":      responseID,
		"questionnaireID": req.ID,
		"temporarily":     req.Temporarily,
		"submitted_at":    submittedAt,
		"body":            req.Body,
	})
}

// GetResponse GET /responses/:responseID
func (r *Response) GetResponse(c echo.Context) error {
	strResponseID := c.Param("responseID")
	responseID, err := strconv.Atoi(strResponseID)
	if err != nil {
		c.Logger().Infof("failed to convert responseID to int: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to parse responseID(%s) to integer: %w", strResponseID, err))
	}

	respondentDetail, err := r.GetRespondentDetail(c.Request().Context(), responseID)
	if errors.Is(err, model.ErrRecordNotFound) {
		c.Logger().Infof("response not found: %+v", err)
		return echo.NewHTTPError(http.StatusNotFound, "response not found")
	}
	if err != nil {
		c.Logger().Errorf("failed to get respondent detail: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, respondentDetail)
}

// EditResponse PATCH /responses/:responseID
func (r *Response) EditResponse(c echo.Context) error {
	responseID, err := getResponseID(c)
	if err != nil {
		c.Logger().Errorf("failed to get responseID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get responseID: %w", err))
	}

	req := Responses{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Infof("failed to bind Responses: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	limit, err := r.GetQuestionnaireLimit(c.Request().Context(), req.ID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("questionnaire not found: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		c.Logger().Errorf("failed to get questionnaire limit: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 回答期限を過ぎた回答は許可しない
	if limit.Valid && limit.Time.Before(time.Now()) {
		c.Logger().Info("expired questionnaire")
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	// validationsのパターンマッチ
	questionIDs := make([]int, 0, len(req.Body))
	QuestionTypes := make(map[int]model.ResponseBody, len(req.Body))

	for _, body := range req.Body {
		questionIDs = append(questionIDs, body.QuestionID)
		QuestionTypes[body.QuestionID] = body
	}

	validations, err := r.GetValidations(c.Request().Context(), questionIDs)
	if err != nil {
		c.Logger().Errorf("failed to get validations: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// パターンマッチしてエラーなら返す
	for _, validation := range validations {
		body := QuestionTypes[validation.QuestionID]
		switch body.QuestionType {
		case "Number":
			if err := r.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, model.ErrInvalidNumber) {
					c.Logger().Errorf("invalid number: %+v", err)
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				c.Logger().Infof("invalid number: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Text":
			if err := r.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, model.ErrTextMatching) {
					c.Logger().Infof("invalid text: %+v", err)
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				c.Logger().Errorf("invalid text: %+v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}

	scaleLabelIDs := []int{}
	for _, body := range req.Body {
		switch body.QuestionType {
		case "LinearScale":
			scaleLabelIDs = append(scaleLabelIDs, body.QuestionID)
		}
	}

	scaleLabels, err := r.GetScaleLabels(c.Request().Context(), scaleLabelIDs)
	if err != nil {
		c.Logger().Errorf("failed to get scale labels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	scaleLabelMap := make(map[int]*model.ScaleLabels, len(scaleLabels))
	for _, label := range scaleLabels {
		scaleLabelMap[label.QuestionID] = &label
	}

	// LinearScaleのパターンマッチ
	for _, body := range req.Body {
		switch body.QuestionType {
		case "LinearScale":
			label, ok := scaleLabelMap[body.QuestionID]
			if !ok {
				label = &model.ScaleLabels{}
			}
			if err := r.CheckScaleLabel(*label, body.Body.ValueOrZero()); err != nil {
				c.Logger().Infof("invalid scale label: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}
	}

	if !req.Temporarily {
		err := r.UpdateSubmittedAt(c.Request().Context(), responseID)
		if err != nil {
			c.Logger().Errorf("failed to update submitted at: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to update sbmitted_at: %w", err))
		}
	}

	//全消し&追加(レコード数爆発しそう)
	if err := r.IResponse.DeleteResponse(c.Request().Context(), responseID); err != nil {
		c.Logger().Errorf("failed to delete response: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	responseMetas := make([]*model.ResponseMeta, 0, len(req.Body))
	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				responseMetas = append(responseMetas, &model.ResponseMeta{
					QuestionID: body.QuestionID,
					Data:       option,
				})
			}
		default:
			responseMetas = append(responseMetas, &model.ResponseMeta{
				QuestionID: body.QuestionID,
				Data:       body.Body.ValueOrZero(),
			})
		}
	}

	err = r.InsertResponses(c.Request().Context(), responseID, responseMetas)
	if err != nil {
		c.Logger().Errorf("failed to insert responses: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert responses: %w", err))
	}

	return c.NoContent(http.StatusOK)
}

// DeleteResponse DELETE /responses/:responseID
func (r *Response) DeleteResponse(c echo.Context) error {
	responseID, err := getResponseID(c)
	if err != nil {
		c.Logger().Errorf("failed to get response id: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get responseID: %w", err))
	}

	limit, err := r.GetQuestionnaireLimitByResponseID(c.Request().Context(), responseID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("response not found: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("failed to find limit of responseID:%d(error: %w)", responseID, err))
		}
		c.Logger().Errorf("failed to get questionnaire limit: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get limit of responseID:%d(error: %w)", responseID, err))
	}

	// 回答期限を過ぎた回答の削除は許可しない
	if limit.Valid && limit.Time.Before(time.Now()) {
		c.Logger().Info("expired response")
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	err = r.DeleteRespondent(c.Request().Context(), responseID)
	if err != nil {
		c.Logger().Errorf("failed to delete respondent: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = r.IResponse.DeleteResponse(c.Request().Context(), responseID)
	if err != nil {
		c.Logger().Errorf("failed to delete response: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
