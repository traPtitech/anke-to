package controller

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
	"gopkg.in/guregu/null.v4"
)

func questionnaireInfo2questionnaireSummary(questionnaireInfo model.QuestionnaireInfo, allResponded bool, hasMyDraft bool, hasMyResponse bool, respondedDateTimeByMe null.Time) *openapi.QuestionnaireSummary {
	res := openapi.QuestionnaireSummary{
		AllResponded:             allResponded,
		CreatedAt:                questionnaireInfo.CreatedAt,
		Description:              questionnaireInfo.Description,
		HasMyDraft:               hasMyDraft,
		HasMyResponse:            hasMyResponse,
		IsDuplicateAnswerAllowed: questionnaireInfo.IsDuplicateAnswerAllowed,
		IsAnonymous:              questionnaireInfo.IsAnonymous,
		IsPublished:              questionnaireInfo.IsPublished,
		IsTargetingMe:            questionnaireInfo.IsTargeted,
		ModifiedAt:               questionnaireInfo.ModifiedAt,
		QuestionnaireId:          questionnaireInfo.ID,
		Title:                    questionnaireInfo.Title,
	}
	if respondedDateTimeByMe.Valid {
		res.RespondedDateTimeByMe = &respondedDateTimeByMe.Time
	} else {
		res.RespondedDateTimeByMe = nil
	}
	if questionnaireInfo.ResTimeLimit.Valid {
		res.ResponseDueDateTime = &questionnaireInfo.ResTimeLimit.Time
	} else {
		res.ResponseDueDateTime = nil
	}
	return &res
}

func convertResponseViewableBy(resShareType openapi.ResShareType) string {
	switch resShareType {
	case "admins":
		return "administrators"
	case "respondents":
		return "respondents"
	case "anyone":
		return "public"
	default:
		return "administrators"
	}
}

func convertResSharedTo(resSharedTo string) openapi.ResShareType {
	switch resSharedTo {
	case "administrators":
		return "admins"
	case "respondents":
		return "respondents"
	case "public":
		return "anyone"
	default:
		return "admins"
	}

}

func createUsersAndGroups(users []string, groups uuid.UUIDs) openapi.UsersAndGroups {
	res := openapi.UsersAndGroups{
		Users:  users,
		Groups: groups,
	}
	return res
}

func convertOptions(options []model.Options) openapi.QuestionSettingsSingleChoice {
	res := openapi.QuestionSettingsSingleChoice{}
	for _, option := range options {
		res.Options = append(res.Options, option.Body)
	}
	return res
}

func convertQuestions(questions []model.Questions) ([]openapi.Question, error) {
	res := []openapi.Question{}
	for _, question := range questions {
		q := openapi.Question{
			CreatedAt:   &question.CreatedAt,
			Title:       question.Body,
			Description: question.Description,
			IsRequired:  question.IsRequired,
			QuestionId:  &question.ID,
		}
		switch question.Type {
		case "Text":
			err := q.FromQuestionSettingsText(
				openapi.QuestionSettingsText{
					QuestionType: "Text",
				},
			)
			if err != nil {
				return nil, err
			}
		case "TextArea":
			err := q.FromQuestionSettingsText(
				openapi.QuestionSettingsText{
					QuestionType: "TextLong",
				},
			)
			if err != nil {
				return nil, err
			}
		case "Number":
			err := q.FromQuestionSettingsNumber(
				openapi.QuestionSettingsNumber{
					QuestionType: "Number",
				},
			)
			if err != nil {
				return nil, err
			}
		case "MultipleChoice":
			var err error
			question.Options, err = model.NewOption().GetOptions(context.Background(), []int{question.ID})
			if err != nil {
				return nil, err
			}
			err = q.FromQuestionSettingsSingleChoice(
				openapi.QuestionSettingsSingleChoice{
					QuestionType: "SingleChoice",
					Options:      convertOptions(question.Options).Options,
				},
			)
			if err != nil {
				return nil, err
			}
		case "Checkbox":
			var err error
			question.Options, err = model.NewOption().GetOptions(context.Background(), []int{question.ID})
			if err != nil {
				return nil, err
			}
			err = q.FromQuestionSettingsMultipleChoice(
				openapi.QuestionSettingsMultipleChoice{
					QuestionType: "MultipleChoice",
					Options:      convertOptions(question.Options).Options,
				},
			)
			if err != nil {
				return nil, err
			}
		case "LinearScale":
			var err error
			question.ScaleLabels, err = model.NewScaleLabel().GetScaleLabels(context.Background(), []int{question.ID})
			if err != nil {
				return nil, err
			}
			err = q.FromQuestionSettingsScale(
				openapi.QuestionSettingsScale{
					QuestionType: "Scale",
					MinLabel:     &question.ScaleLabels[0].ScaleLabelLeft,
					MaxLabel:     &question.ScaleLabels[0].ScaleLabelRight,
					MinValue:     question.ScaleLabels[0].ScaleMin,
					MaxValue:     question.ScaleLabels[0].ScaleMax,
				},
			)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, q)
	}
	return res, nil
}

