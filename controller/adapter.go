package controller

import (
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

func questionnaire2QuestionnaireDetail(questionnaires model.Questionnaires, adminUsers []string, adminGroups []string, targetUsers []string, targetGroups []string) openapi.QuestionnaireDetail {
	res := openapi.QuestionnaireDetail{
		Admins:      createUsersAndGroups(adminUsers, adminGroups),
		CreatedAt:   questionnaires.CreatedAt,
		Description: questionnaires.Description,
		// IsAllowingMultipleResponses: questionnaires.IsAllowingMultipleResponses,
		// IsAnonymous:                 questionnaires.IsAnonymous,
		// IsPublished:                 questionnaires.IsPublished,
		ModifiedAt:          questionnaires.ModifiedAt,
		QuestionnaireId:     questionnaires.ID,
		Questions:           convertQuestions(questionnaires.Questions),
		Respondents:         convertRespondents(questionnaires.Respondents),
		ResponseDueDateTime: &questionnaires.ResTimeLimit.Time,
		ResponseViewableBy:  convertResSharedTo(questionnaires.ResSharedTo),
		Targets:             createUsersAndGroups(targetUsers, targetGroups),
		Title:               questionnaires.Title,
	}
	return res
}
