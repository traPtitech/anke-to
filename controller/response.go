package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
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
	model.IOption
	model.IValidation
	model.IScaleLabel
	model.ITransaction
}

func NewResponse(
	questionnaire model.IQuestionnaire,
	respondent model.IRespondent,
	response model.IResponse,
	target model.ITarget,
	question model.IQuestion,
	option model.IOption,
	validation model.IValidation,
	scaleLabel model.IScaleLabel,
	transaction model.ITransaction,
) *Response {
	return &Response{
		IQuestionnaire: questionnaire,
		IRespondent:    respondent,
		IResponse:      response,
		ITarget:        target,
		IQuestion:      question,
		IOption:        option,
		IValidation:    validation,
		IScaleLabel:    scaleLabel,
		ITransaction:   transaction,
	}
}

func (r *Response) GetMyResponses(ctx echo.Context, params openapi.GetMyResponsesParams, userID string) (openapi.ResponsesWithQuestionnaireInfo, error) {
	res := openapi.ResponsesWithQuestionnaireInfo{}

	var sort string
	if params.Sort == nil {
		sort = ""
	} else {
		sort = string(*params.Sort)
	}
	var questionnaireIDs []int
	if params.QuestionnaireIDs == nil {
		questionnaireIDs = nil
	} else {
		questionnaireIDs = *params.QuestionnaireIDs
	}
	responsesID, err := r.IRespondent.GetMyResponseIDs(ctx.Request().Context(), sort, userID, questionnaireIDs, params.IsDraft)
	if err != nil {
		ctx.Logger().Errorf("failed to get my responses ID: %+v", err)
		return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire responses: %w", err))
	}

	responseLists := make(map[int][]openapi.Response)
	var responseQuestionnaireIDs []int

	for _, responseID := range responsesID {
		responseDetail, err := r.IRespondent.GetRespondentDetail(ctx.Request().Context(), responseID)
		if err != nil {
			ctx.Logger().Errorf("failed to get respondent detail: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondent detail: %w", err))
		}

		response, err := respondentDetail2Response(ctx, responseDetail)
		if err != nil {
			ctx.Logger().Errorf("failed to convert respondent detail into response: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to convert respondent detail into response: %w", err))
		}

		tmp := openapi.Response{
			Body:            response.Body,
			IsDraft:         response.IsDraft,
			ModifiedAt:      response.ModifiedAt,
			QuestionnaireId: response.QuestionnaireId,
			Respondent:      &userID,
			ResponseId:      response.ResponseId,
			SubmittedAt:     response.SubmittedAt,
			IsAnonymous:     response.IsAnonymous,
		}
		responseLists[responseDetail.QuestionnaireID] = append(responseLists[responseDetail.QuestionnaireID], tmp)
		responseQuestionnaireIDs = append(responseQuestionnaireIDs, responseDetail.QuestionnaireID)
	}

	slices.Sort(responseQuestionnaireIDs)
	for i, questionnaireID := range responseQuestionnaireIDs {
		if i != 0 && responseQuestionnaireIDs[i-1] == questionnaireID {
			continue
		}

		questionnaire, _, _, _, _, _, _, _, err := r.IQuestionnaire.GetQuestionnaireInfo(ctx.Request().Context(), questionnaireID)
		if err != nil {
			ctx.Logger().Errorf("failed to get questionnaire info: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire info: %w", err))
		}

		isTargetingMe, err := r.ITarget.IsTargetingMe(ctx.Request().Context(), questionnaireID, userID)
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

		responses := responseLists[questionnaireID]
		res = append(res, openapi.ResponseWithQuestionnaireInfoItem{
			QuestionnaireInfo: &questionnaireInfo,
			Responses:         &responses,
		})
	}

	return res, nil
}