func questionnaire2QuestionnaireDetail(questionnaires model.Questionnaires, admins []string, adminUsers []string, adminGroups []uuid.UUID, targets []string, targetUsers []string, targetGroups []uuid.UUID, respondents []string) (openapi.QuestionnaireDetail, error) {
	questions, err := model.NewQuestion().GetQuestions(context.Background(), questionnaires.ID)
	if err != nil {
		return openapi.QuestionnaireDetail{}, err
	}
	questionsConverted, err := convertQuestions(questions)
	if err != nil {
		return openapi.QuestionnaireDetail{}, err
	}
	responseDueDateTime := &questionnaires.ResTimeLimit.Time
	if !questionnaires.ResTimeLimit.Valid {
		responseDueDateTime = nil
	}
	res := openapi.QuestionnaireDetail{
		Admin:                    createUsersAndGroups(adminUsers, adminGroups),
		Admins:                   admins,
		CreatedAt:                questionnaires.CreatedAt,
		Description:              questionnaires.Description,
		IsDuplicateAnswerAllowed: questionnaires.IsDuplicateAnswerAllowed,
		IsAnonymous:              questionnaires.IsAnonymous,
		IsPublished:              questionnaires.IsPublished,
		ModifiedAt:               questionnaires.ModifiedAt,
		QuestionnaireId:          questionnaires.ID,
		Questions:                questionsConverted,
		Respondents:              respondents,
		ResponseDueDateTime:      responseDueDateTime,
		ResponseViewableBy:       convertResSharedTo(questionnaires.ResSharedTo),
		Target:                   createUsersAndGroups(targetUsers, targetGroups),
		Targets:                  targets,
		Title:                    questionnaires.Title,
	}
	return res, nil
}

