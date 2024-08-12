package controller

import (
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
	model.ITarget
}

func NewResponse() *Response {
	return &Response{}
}

func (r Response) GetMyResponses(ctx echo.Context, params openapi.GetMyResponsesParams, userID string) (openapi.ResponsesWithQuestionnaireInfo, error) {
	res := openapi.ResponsesWithQuestionnaireInfo{}

	sort := string(*params.Sort)
	responsesID := []int{}
	responsesID, err := r.IRespondent.GetMyResponsesID(ctx.Request().Context(), sort, userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get my responses ID: %+v", err)
		return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire responses: %+w", err))
	}

	for _, responseID := range responsesID {
		responseDetail, err := r.IRespondent.GetRespondentDetail(ctx.Request().Context(), responseID)
		if err != nil {
			ctx.Logger().Errorf("failed to get respondent detail: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondent detail: %+w", err))
		}

		questionnaire, _, _, _, _, _, err := r.IQuestionnaire.GetQuestionnaireInfo(ctx.Request().Context(), responseDetail.QuestionnaireID)
		if err != nil {
			ctx.Logger().Errorf("failed to get questionnaire info: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire info: %+w", err))
		}

		isTargetingMe, err := r.ITarget.IsTargetingMe(ctx.Request().Context(), responseDetail.QuestionnaireID, userID)
		if err != nil {
			ctx.Logger().Errorf("failed to get target info: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get target info: %+w", err))
		}

		questionnaireInfo := struct {
			CreatedAt           time.Time  `json:"created_at"`
			IsTargetingMe       bool       `json:"is_targeting_me"`
			ModifiedAt          time.Time  `json:"modified_at"`
			ResponseDueDateTime *time.Time `json:"response_due_date_time,omitempty"`
			Title               string     `json:"title"`
		}{
			CreatedAt:           questionnaire.CreatedAt,
			IsTargetingMe:       isTargetingMe,
			ModifiedAt:          questionnaire.ModifiedAt,
			ResponseDueDateTime: &questionnaire.ResTimeLimit.Time,
			Title:               questionnaire.Title,
		}

		response, err := respondentDetail2Response(ctx, responseDetail)
		if err != nil {
			ctx.Logger().Errorf("failed to convert respondent detail into response: %+v", err)
			return openapi.ResponsesWithQuestionnaireInfo{}, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to convert respondent detail into response: %+w", err))
		}

		tmp := struct {
			Body              []openapi.ResponseBody `json:"body"`
			IsDraft           bool                   `json:"is_draft"`
			ModifiedAt        time.Time              `json:"modified_at"`
			QuestionnaireId   int                    `json:"questionnaire_id"`
			QuestionnaireInfo *struct {
				CreatedAt           time.Time  `json:"created_at"`
				IsTargetingMe       bool       `json:"is_targeting_me"`
				ModifiedAt          time.Time  `json:"modified_at"`
				ResponseDueDateTime *time.Time `json:"response_due_date_time,omitempty"`
				Title               string     `json:"title"`
			} `json:"questionnaire_info,omitempty"`
			Respondent  openapi.TraqId `json:"respondent"`
			ResponseId  int            `json:"response_id"`
			SubmittedAt time.Time      `json:"submitted_at"`
		}{
			Body:              response.Body,
			IsDraft:           response.IsDraft,
			ModifiedAt:        response.ModifiedAt,
			QuestionnaireId:   response.QuestionnaireId,
			QuestionnaireInfo: &questionnaireInfo,
			Respondent:        userID,
			ResponseId:        response.ResponseId,
			SubmittedAt:       response.SubmittedAt,
		}
		res = append(res, tmp)
	}

	return res, nil
}
