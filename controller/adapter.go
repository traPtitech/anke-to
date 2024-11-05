package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
	"gopkg.in/guregu/null.v4"
)

func questionnaireInfo2questionnaireSummary(questionnaireInfo model.QuestionnaireInfo, allResponded bool, hasMyDraft bool, hasMyResponse bool, respondedDateTimeByMe null.Time) *openapi.QuestionnaireSummary {
	res := openapi.QuestionnaireSummary{
		AllResponded:  allResponded,
		CreatedAt:     questionnaireInfo.CreatedAt,
		Description:   questionnaireInfo.Description,
		HasMyDraft:    hasMyDraft,
		HasMyResponse: hasMyResponse,
		// IsAllowingMultipleResponses: questionnaireInfo.IsAllowingMultipleResponses,
		// IsAnonymous:                 questionnaireInfo.IsAnonymous,
		// IsPublished:                 questionnaireInfo.IsPublished,
		IsTargetingMe:   questionnaireInfo.IsTargeted,
		ModifiedAt:      questionnaireInfo.ModifiedAt,
		QuestionnaireId: questionnaireInfo.ID,
		Title:           questionnaireInfo.Title,
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

func createUsersAndGroups(users []string, groups []string) openapi.UsersAndGroups {
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

func convertQuestions(questions []model.Questions) []openapi.Question {
	res := []openapi.Question{}
	for _, question := range questions {
		q := openapi.Question{
			CreatedAt: question.CreatedAt,
			// Description:     question.Description,
			IsRequired:      question.IsRequired,
			QuestionId:      question.ID,
			QuestionnaireId: question.QuestionnaireID,
			Title:           question.Body,
		}
		switch question.Type {
		case "Text":
			q.FromQuestionSettingsText(
				openapi.QuestionSettingsText{
					QuestionType: "Text",
				},
			)
		case "TextArea":
			q.FromQuestionSettingsText(
				openapi.QuestionSettingsText{
					QuestionType: "TextLong",
				},
			)
		case "Number":
			q.FromQuestionSettingsNumber(
				openapi.QuestionSettingsNumber{
					QuestionType: "Number",
				},
			)
		case "Radio":
			q.FromQuestionSettingsSingleChoice(
				openapi.QuestionSettingsSingleChoice{
					QuestionType: "Radio",
					Options:      convertOptions(question.Options).Options,
				},
			)
		case "MultipleChoice":
			q.FromQuestionSettingsMultipleChoice(
				openapi.QuestionSettingsMultipleChoice{
					QuestionType: "MultipleChoice",
					Options:      convertOptions(question.Options).Options,
				},
			)
		case "LinearScale":
			q.FromQuestionSettingsScale(
				openapi.QuestionSettingsScale{
					QuestionType: "LinearScale",
					MinLabel:     &question.ScaleLabels[0].ScaleLabelLeft,
					MaxLabel:     &question.ScaleLabels[0].ScaleLabelRight,
					MinValue:     question.ScaleLabels[0].ScaleMin,
					MaxValue:     question.ScaleLabels[0].ScaleMax,
				},
			)
		}
	}
	return res
}

func convertRespondents(respondents []model.Respondents) []string {
	res := []string{}
	for _, respondent := range respondents {
		res = append(res, respondent.UserTraqid)
	}
	return res
}

func questionnaire2QuestionnaireDetail(questionnaires model.Questionnaires, adminUsers []string, adminGroups []string, targetUsers []string, targetGroups []string, respondents []string) openapi.QuestionnaireDetail {
	res := openapi.QuestionnaireDetail{
		Admins:      createUsersAndGroups(adminUsers, adminGroups),
		CreatedAt:   questionnaires.CreatedAt,
		Description: questionnaires.Description,
		IsAllowingMultipleResponses: questionnaires.IsAllowingMultipleResponses,
		IsAnonymous:                 questionnaires.IsAnonymous,
		IsPublished:                 questionnaires.IsPublished,
		ModifiedAt:          questionnaires.ModifiedAt,
		QuestionnaireId:     questionnaires.ID,
		Questions:           convertQuestions(questionnaires.Questions),
		Respondents:         respondents,
		ResponseDueDateTime: &questionnaires.ResTimeLimit.Time,
		ResponseViewableBy:  convertResSharedTo(questionnaires.ResSharedTo),
		Targets:             createUsersAndGroups(targetUsers, targetGroups),
		Title:               questionnaires.Title,
	}
	return res
}

func respondentDetail2Response(ctx echo.Context, respondentDetail model.RespondentDetail) (openapi.Response, error) {
	oResponseBodies := []openapi.ResponseBody{}
	for j, r := range respondentDetail.Responses {
		oResponseBody := openapi.ResponseBody{}
		switch r.QuestionType {
		case "Text":
			if r.Body.Valid {
				oResponseBody.FromResponseBodyText(
					openapi.ResponseBodyText{
						Answer:       r.Body.String,
						QuestionType: "Text",
					},
				)
			}
		case "TextArea":
			if r.Body.Valid {
				oResponseBody.FromResponseBodyText(
					openapi.ResponseBodyText{
						Answer:       r.Body.String,
						QuestionType: "TextLong",
					},
				)
			}
		case "Number":
			if r.Body.Valid {
				answer, err := strconv.ParseFloat(r.Body.String, 32)
				if err != nil {
					ctx.Logger().Errorf("failed to convert string to float: %+v", err)
					return openapi.Response{}, err
				}
				oResponseBody.FromResponseBodyNumber(
					openapi.ResponseBodyNumber{
						Answer:       float32(answer),
						QuestionType: "Number",
					},
				)
			}
		case "MultipleChoice":
			if r.Body.Valid {
				answer := []int{}
				questionnaire, _, _, _, _, _, err := model.NewQuestionnaire().GetQuestionnaireInfo(ctx.Request().Context(), r.QuestionID)
				if err != nil {
					ctx.Logger().Errorf("failed to get questionnaire info: %+v", err)
					return openapi.Response{}, err
				}
				for _, a := range r.OptionResponse {
					for i, o := range questionnaire.Questions[j].Options {
						if a == o.Body {
							answer = append(answer, i)
						}
					}
				}
				oResponseBody.FromResponseBodyMultipleChoice(
					openapi.ResponseBodyMultipleChoice{
						Answer:       answer,
						QuestionType: "MultipleChoice",
					},
				)
			}
		case "Checkbox":
			if r.Body.Valid {
				questionnaire, _, _, _, _, _, err := model.NewQuestionnaire().GetQuestionnaireInfo(ctx.Request().Context(), r.QuestionID)
				if err != nil {
					ctx.Logger().Errorf("failed to get questionnaire info: %+v", err)
					return openapi.Response{}, err
				}
				for _, a := range r.OptionResponse {
					for i, o := range questionnaire.Questions[j].Options {
						if a == o.Body {
							oResponseBody.FromResponseBodySingleChoice(
								openapi.ResponseBodySingleChoice{
									Answer:       i,
									QuestionType: "SingleChoice",
								},
							)
						}
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
				oResponseBody.FromResponseBodyScale(
					openapi.ResponseBodyScale{
						Answer:       answer,
						QuestionType: "LinearScale",
					},
				)
			}
		}
		oResponseBodies = append(oResponseBodies, oResponseBody)
	}

	res := openapi.Response{
		Body:            oResponseBodies,
		IsDraft:         respondentDetail.SubmittedAt.Valid,
		ModifiedAt:      respondentDetail.ModifiedAt,
		QuestionnaireId: respondentDetail.QuestionnaireID,
		Respondent:      respondentDetail.TraqID,
		ResponseId:      respondentDetail.ResponseID,
		SubmittedAt:     respondentDetail.SubmittedAt.Time,
	}

	return res, nil
}

func responseBody2ResponseMetas(body []openapi.ResponseBody, questions []model.Questions) ([]*model.ResponseMeta, error) {
	res := []*model.ResponseMeta{}

	for i, b := range body {
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
		case "TextLong":
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
		case "SingleChoice":
			bSingleChoice, err := b.AsResponseBodySingleChoice()
			if err != nil {
				return nil, err
			}
			res = append(res, &model.ResponseMeta{
				QuestionID: questions[i].ID,
				Data:       strconv.FormatInt(int64(bSingleChoice.Answer), 10),
			})
		case "MultipleChoice":
			bMultipleChoice, err := b.AsResponseBodyMultipleChoice()
			if err != nil {
				return nil, err
			}
			for _, a := range bMultipleChoice.Answer {
				res = append(res, &model.ResponseMeta{
					QuestionID: questions[i].ID,
					Data:       strconv.FormatInt(int64(a), 10),
				})
			}
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
