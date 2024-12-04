package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
)

// Response Responseの構造体
type Response struct {
	model.IQuestionnaire
	model.IRespondent
	model.IResponse
	model.ITarget
	model.IQuestion
	model.IValidation
	model.IScaleLabel
}

func NewResponse() *Response {
	return &Response{}
}

func (r Response) GetMyResponses(ctx echo.Context, params openapi.GetMyResponsesParams, userID string) (openapi.ResponsesWithQuestionnaireInfo, error) {
	res := openapi.ResponsesWithQuestionnaireInfo{}

	sort := string(*params.Sort)
	responsesID := []int{}
	responsesID, err := r.IRespondent.GetMyResponseIDs(ctx.Request().Context(), sort, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get my responses ID: %+v", err)
		return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire responses: %w", err))
	}

	for _, responseID := range responsesID {
		responseDetail, err := r.IRespondent.GetRespondentDetail(ctx.Request().Context(), responseID)
		if err != nil {
			ctx.Logger().Errorf("failed to get respondent detail: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondent detail: %w", err))
		}

		questionnaire, _, _, _, _, _, err := r.IQuestionnaire.GetQuestionnaireInfo(ctx.Request().Context(), responseDetail.QuestionnaireID)
		if err != nil {
			ctx.Logger().Errorf("failed to get questionnaire info: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire info: %w", err))
		}

		isTargetingMe, err := r.ITarget.IsTargetingMe(ctx.Request().Context(), responseDetail.QuestionnaireID, userID)
		if err != nil {
			ctx.Logger().Errorf("failed to get target info: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get target info: %w", err))
		}

		questionnaireInfo := openapi.QuestionnaireInfo{
			CreatedAt:           questionnaire.CreatedAt,
			IsTargetingMe:       isTargetingMe,
			ModifiedAt:          questionnaire.ModifiedAt,
			ResponseDueDateTime: &questionnaire.ResTimeLimit.Time,
			Title:               questionnaire.Title,
		}

		response, err := respondentDetail2Response(ctx, responseDetail)
		if err != nil {
			ctx.Logger().Errorf("failed to convert respondent detail into response: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to convert respondent detail into response: %w", err))
		}

		tmp := openapi.ResponseWithQuestionnaireInfoItem{
			Body:              response.Body,
			IsDraft:           response.IsDraft,
			ModifiedAt:        response.ModifiedAt,
			QuestionnaireId:   response.QuestionnaireId,
			QuestionnaireInfo: &questionnaireInfo,
			Respondent:        &userID,
			ResponseId:        response.ResponseId,
			SubmittedAt:       response.SubmittedAt,
			IsAnonymous:       response.IsAnonymous,
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (r Response) GetResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) (openapi.Response, error) {
	responseDetail, err := r.IRespondent.GetRespondentDetail(ctx.Request().Context(), responseID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			ctx.Logger().Errorf("failed to find response by response ID: %+v", err)
			return openapi.Response{}, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("failed to find response by response ID: %w", err))
		}
		ctx.Logger().Errorf("failed to get respondent detail: %+v", err)
		return openapi.Response{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondent detail: %w", err))
	}

	res, err := respondentDetail2Response(ctx, responseDetail)
	if err != nil {
		ctx.Logger().Errorf("failed to convert respondent detail into response: %+v", err)
		return openapi.Response{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to convert respondent detail into response: %w", err))
	}

	return res, nil
}

func (r Response) DeleteResponse(ctx echo.Context, responseID openapi.ResponseIDInPath, userID string) error {
	limit, err := r.IQuestionnaire.GetQuestionnaireLimitByResponseID(ctx.Request().Context(), responseID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			ctx.Logger().Errorf("failed to find response by response ID: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("failed to find response by response ID: %w", err))
		}
		ctx.Logger().Errorf("failed to get questionnaire limit by response ID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire limit by response ID: %w", err))
	}
	if limit.Valid && limit.Time.Before(time.Now()) {
		ctx.Logger().Errorf("unable delete the expired response")
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Errorf("unable delete the expired response"))
	}

	err = r.IRespondent.DeleteRespondent(ctx.Request().Context(), responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to delete respondent: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete respondent: %w", err))
	}

	err = r.IResponse.DeleteResponse(ctx.Request().Context(), responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to delete response: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete response: %w", err))
	}

	return nil
}

