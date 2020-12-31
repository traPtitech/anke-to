package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/traPtitech/anke-to/model"
)

// PostResponse POST /responses
func PostResponse(c echo.Context) error {

	req := model.Responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	limit, err := model.GetQuestionnaireLimit(req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 回答期限を過ぎた回答は許可しない
	if limit != "NULL" && limit < time.Now().Format(time.RFC3339) {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	// validationsのパターンマッチ
	questionIDs := make([]int, 0, len(req.Body))
	QuestionTypes := make(map[int]model.ResponseBody, len(req.Body))

	for _, body := range req.Body {
		questionIDs = append(questionIDs, body.QuestionID)
		QuestionTypes[body.QuestionID] = body
	}

	validations, err := model.GetValidations(questionIDs)

	// パターンマッチしてエラーなら返す
	for _, validation := range validations {
		body := QuestionTypes[validation.QuestionID]
		switch body.QuestionType {
		case "Number":
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := model.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.NumberValidError{}) {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Text":
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := model.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.TextMatchError{}) {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
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

	scaleLabels, err := model.GetScaleLabels(scaleLabelIDs)
	if err != nil {
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
			if err := model.CheckScaleLabel(*label, body.Body.ValueOrZero()); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}
	}

	userID := model.GetUserID(c)
	responseID, err := model.InsertRespondent(userID, req.ID, req.SubmittedAt)
	if err != nil {
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

	err = model.InsertResponses(responseID, responseMetas)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert responses: %w", err))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"responseID":      responseID,
		"questionnaireID": req.ID,
		"submitted_at":    req.SubmittedAt,
		"body":            req.Body,
	})
}

// GetResponse GET /responses/:responseID
func GetResponse(c echo.Context) error {
	strResponseID := c.Param("responseID")
	responseID, err := strconv.Atoi(strResponseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to parse responseID(%s) to integer: %w", strResponseID, err))
	}

	respondentDetail, err := model.GetRespondentDetail(responseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, respondentDetail)
}

// EditResponse PATCH /responses/:responseID
func EditResponse(c echo.Context) error {
	responseID, err := getResponseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get responseID: %w", err))
	}

	req := model.Responses{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	limit, err := model.GetQuestionnaireLimit(req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 回答期限を過ぎた回答は許可しない
	if limit != "NULL" && limit < time.Now().Format(time.RFC3339) {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	// validationsのパターンマッチ
	questionIDs := make([]int, 0, len(req.Body))
	QuestionTypes := make(map[int]model.ResponseBody, len(req.Body))

	for _, body := range req.Body {
		questionIDs = append(questionIDs, body.QuestionID)
		QuestionTypes[body.QuestionID] = body
	}

	validations, err := model.GetValidations(questionIDs)

	// パターンマッチしてエラーなら返す
	for _, validation := range validations {
		body := QuestionTypes[validation.QuestionID]
		switch body.QuestionType {
		case "Number":
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := model.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.NumberValidError{}) {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Text":
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := model.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.TextMatchError{}) {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
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

	scaleLabels, err := model.GetScaleLabels(scaleLabelIDs)
	if err != nil {
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
			if err := model.CheckScaleLabel(*label, body.Body.ValueOrZero()); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}
	}

	if req.SubmittedAt.Valid {
		err := model.UpdateSubmittedAt(responseID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to update sbmitted_at: %w", err))
		}
	}

	//全消し&追加(レコード数爆発しそう)
	if err := model.DeleteResponse(responseID); err != nil {
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

	err = model.InsertResponses(responseID, responseMetas)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert responses: %w", err))
	}

	return c.NoContent(http.StatusOK)
}

// DeleteResponse DELETE /responses/:responseID
func DeleteResponse(c echo.Context) error {
	responseID, err := getResponseID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get responseID: %w", err))
	}

	userID := model.GetUserID(c)
	if err := model.DeleteRespondent(userID, responseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