func (r *Response) GetResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) (openapi.Response, error) {
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

func (r *Response) DeleteResponse(ctx echo.Context, responseID openapi.ResponseIDInPath) error {
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

func (r *Response) EditResponse(ctx echo.Context, responseID openapi.ResponseIDInPath, req openapi.EditResponseJSONRequestBody) error {
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

	respondentDetail, err := r.IRespondent.GetRespondentDetail(ctx.Request().Context(), responseID)
	if err != nil {
		ctx.Logger().Errorf("failed to get respondent detail: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondent detail: %w", err))
	}

	questions, err := r.IQuestion.GetQuestions(ctx.Request().Context(), respondentDetail.QuestionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to get questions: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questions: %w", err))
	}

	responseMetas, err := responseBody2ResponseMetas(req.Body, questions)
	if err != nil {
		ctx.Logger().Errorf("failed to convert response body into response metas: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to convert response body into response metas: %w", err))
	}

	// validationでチェック
	questionIDs := make([]int, len(questions))
	questionTypes := make(map[int]string, len(questions))
	questionRequired := make(map[int]bool, len(questions))
	for i, question := range questions {
		questionIDs[i] = question.ID
		questionTypes[question.ID] = question.Type
		questionRequired[question.ID] = question.IsRequired
	}

	validations, err := r.IValidation.GetValidations(ctx.Request().Context(), questionIDs)
	if err != nil {
		ctx.Logger().Errorf("failed to get validations: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	validationMap := make(map[int]model.Validations, len(validations))
	for _, validation := range validations {
		validationMap[validation.QuestionID] = validation
	}

	options, err := r.IOption.GetOptions(ctx.Request().Context(), questionIDs)
	if err != nil {
		ctx.Logger().Errorf("failed to get options: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	optionMap := make(map[int][]model.Options, len(options))
	for _, option := range options {
		optionMap[option.QuestionID] = append(optionMap[option.QuestionID], option)
	}

	scaleLabels, err := r.IScaleLabel.GetScaleLabels(ctx.Request().Context(), questionIDs)
	if err != nil {
		ctx.Logger().Errorf("failed to get scale labels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, scaleLabel := range scaleLabels {
		scaleLabelMap[scaleLabel.QuestionID] = scaleLabel
	}

	for _, responseMeta := range responseMetas {
		questionRequired[responseMeta.QuestionID] = false
		switch questionTypes[responseMeta.QuestionID] {
		case "Text", "TextArea":
			validation, ok := validationMap[responseMeta.QuestionID]
			if !ok {
				validation = model.Validations{}
			}
			err := r.IValidation.CheckTextValidation(validation, responseMeta.Data)
			if err != nil {
				if errors.Is(err, model.ErrTextMatching) {
					ctx.Logger().Errorf("invalid text: %+v", err)
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				ctx.Logger().Errorf("invalid text: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Number":
			validation, ok := validationMap[responseMeta.QuestionID]
			if !ok {
				validation = model.Validations{}
			}
			err := r.IValidation.CheckNumberValidation(validation, responseMeta.Data)
			if err != nil {
				if errors.Is(err, model.ErrInvalidNumber) {
					ctx.Logger().Errorf("invalid number: %+v", err)
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				ctx.Logger().Errorf("invalid number: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Checkbox", "MultipleChoice":
		case "LinearScale":
			label, ok := scaleLabelMap[responseMeta.QuestionID]
			if !ok {
				label = model.ScaleLabels{}
			}
			err := r.IScaleLabel.CheckScaleLabel(label, responseMeta.Data)
			if err != nil {
				ctx.Logger().Errorf("invalid scale: %+v", err)
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		default:
			ctx.Logger().Errorf("invalid question id: %+v", responseMeta.QuestionID)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid question id: %d", responseMeta.QuestionID))
		}
	}

	for _, question := range questions {
		if questionRequired[question.ID] {
			ctx.Logger().Errorf("required question is not answered: %+v", question.ID)
			return echo.NewHTTPError(http.StatusBadRequest, "required question is not answered")
		}
	}

	err = r.ITransaction.Do(ctx.Request().Context(), nil, func(c context.Context) error {
		if !respondentDetail.SubmittedAt.Valid {
			if !req.IsDraft {
				err := r.IRespondent.UpdateSubmittedAt(c, responseID)
				if err != nil {
					ctx.Logger().Errorf("failed to update submitted at: %+v", err)
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to update submitted at: %w", err))
				}
			}
		} else {
			if req.IsDraft {
				ctx.Logger().Errorf("unable to update the response to draft")
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to update the response to draft"))
			}
		}
		err := r.IRespondent.UpdateModifiedAt(c, responseID)
		if err != nil {
			ctx.Logger().Errorf("failed to update modified at: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to update modified at: %w", err))
		}

		if len(responseMetas) > 0 {
			err = r.IResponse.InsertResponses(c, responseID, responseMetas)
			if err != nil {
				ctx.Logger().Errorf("failed to insert responses: %+v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert responses: %w", err))
			}
		}

		return nil
	})
	if err != nil {
		ctx.Logger().Errorf("failed to update response: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to update response: %w", err))
	}

	return nil
}