func respondentDetail2Response(ctx echo.Context, respondentDetail model.RespondentDetail) (openapi.Response, error) {
	oResponseBodies := []openapi.ResponseBody{}
	for _, r := range respondentDetail.Responses {
		oResponseBody := openapi.ResponseBody{}
		oResponseBody.QuestionId = r.QuestionID
		switch r.QuestionType {
		case "Text":
			if r.Body.Valid {
				err := oResponseBody.FromResponseBodyText(
					openapi.ResponseBodyText{
						Answer:       r.Body.String,
						QuestionType: "Text",
					},
				)
				if err != nil {
					return openapi.Response{}, err
				}
			}
		case "TextArea":
			if r.Body.Valid {
				err := oResponseBody.FromResponseBodyTextLong(
					openapi.ResponseBodyTextLong{
						Answer:       r.Body.String,
						QuestionType: "TextLong",
					},
				)
				if err != nil {
					return openapi.Response{}, err
				}
			}
		case "Number":
			if r.Body.Valid {
				answer, err := strconv.ParseFloat(r.Body.String, 32)
				if err != nil {
					ctx.Logger().Errorf("failed to convert string to float: %+v", err)
					return openapi.Response{}, err
				}
				err = oResponseBody.FromResponseBodyNumber(
					openapi.ResponseBodyNumber{
						Answer:       float32(answer),
						QuestionType: "Number",
					},
				)
				if err != nil {
					return openapi.Response{}, err
				}
			}
		case "Checkbox":
			if len(r.OptionResponse) > 0 {
				if len(r.OptionResponse) > 1 {
					return openapi.Response{}, errors.New("too many responses")
				}
				answer := []int{}
				for _, o := range r.OptionResponse {
					err := json.Unmarshal([]byte(o), &answer)
					if err != nil {
						return openapi.Response{}, err
					}
					err = oResponseBody.FromResponseBodyMultipleChoice(
						openapi.ResponseBodyMultipleChoice{
							Answer:       answer,
							QuestionType: "MultipleChoice",
						},
					)
					if err != nil {
						return openapi.Response{}, err
					}
				}
			}
		case "MultipleChoice":
			if len(r.OptionResponse) > 0 {
				if len(r.OptionResponse) > 1 {
					return openapi.Response{}, errors.New("too many responses")
				}
				for _, o := range r.OptionResponse {
					var option int
					err := json.Unmarshal([]byte(o), &option)
					if err != nil {
						return openapi.Response{}, err
					}
					err = oResponseBody.FromResponseBodySingleChoice(
						openapi.ResponseBodySingleChoice{
							Answer:       option,
							QuestionType: "SingleChoice",
						},
					)
					if err != nil {
						return openapi.Response{}, err
					}
				}
			}
		case "LinearScale":
			if r.Body.Valid {
				answer, err := strconv.Atoi(r.Body.String)
				if err != nil {
					ctx.Logger().Errorf("failed to convert string to int: %+v", err)
					return openapi.Response{}, err
				}
				err = oResponseBody.FromResponseBodyScale(
					openapi.ResponseBodyScale{
						Answer:       answer,
						QuestionType: "Scale",
					},
				)
				if err != nil {
					return openapi.Response{}, err
				}
			}
		}
		oResponseBodies = append(oResponseBodies, oResponseBody)
	}

	isAnonymous, err := model.NewQuestionnaire().GetResponseIsAnonymousByQuestionnaireID(ctx.Request().Context(), respondentDetail.QuestionnaireID)
	if err != nil {
		ctx.Logger().Errorf("failed to get response is anonymous: %+v", err)
		return openapi.Response{}, err
	}

	respondent := &respondentDetail.TraqID
	if isAnonymous {
		respondent = nil
	}

	res := openapi.Response{
		Body:            oResponseBodies,
		IsAnonymous:     &isAnonymous,
		IsDraft:         !respondentDetail.SubmittedAt.Valid,
		ModifiedAt:      respondentDetail.ModifiedAt,
		QuestionnaireId: respondentDetail.QuestionnaireID,
		Respondent:      respondent,
		ResponseId:      respondentDetail.ResponseID,
		SubmittedAt:     respondentDetail.SubmittedAt.Time,
	}

	return res, nil
}

func responseBody2ResponseMetas(body []openapi.ResponseBody, questions []model.Questions) ([]*model.ResponseMeta, error) {
	res := []*model.ResponseMeta{}

	var questionIDMap = make(map[int]int, len(questions))
	for i, question := range questions {
		questionIDMap[question.ID] = i
	}
	for _, b := range body {
		i, ok := questionIDMap[b.QuestionId]
		if !ok {
			return nil, errors.New("question not found")
		}
		switch questions[i].Type {
		case "Text":
			bText, err := b.AsResponseBodyText()
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       bText.Answer,
			})
		case "TextArea":
			bTextLong, err := b.AsResponseBodyTextLong()
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       bTextLong.Answer,
			})
		case "Number":
			bNumber, err := b.AsResponseBodyNumber()
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       strconv.FormatFloat(float64(bNumber.Answer), 'f', -1, 32),
			})
		case "MultipleChoice":
			bSingleChoice, err := b.AsResponseBodySingleChoice()
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(bSingleChoice.Answer)
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       string(data),
			})
		case "Checkbox":
			bMultipleChoice, err := b.AsResponseBodyMultipleChoice()
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(bMultipleChoice.Answer)
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       string(data),
			})
		case "LinearScale":
			bScale, err := b.AsResponseBodyScale()
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       strconv.FormatInt(int64(bScale.Answer), 10),
			})
		}
	}
	return res, nil
}