func (r Response) EditResponse(ctx echo.Context, responseID openapi.ResponseIDInPath, req openapi.EditResponseJSONRequestBody) error {
	limit, err := r.IQuestionnaire.GetQuestionnaireLimitByResponseID(ctx.Request().Context(), responseID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			ctx.Logger().Infof("failed to find response by response ID: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("failed to find response by response ID: %w", err))
		}
		ctx.Logger().Errorf("failed to get questionnaire limit by response ID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire limit by response ID: %w", err))
	}

	if limit.Valid && limit.Time.Before(time.Now()) {
		ctx.Logger().Info("unable to edit the expired response")
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Errorf("unable edit the expired response"))
	}

	err = r.IResponse.DeleteResponse(ctx.Request().Context(), responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to delete response: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete response: %w", err))
	}

	questions, err := r.IQuestion.GetQuestions(ctx.Request().Context(), req.QuestionnaireId)

	responseMetas, err := responseBody2ResponseMetas(req.Body, questions)
	if err != nil {
		ctx.Logger().Errorf("failed to convert response body into response metas: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to convert response body into response metas: %w", err))
	}

	// validationでチェック
	questionIDs := make([]int, len(questions))
	questionTypes := make(map[int]string, len(questions))
	for i, question := range questions {
		questionIDs[i] = question.ID
		questionTypes[question.ID] = question.Type
	}

	validations, err := r.IValidation.GetValidations(ctx.Request().Context(), questionIDs)
	if err != nil {
		ctx.Logger().Errorf("failed to get validations: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get validations: %w", err))
	}

	for i, validation := range validations {
		switch questionTypes[validation.QuestionID] {
		case "Text", "TextLong":
			err := r.IValidation.CheckTextValidation(validation, responseMetas[i].Data)
			if err != nil {
				if errors.Is(err, model.ErrTextMatching) {
					ctx.Logger().Errorf("invalid text: %+v", err)
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid text: %w", err))
				}
				ctx.Logger().Errorf("invalid text: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid text: %w", err))
			}
		case "Number":
			err := r.IValidation.CheckNumberValidation(validation, responseMetas[i].Data)
			if err != nil {
				if errors.Is(err, model.ErrInvalidNumber) {
					ctx.Logger().Errorf("invalid number: %+v", err)
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid number: %w", err))
				}
				ctx.Logger().Errorf("invalid number: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid number: %w", err))
			}
		}
	}

	// scaleのvalidation
	scaleLabelIDs := []int{}
	for _, question := range questions {
		if question.Type == "Scale" {
			scaleLabelIDs = append(scaleLabelIDs, question.ID)
		}
	}

	scaleLabels, err := r.IScaleLabel.GetScaleLabels(ctx.Request().Context(), scaleLabelIDs)
	if err != nil {
		ctx.Logger().Errorf("failed to get scale labels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get scale labels: %w", err))
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, scaleLabel := range scaleLabels {
		scaleLabelMap[scaleLabel.QuestionID] = scaleLabel
	}

	for i, question := range questions {
		if question.Type == "Scale" {
			label, ok := scaleLabelMap[question.ID]
			if !ok {
				label = model.ScaleLabels{}
			}
			err := r.IScaleLabel.CheckScaleLabel(label, responseMetas[i].Data)
			if err != nil {
				ctx.Logger().Errorf("invalid scale: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid scale: %w", err))
			}
		}
	}

	if len(responseMetas) > 0 {
		err = r.IResponse.InsertResponses(ctx.Request().Context(), responseID, responseMetas)
		if err != nil {
			ctx.Logger().Errorf("failed to insert responses: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert responses: %w", err))
		}
	}

	return nil
}
