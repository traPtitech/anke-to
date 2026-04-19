package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/openapi"
	"github.com/traPtitech/anke-to/traq"
	"gopkg.in/guregu/null.v4"
)

// Questionnaire Questionnaireの構造体
type Questionnaire struct {
	model.IQuestionnaire
	model.ITarget
	model.ITargetGroup
	model.ITargetUser
	model.IAdministrator
	model.IAdministratorGroup
	model.IAdministratorUser
	model.IQuestion
	model.IOption
	model.IScaleLabel
	model.IValidation
	model.ITransaction
	model.IRespondent
	traq.IWebhook
	*Response
	*Reminder
}

func NewQuestionnaire(
	questionnaire model.IQuestionnaire,
	target model.ITarget,
	targetGroup model.ITargetGroup,
	targetUser model.ITargetUser,
	administrator model.IAdministrator,
	administratorGroup model.IAdministratorGroup,
	administratorUser model.IAdministratorUser,
	question model.IQuestion,
	option model.IOption,
	scaleLabel model.IScaleLabel,
	validation model.IValidation,
	transaction model.ITransaction,
	respondent model.IRespondent,
	webhook traq.IWebhook,
	response *Response,
	reminder *Reminder,
) *Questionnaire {
	return &Questionnaire{
		IQuestionnaire:      questionnaire,
		ITarget:             target,
		ITargetGroup:        targetGroup,
		ITargetUser:         targetUser,
		IAdministrator:      administrator,
		IAdministratorGroup: administratorGroup,
		IAdministratorUser:  administratorUser,
		IQuestion:           question,
		IOption:             option,
		IScaleLabel:         scaleLabel,
		IValidation:         validation,
		ITransaction:        transaction,
		IRespondent:         respondent,
		IWebhook:            webhook,
		Response:            response,
		Reminder:            reminder,
	}
}

const (
	MaxTitleLength               = 1024
	responseDueDateTimeTolerance = 5 * time.Second
)

func normalizeResponseDueDateTime(responseDueDateTime *null.Time, now time.Time) error {
	if !responseDueDateTime.Valid {
		return nil
	}

	if responseDueDateTime.ValueOrZero().Before(now) {
		if now.Sub(responseDueDateTime.ValueOrZero()) > responseDueDateTimeTolerance {
			return errors.New("invalid resTimeLimit")
		}
		responseDueDateTime.Time = now
	}

	return nil
}

func maxLengthPattern(maxLength *int) string {
	if maxLength == nil {
		return ""
	}
	return "^.{0," + strconv.Itoa(*maxLength) + "}$"
}

func formatNumberBound(value *float64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatFloat(*value, 'f', -1, 64)
}

func (q *Questionnaire) GetQuestionnaires(ctx echo.Context, userID string, params openapi.GetQuestionnairesParams) (openapi.QuestionnaireList, error) {
	res := openapi.QuestionnaireList{
		Questionnaires: []openapi.QuestionnaireSummary{},
	}
	var sort string
	if params.Sort == nil {
		sort = ""
	} else {
		sort = string(*params.Sort)
	}
	var search string
	if params.Search == nil {
		search = ""
	} else {
		search = string(*params.Search)
	}
	var pageNum int
	if params.Page == nil {
		pageNum = 1
	} else {
		pageNum = int(*params.Page)
	}
	if pageNum < 1 {
		pageNum = 1
	}

	var onlyTargetingMe, onlyAdministratedByMe, notOverDue bool
	if params.OnlyTargetingMe == nil {
		onlyTargetingMe = false
	} else {
		onlyTargetingMe = *params.OnlyTargetingMe
	}
	if params.OnlyAdministratedByMe == nil {
		onlyAdministratedByMe = false
	} else {
		onlyAdministratedByMe = *params.OnlyAdministratedByMe
	}
	if params.NotOverDue == nil {
		notOverDue = false
	} else {
		notOverDue = *params.NotOverDue
	}

	var hasMyResponse, hasMyDraft, isDraft *bool
	hasMyResponse = params.HasMyResponse
	hasMyDraft = params.HasMyDraft
	isDraft = params.IsDraft
	// When fetching draft (unpublished) questionnaires, restrict to ones the user administrates
	if isDraft != nil && *isDraft {
		onlyAdministratedByMe = true
	}
	countOnly := params.CountOnly != nil && *params.CountOnly

	questionnaireList, totalRecords, pageMax, err := q.IQuestionnaire.GetQuestionnaires(ctx.Request().Context(), userID, sort, search, pageNum, onlyTargetingMe, onlyAdministratedByMe, notOverDue, hasMyResponse, hasMyDraft, isDraft, countOnly)
	if err != nil {
		return res, err
	}
	res.PageMax = pageMax
	res.TotalRecords = totalRecords
	if countOnly || len(questionnaireList) == 0 {
		return res, nil
	}

	questionnaireIDs := make([]int, 0, len(questionnaireList))
	for _, questionnaire := range questionnaireList {
		questionnaireIDs = append(questionnaireIDs, questionnaire.ID)
	}

	targets, err := q.ITarget.GetTargets(ctx.Request().Context(), questionnaireIDs)
	if err != nil {
		return res, err
	}
	respondents, err := q.IRespondent.GetRespondentsUserIDs(ctx.Request().Context(), questionnaireIDs)
	if err != nil {
		return res, err
	}
	myRespondents, err := q.GetRespondentInfos(ctx.Request().Context(), userID, questionnaireIDs...)
	if err != nil {
		return res, err
	}

	targetsByQuestionnaireID := make(map[int][]model.Targets, len(questionnaireList))
	for _, target := range targets {
		targetsByQuestionnaireID[target.QuestionnaireID] = append(targetsByQuestionnaireID[target.QuestionnaireID], target)
	}
	respondentsByQuestionnaireID := make(map[int][]model.Respondents, len(questionnaireList))
	for _, respondent := range respondents {
		respondentsByQuestionnaireID[respondent.QuestionnaireID] = append(respondentsByQuestionnaireID[respondent.QuestionnaireID], respondent)
	}
	myRespondentsByQuestionnaireID := make(map[int][]model.RespondentInfo, len(questionnaireList))
	for _, respondent := range myRespondents {
		myRespondentsByQuestionnaireID[respondent.QuestionnaireID] = append(myRespondentsByQuestionnaireID[respondent.QuestionnaireID], respondent)
	}

	for _, questionnaire := range questionnaireList {
		allRespondend := len(targetsByQuestionnaireID[questionnaire.ID]) == 0 || isAllTargetsReponded(targetsByQuestionnaireID[questionnaire.ID], respondentsByQuestionnaireID[questionnaire.ID])
		hasMyDraftForQuestionnaire := false
		hasMyResponseForQuestionnaire := false
		respondendDateTimeByMe := null.Time{}

		for _, respondent := range myRespondentsByQuestionnaireID[questionnaire.ID] {
			if !respondent.SubmittedAt.Valid {
				hasMyDraftForQuestionnaire = true
				continue
			}
			if !respondendDateTimeByMe.Valid {
				respondendDateTimeByMe = respondent.SubmittedAt
			}
			hasMyResponseForQuestionnaire = true
		}

		res.Questionnaires = append(res.Questionnaires, *questionnaireInfo2questionnaireSummary(questionnaire, allRespondend, hasMyDraftForQuestionnaire, hasMyResponseForQuestionnaire, respondendDateTimeByMe))
	}
	return res, nil
}

func (q *Questionnaire) PostQuestionnaire(c echo.Context, params openapi.PostQuestionnaireJSONRequestBody) (openapi.QuestionnaireDetail, error) {
	responseDueDateTime := null.Time{}
	if params.ResponseDueDateTime != nil {
		responseDueDateTime.Valid = true
		responseDueDateTime.Time = *params.ResponseDueDateTime
	}
	if err := normalizeResponseDueDateTime(&responseDueDateTime, time.Now()); err != nil {
		c.Logger().Infof("invalid resTimeLimit: %+v", responseDueDateTime)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusBadRequest, "invalid resTimeLimit")
	}

	questionnaireID := 0
	var err error

	if len(params.Title) == 0 || len(params.Title) > MaxTitleLength {
		c.Logger().Infof("invalid title: %+v", params.Title)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusBadRequest, "invalid title")
	}

	var notificationMessages []string
	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		questionnaireID, err = q.InsertQuestionnaire(ctx, params.Title, params.Description, responseDueDateTime, convertResponseViewableBy(params.ResponseViewableBy), params.IsPublished, params.IsAnonymous, params.IsDuplicateAnswerAllowed)
		if err != nil {
			c.Logger().Errorf("failed to insert questionnaire: %+v", err)
			return err
		}
		allTargetUsers, err := rollOutUsersAndGroups(params.Target.Users, params.Target.Groups)
		if err != nil {
			c.Logger().Errorf("failed to roll out users and groups: %+v", err)
			return err
		}
		targetGroupNames, err := uuid2GroupNames(params.Target.Groups)
		if err != nil {
			c.Logger().Errorf("failed to get group names: %+v", err)
			return err
		}
		err = q.InsertTargets(ctx, questionnaireID, allTargetUsers)
		if err != nil {
			c.Logger().Errorf("failed to insert targets: %+v", err)
			return err
		}
		err = q.InsertTargetUsers(ctx, questionnaireID, params.Target.Users)
		if err != nil {
			c.Logger().Errorf("failed to insert target groups: %+v", err)
			return err
		}
		err = q.InsertTargetGroups(ctx, questionnaireID, params.Target.Groups)
		if err != nil {
			c.Logger().Errorf("failed to insert target groups: %+v", err)
			return err
		}
		allAdminUsers, err := rollOutUsersAndGroups(params.Admin.Users, params.Admin.Groups)
		if err != nil {
			c.Logger().Errorf("failed to roll out administrators: %+v", err)
			return err
		}
		if len(allAdminUsers) == 0 {
			c.Logger().Errorf("no administrators")
			return errors.New("no administrators")
		}
		adminGroupNames, err := uuid2GroupNames(params.Admin.Groups)
		if err != nil {
			c.Logger().Errorf("failed to get group names: %+v", err)
			return err
		}
		err = q.InsertAdministrators(ctx, questionnaireID, allAdminUsers)
		if err != nil {
			c.Logger().Errorf("failed to insert administrators: %+v", err)
			return err
		}
		err = q.InsertAdministratorUsers(ctx, questionnaireID, params.Admin.Users)
		if err != nil {
			c.Logger().Errorf("failed to insert administrator users: %+v", err)
			return err
		}
		err = q.InsertAdministratorGroups(ctx, questionnaireID, params.Admin.Groups)
		if err != nil {
			c.Logger().Errorf("failed to insert administrator groups: %+v", err)
			return err
		}
		for questoinNum, question := range params.Questions {
			b, err := question.MarshalJSON()
			if err != nil {
				c.Logger().Errorf("failed to marshal new question: %+v", err)
				return err
			}
			var questionParsed map[string]interface{}
			err = json.Unmarshal([]byte(b), &questionParsed)
			if err != nil {
				c.Logger().Errorf("failed to unmarshal new question: %+v", err)
				return err
			}
			questionTypeRaw, ok := questionParsed["question_type"]
			if !ok {
				c.Logger().Errorf("question type is required")
				return errors.New("question type is required")
			}
			questionType, ok := questionTypeRaw.(string)
			if !ok {
				c.Logger().Errorf("question type must be string")
				return errors.New("question type must be string")
			}
			switch questionType {
			case "Text":
				questionType = "Text"
			case "TextLong":
				questionType = "TextArea"
			case "Number":
				questionType = "Number"
			case "SingleChoice":
				questionType = "MultipleChoice"
			case "MultipleChoice":
				questionType = "Checkbox"
			case "Scale":
				questionType = "LinearScale"
			default:
				c.Logger().Errorf("invalid question type")
				return errors.New("invalid question type")
			}
			questionID, err := q.InsertQuestion(ctx, questionnaireID, 1, questoinNum+1, questionType, question.Title, question.Description, question.IsRequired)
			if err != nil {
				c.Logger().Errorf("failed to insert question: %+v", err)
				return err
			}

			// insert validations
			switch questionType {
			case "MultipleChoice":
				b, err := question.AsQuestionSettingsSingleChoice()
				if err != nil {
					c.Logger().Errorf("failed to get question settings: %+v", err)
					return errors.New("failed to get question settings")
				}
				for i, v := range b.Options {
					err := q.IOption.InsertOption(ctx, questionID, i+1, v)
					if err != nil {
						c.Logger().Errorf("failed to insert option: %+v", err)
						return errors.New("failed to insert option")
					}
				}
			case "Checkbox":
				b, err := question.AsQuestionSettingsMultipleChoice()
				if err != nil {
					c.Logger().Errorf("failed to get question settings: %+v", err)
					return errors.New("failed to get question settings")
				}
				for i, v := range b.Options {
					err := q.IOption.InsertOption(ctx, questionID, i+1, v)
					if err != nil {
						c.Logger().Errorf("failed to insert option: %+v", err)
						return errors.New("failed to insert option")
					}
				}
			case "LinearScale":
				b, err := question.AsQuestionSettingsScale()
				if err != nil {
					c.Logger().Errorf("failed to get question settings: %+v", err)
					return errors.New("failed to get question settings")
				}
				if b.MaxValue < b.MinValue {
					c.Logger().Errorf("invalid scale")
					return errors.New("invalid scale")
				}
				minLabel := ""
				maxLabel := ""
				if b.MinLabel != nil {
					minLabel = *b.MinLabel
				}
				if b.MaxLabel != nil {
					maxLabel = *b.MaxLabel
				}
				err = q.IScaleLabel.InsertScaleLabel(ctx, questionID,
					model.ScaleLabels{
						ScaleLabelLeft:  minLabel,
						ScaleLabelRight: maxLabel,
						ScaleMax:        b.MaxValue,
						ScaleMin:        b.MinValue,
					})
				if err != nil {
					c.Logger().Errorf("failed to insert scale label: %+v", err)
					return errors.New("failed to insert scale label")
				}
			case "Text":
				b, err := question.AsQuestionSettingsText()
				if err != nil {
					c.Logger().Errorf("failed to get question settings: %+v", err)
					return errors.New("failed to get question settings")
				}
				err = q.IValidation.InsertValidation(ctx, questionID,
					model.Validations{
						RegexPattern: maxLengthPattern(b.MaxLength),
					})
				if err != nil {
					c.Logger().Errorf("failed to insert validation: %+v", err)
					return errors.New("failed to insert validation")
				}
			case "TextArea":
				b, err := question.AsQuestionSettingsTextLong()
				if err != nil {
					c.Logger().Errorf("failed to get question settings: %+v", err)
					return errors.New("failed to get question settings")
				}
				err = q.IValidation.InsertValidation(ctx, questionID,
					model.Validations{
						RegexPattern: maxLengthPattern(b.MaxLength),
					})
				if err != nil {
					c.Logger().Errorf("failed to insert validation: %+v", err)
					return errors.New("failed to insert validation")
				}
			case "Number":
				b, err := question.AsQuestionSettingsNumber()
				if err != nil {
					c.Logger().Errorf("failed to get question settings: %+v", err)
					return errors.New("failed to get question settings")
				}
				// 数字かどうか，min<=maxになっているかどうか
				minValueStr := formatNumberBound(b.MinValue)
				maxValueStr := formatNumberBound(b.MaxValue)
				err = q.IValidation.CheckNumberValid(minValueStr, maxValueStr)
				if err != nil {
					c.Logger().Errorf("invalid number: %+v", err)
					return errors.New("invalid number")
				}
				err = q.IValidation.InsertValidation(ctx, questionID,
					model.Validations{
						MinBound: minValueStr,
						MaxBound: maxValueStr,
					})
				if err != nil {
					c.Logger().Errorf("failed to insert validation: %+v", err)
					return errors.New("failed to insert validation")
				}
			}
		}

		notificationMessages = createQuestionnaireMessage(
			questionnaireID,
			params.Title,
			params.Description,
			append(allAdminUsers, adminGroupNames...),
			responseDueDateTime,
			append(allTargetUsers, targetGroupNames...),
		)

		if params.ResponseDueDateTime != nil && params.IsPublished {
			dueDateTime := responseDueDateTime.Time
			err = q.PushReminder(questionnaireID, &dueDateTime)
			if err != nil {
				c.Logger().Errorf("failed to push reminder: %+v", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		c.Logger().Errorf("failed to create a questionnaire: %+v", err)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to create a questionnaire")
	}

	// Send traQ notifications after the DB transaction commits.
	// Failures are only logged; the questionnaire creation itself is treated as successful.
	for _, message := range notificationMessages {
		if err := q.PostMessage(message); err != nil {
			c.Logger().Errorf("failed to post questionnaire creation message (questionnaireID: %d): %+v", questionnaireID, err)
		}
	}

	questionnaireInfo, targets, targetUsers, targetGroups, admins, adminUsers, adminGroups, respondents, err := q.GetQuestionnaireInfo(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get questionnaire info: %+v", err)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to get questionnaire info")
	}

	questionnaireDetail, err := questionnaire2QuestionnaireDetail(*questionnaireInfo, admins, adminUsers, adminGroups, targets, targetUsers, targetGroups, respondents)
	if err != nil {
		c.Logger().Errorf("failed to convert questionnaire to questionnaire detail: %+v", err)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to convert questionnaire to questionnaire detail")
	}
	return questionnaireDetail, nil
}
func (q *Questionnaire) GetQuestionnaire(ctx echo.Context, questionnaireID int) (openapi.QuestionnaireDetail, error) {
	questionnaireInfo, targets, targetUsers, targetGroups, admins, adminUsers, adminGroups, respondents, err := q.GetQuestionnaireInfo(ctx.Request().Context(), questionnaireID)
	if err != nil {
		return openapi.QuestionnaireDetail{}, err
	}
	questionnaireDetail, err := questionnaire2QuestionnaireDetail(*questionnaireInfo, admins, adminUsers, adminGroups, targets, targetUsers, targetGroups, respondents)
	if err != nil {
		ctx.Logger().Errorf("failed to convert questionnaire to questionnaire detail: %+v", err)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to convert questionnaire to questionnaire detail")
	}
	return questionnaireDetail, nil
}

func (q *Questionnaire) EditQuestionnaire(c echo.Context, questionnaireID int, params openapi.EditQuestionnaireJSONRequestBody) error {
	// unable to change the questionnaire from anonymous to non-anonymous
	isAnonymous, err := q.GetResponseIsAnonymousByQuestionnaireID(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get anonymous info: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get anonymous info")
	}
	if isAnonymous && !params.IsAnonymous {
		c.Logger().Info("unable to change the questionnaire from anonymous to non-anonymous")
		return echo.NewHTTPError(http.StatusMethodNotAllowed, "unable to change the questionnaire from anonymous to non-anonymous")
	}

	responseDueDateTime := null.Time{}
	if params.ResponseDueDateTime != nil {
		responseDueDateTime.Valid = true
		responseDueDateTime.Time = *params.ResponseDueDateTime
	}
	if err := normalizeResponseDueDateTime(&responseDueDateTime, time.Now()); err != nil {
		c.Logger().Infof("invalid resTimeLimit: %+v", responseDueDateTime)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid resTimeLimit")
	}

	if len(params.Title) == 0 || len(params.Title) > MaxTitleLength {
		c.Logger().Infof("invalid title: %+v", params.Title)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid title")
	}

	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		err := q.UpdateQuestionnaire(ctx, params.Title, params.Description, responseDueDateTime, convertResponseViewableBy(params.ResponseViewableBy), questionnaireID, params.IsPublished, params.IsAnonymous, params.IsDuplicateAnswerAllowed)
		if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			c.Logger().Errorf("failed to update questionnaire: %+v", err)
			return err
		}
		if params.Target != nil {
			err = q.DeleteTargets(ctx, questionnaireID)
			if err != nil {
				c.Logger().Errorf("failed to delete targets: %+v", err)
				return err
			}
			err = q.DeleteTargetUsers(ctx, questionnaireID)
			if err != nil {
				c.Logger().Errorf("failed to delete target users: %+v", err)
				return err
			}
			err = q.DeleteTargetGroups(ctx, questionnaireID)
			if err != nil {
				c.Logger().Errorf("failed to delete target groups: %+v", err)
				return err
			}
			allTargetUsers, err := rollOutUsersAndGroups((*params.Target).Users, params.Target.Groups)
			if err != nil {
				c.Logger().Errorf("failed to roll out users and groups: %+v", err)
				return err
			}
			err = q.InsertTargets(ctx, questionnaireID, allTargetUsers)
			if err != nil {
				c.Logger().Errorf("failed to insert targets: %+v", err)
				return err
			}
			err = q.InsertTargetUsers(ctx, questionnaireID, params.Target.Users)
			if err != nil {
				c.Logger().Errorf("failed to insert target users: %+v", err)
				return err
			}
			err = q.InsertTargetGroups(ctx, questionnaireID, params.Target.Groups)
			if err != nil {
				c.Logger().Errorf("failed to insert target groups: %+v", err)
				return err
			}
		}
		if params.Admin != nil {
			err = q.DeleteAdministrators(ctx, questionnaireID)
			if err != nil {
				c.Logger().Errorf("failed to delete administrators: %+v", err)
				return err
			}
			err = q.DeleteAdministratorUsers(ctx, questionnaireID)
			if err != nil {
				c.Logger().Errorf("failed to delete administrator users: %+v", err)
				return err
			}
			err = q.DeleteAdministratorGroups(ctx, questionnaireID)
			if err != nil {
				c.Logger().Errorf("failed to delete administrator groups: %+v", err)
				return err
			}
			allAdminUsers, err := rollOutUsersAndGroups(params.Admin.Users, params.Admin.Groups)
			if err != nil {
				c.Logger().Errorf("failed to roll out administrators: %+v", err)
				return err
			}
			if len(allAdminUsers) == 0 {
				c.Logger().Errorf("no administrators")
				return errors.New("no administrators")
			}
			err = q.InsertAdministrators(ctx, questionnaireID, allAdminUsers)
			if err != nil {
				c.Logger().Errorf("failed to insert administrators: %+v", err)
				return err
			}
			err = q.InsertAdministratorUsers(ctx, questionnaireID, params.Admin.Users)
			if err != nil {
				c.Logger().Errorf("failed to insert administrator users: %+v", err)
				return err
			}
			err = q.InsertAdministratorGroups(ctx, questionnaireID, params.Admin.Groups)
			if err != nil {
				c.Logger().Errorf("failed to insert administrator groups: %+v", err)
				return err
			}
		}

		var ifQuestionExist = make(map[int]bool)
		for questoinNum, question := range params.Questions {
			b, err := question.MarshalJSON()
			if err != nil {
				c.Logger().Errorf("failed to marshal new question: %+v", err)
				return err
			}
			var questionParsed map[string]interface{}
			err = json.Unmarshal([]byte(b), &questionParsed)
			if err != nil {
				c.Logger().Errorf("failed to unmarshal new question: %+v", err)
				return err
			}
			questionTypeRaw, ok := questionParsed["question_type"]
			if !ok {
				c.Logger().Errorf("question type is required")
				return errors.New("question type is required")
			}
			questionType, ok := questionTypeRaw.(string)
			if !ok {
				c.Logger().Errorf("question type must be string")
				return errors.New("question type must be string")
			}
			switch questionType {
			case "Text":
				questionType = "Text"
			case "TextLong":
				questionType = "TextArea"
			case "Number":
				questionType = "Number"
			case "SingleChoice":
				questionType = "MultipleChoice"
			case "MultipleChoice":
				questionType = "Checkbox"
			case "Scale":
				questionType = "LinearScale"
			default:
				c.Logger().Errorf("invalid question type")
				return errors.New("invalid question type")
			}
			if question.QuestionId == nil {
				questionID, err := q.InsertQuestion(ctx, questionnaireID, 1, questoinNum+1, questionType, question.Title, question.Description, question.IsRequired)
				if err != nil {
					c.Logger().Errorf("failed to insert question: %+v", err)
					return err
				}
				ifQuestionExist[questionID] = true
				// insert validations
				switch questionType {
				case "MultipleChoice":
					b, err := question.AsQuestionSettingsSingleChoice()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					for i, v := range b.Options {
						err := q.IOption.InsertOption(ctx, questionID, i+1, v)
						if err != nil {
							c.Logger().Errorf("failed to insert option: %+v", err)
							return errors.New("failed to insert option")
						}
					}
				case "Checkbox":
					b, err := question.AsQuestionSettingsMultipleChoice()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					for i, v := range b.Options {
						err := q.IOption.InsertOption(ctx, questionID, i+1, v)
						if err != nil {
							c.Logger().Errorf("failed to insert option: %+v", err)
							return errors.New("failed to insert option")
						}
					}
				case "LinearScale":
					b, err := question.AsQuestionSettingsScale()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					if b.MaxValue < b.MinValue {
						c.Logger().Errorf("invalid scale")
						return errors.New("invalid scale")
					}
					minLabel := ""
					maxLabel := ""
					if b.MinLabel != nil {
						minLabel = *b.MinLabel
					}
					if b.MaxLabel != nil {
						maxLabel = *b.MaxLabel
					}
					err = q.IScaleLabel.InsertScaleLabel(ctx, questionID,
						model.ScaleLabels{
							ScaleLabelLeft:  minLabel,
							ScaleLabelRight: maxLabel,
							ScaleMax:        b.MaxValue,
							ScaleMin:        b.MinValue,
						})
					if err != nil {
						c.Logger().Errorf("failed to insert scale label: %+v", err)
						return errors.New("failed to insert scale label")
					}
				case "Text":
					b, err := question.AsQuestionSettingsText()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					err = q.IValidation.InsertValidation(ctx, questionID,
						model.Validations{
							RegexPattern: maxLengthPattern(b.MaxLength),
						})
					if err != nil {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				case "TextArea":
					b, err := question.AsQuestionSettingsTextLong()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					err = q.IValidation.InsertValidation(ctx, questionID,
						model.Validations{
							RegexPattern: maxLengthPattern(b.MaxLength),
						})
					if err != nil {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				case "Number":
					b, err := question.AsQuestionSettingsNumber()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					// 数字かどうか，min<=maxになっているかどうか
					minValueStr := formatNumberBound(b.MinValue)
					maxValueStr := formatNumberBound(b.MaxValue)
					err = q.IValidation.CheckNumberValid(minValueStr, maxValueStr)
					if err != nil {
						c.Logger().Errorf("invalid number: %+v", err)
						return errors.New("invalid number")
					}
					err = q.IValidation.InsertValidation(ctx, questionID,
						model.Validations{
							MinBound: minValueStr,
							MaxBound: maxValueStr,
						})
					if err != nil {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				}
			} else {
				ifQuestionExist[*question.QuestionId] = true
				err = q.UpdateQuestion(ctx, questionnaireID, 1, questoinNum+1, questionType, question.Title, question.Description, question.IsRequired, *question.QuestionId)
				if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
					c.Logger().Errorf("failed to update question: %+v", err)
					return err
				}
				// update validations
				switch questionType {
				case "MultipleChoice":
					b, err := question.AsQuestionSettingsSingleChoice()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					err = q.IOption.UpdateOptions(ctx, b.Options, *question.QuestionId)
					if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
						c.Logger().Errorf("failed to update options: %+v", err)
						return errors.New("failed to update options")
					}
				case "Checkbox":
					b, err := question.AsQuestionSettingsMultipleChoice()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					err = q.IOption.UpdateOptions(ctx, b.Options, *question.QuestionId)
					if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
						c.Logger().Errorf("failed to update options: %+v", err)
						return errors.New("failed to update options")
					}
				case "LinearScale":
					b, err := question.AsQuestionSettingsScale()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					if b.MaxValue < b.MinValue {
						c.Logger().Errorf("invalid scale")
						return errors.New("invalid scale")
					}
					minLabel := ""
					maxLabel := ""
					if b.MinLabel != nil {
						minLabel = *b.MinLabel
					}
					if b.MaxLabel != nil {
						maxLabel = *b.MaxLabel
					}
					err = q.IScaleLabel.UpdateScaleLabel(ctx, *question.QuestionId,
						model.ScaleLabels{
							ScaleLabelLeft:  minLabel,
							ScaleLabelRight: maxLabel,
							ScaleMax:        b.MaxValue,
							ScaleMin:        b.MinValue,
						})
					if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
						c.Logger().Errorf("failed to insert scale label: %+v", err)
						return errors.New("failed to insert scale label")
					}
				case "Text":
					b, err := question.AsQuestionSettingsText()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					err = q.IValidation.UpdateValidation(ctx, *question.QuestionId,
						model.Validations{
							RegexPattern: maxLengthPattern(b.MaxLength),
						})
					if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				case "TextArea":
					b, err := question.AsQuestionSettingsTextLong()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					err = q.IValidation.UpdateValidation(ctx, *question.QuestionId,
						model.Validations{
							RegexPattern: maxLengthPattern(b.MaxLength),
						})
					if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				case "Number":
					b, err := question.AsQuestionSettingsNumber()
					if err != nil {
						c.Logger().Errorf("failed to get question settings: %+v", err)
						return errors.New("failed to get question settings")
					}
					// 数字かどうか，min<=maxになっているかどうか
					minValueStr := formatNumberBound(b.MinValue)
					maxValueStr := formatNumberBound(b.MaxValue)
					err = q.IValidation.CheckNumberValid(minValueStr, maxValueStr)
					if err != nil {
						c.Logger().Errorf("invalid number: %+v", err)
						return errors.New("invalid number")
					}
					err = q.IValidation.UpdateValidation(ctx, *question.QuestionId,
						model.Validations{
							MinBound: minValueStr,
							MaxBound: maxValueStr,
						})
					if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				}
			}
		}
		questions, err := q.IQuestion.GetQuestions(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to get questions: %+v", err)
			return err
		}
		for _, question := range questions {
			if !ifQuestionExist[question.ID] {
				err = q.DeleteQuestion(ctx, question.ID)
				if err != nil {
					c.Logger().Errorf("failed to delete question: %+v", err)
					return err
				}
			}
		}

		err = q.DeleteReminder(questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete reminder: %+v", err)
			return err
		}
		if params.ResponseDueDateTime != nil && params.IsPublished {
			dueDateTime := responseDueDateTime.Time
			err = q.PushReminder(questionnaireID, &dueDateTime)
			if err != nil {
				c.Logger().Errorf("failed to push reminder: %+v", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		c.Logger().Errorf("failed to update a questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update a questionnaire")
	}

	return nil
}

func (q *Questionnaire) DeleteQuestionnaire(c echo.Context, questionnaireID int) error {
	err := q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		respondentDetails, err := q.GetRespondentDetails(ctx, questionnaireID, "", false, "", nil)
		if err != nil {
			c.Logger().Errorf("failed to get respondent details: %+v", err)
			return err
		}

		err = q.IQuestionnaire.DeleteQuestionnaire(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete questionnaire: %+v", err)
			return err
		}

		err = q.DeleteTargets(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete targets: %+v", err)
			return err
		}

		err = q.DeleteTargetUsers(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete target users: %+v", err)
			return err
		}

		err = q.DeleteTargetGroups(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete target groups: %+v", err)
			return err
		}

		err = q.DeleteAdministrators(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrators: %+v", err)
			return err
		}

		err = q.DeleteAdministratorUsers(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrator users: %+v", err)
			return err
		}
		err = q.DeleteAdministratorGroups(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrator groups: %+v", err)
			return err
		}

		questions, err := q.GetQuestions(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to get questions: %+v", err)
			return err
		}
		for _, question := range questions {
			err = q.DeleteOptions(ctx, question.ID)
			if err != nil {
				c.Logger().Errorf("failed to delete options: %+v", err)
				return err
			}

			if question.Type == "LinearScale" {
				err = q.DeleteScaleLabel(ctx, question.ID)
				if err != nil {
					c.Logger().Errorf("failed to delete scale label: %+v", err)
					return err
				}
			}

			if question.Type == "Text" || question.Type == "TextArea" || question.Type == "Number" {
				err = q.DeleteValidation(ctx, question.ID)
				if err != nil {
					c.Logger().Errorf("failed to delete validation: %+v", err)
					return err
				}
			}

			err = q.DeleteQuestion(ctx, question.ID)
			if err != nil {
				c.Logger().Errorf("failed to delete question: %+v", err)
				return err
			}
		}

		for _, respondentDetail := range respondentDetails {
			err = q.IResponse.DeleteResponse(ctx, respondentDetail.ResponseID)
			if err != nil && !errors.Is(err, model.ErrNoRecordDeleted) {
				c.Logger().Errorf("failed to delete responses: %+v", err)
				return err
			}

			err = q.DeleteRespondent(ctx, respondentDetail.ResponseID)
			if err != nil {
				c.Logger().Errorf("failed to delete respondents: %+v", err)
				return err
			}
		}

		err = model.NewReminderTarget().DeleteReminderTargets(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete reminder targets: %+v", err)
			return err
		}

		err = q.DeleteReminder(questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete reminder: %+v", err)
			return err
		}

		return nil
	})
	if err != nil {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			return httpError
		}

		c.Logger().Errorf("failed to delete questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete a questionnaire")
	}
	return nil
}

func (q *Questionnaire) GetQuestionnaireMyRemindStatus(c echo.Context, questionnaireID int, userID string) (bool, error) {
	_, _, _, _, _, _, _, _, err := q.GetQuestionnaireInfo(c.Request().Context(), questionnaireID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return false, echo.NewHTTPError(http.StatusNotFound, "questionnaire not found")
		}
		c.Logger().Errorf("failed to get questionnaire info: %+v", err)
		return false, echo.NewHTTPError(http.StatusInternalServerError, "failed to check remind status")
	}

	reminderTarget, err := model.NewReminderTarget().GetReminderTarget(c.Request().Context(), questionnaireID, userID)
	if err == nil {
		return !reminderTarget.IsCanceled, nil
	}
	if err != nil && !errors.Is(err, model.ErrRecordNotFound) {
		c.Logger().Errorf("failed to get reminder target status: %+v", err)
		return false, echo.NewHTTPError(http.StatusInternalServerError, "failed to check remind status")
	}

	status, err := q.GetTargetsCancelStatus(c.Request().Context(), questionnaireID, []string{userID})
	if err != nil {
		if errors.Is(err, model.ErrTargetNotFound) {
			return false, nil
		}
		c.Logger().Errorf("failed to check remind status: %+v", err)
		return false, echo.NewHTTPError(http.StatusInternalServerError, "failed to check remind status")
	}

	return !status[0].IsCanceled, nil
}

func (q *Questionnaire) EditQuestionnaireMyRemindStatus(c echo.Context, questionnaireID int, userID string, isRemindEnabled bool) error {
	_, _, _, _, _, _, _, _, err := q.GetQuestionnaireInfo(c.Request().Context(), questionnaireID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "questionnaire not found")
		}
		c.Logger().Errorf("failed to get questionnaire info: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update remind status")
	}

	err = model.NewReminderTarget().UpsertReminderTarget(c.Request().Context(), questionnaireID, userID, !isRemindEnabled)
	if err != nil {
		c.Logger().Errorf("failed to update remind status: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update remind status")
	}

	return nil
}

func (q *Questionnaire) GetQuestionnaireResponses(c echo.Context, questionnaireID int, params openapi.GetQuestionnaireResponsesParams, userID string) (openapi.Responses, error) {
	res := []openapi.Response{}
	var sort string
	var onlyMyResponse bool
	if params.Sort != nil {
		sort = string(*params.Sort)
	} else {
		sort = ""
	}
	if params.OnlyMyResponse != nil {
		onlyMyResponse = *params.OnlyMyResponse
	} else {
		onlyMyResponse = false
	}
	respondentDetails, err := q.GetRespondentDetails(c.Request().Context(), questionnaireID, sort, onlyMyResponse, userID, params.IsDraft)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return res, echo.NewHTTPError(http.StatusNotFound, "respondent not found")
		}
		c.Logger().Errorf("failed to get respondent details: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, "failed to get respondent details")
	}

	for _, respondentDetail := range respondentDetails {
		response, err := respondentDetail2Response(c, respondentDetail)
		if err != nil {
			c.Logger().Errorf("failed to convert respondent detail to response: %+v", err)
			return res, echo.NewHTTPError(http.StatusInternalServerError, "failed to convert respondent detail to response")
		}
		res = append(res, response)
	}

	return res, nil
}

func (q *Questionnaire) PostQuestionnaireResponse(c echo.Context, questionnaireID int, params openapi.PostQuestionnaireResponseJSONRequestBody, userID string) (openapi.Response, error) {
	res := openapi.Response{}

	limit, err := q.GetQuestionnaireLimit(c.Request().Context(), questionnaireID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info("questionnaire not found")
			return res, echo.NewHTTPError(http.StatusNotFound, err)
		}
		c.Logger().Errorf("failed to get questionnaire limit: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 回答期限を過ぎていたらエラー
	if limit.Valid && limit.Time.Before(time.Now()) {
		c.Logger().Info("expired questionnaire")
		return res, echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("expired questionnaire"))
	}

	questions, err := q.IQuestion.GetQuestions(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get questions: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	responseMetas, err := responseBody2ResponseMetas(params.Body, questions)
	if err != nil {
		c.Logger().Infof("invalid response body: %+v", err)
		return res, echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid response body: %w", err))
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

	validations, err := q.IValidation.GetValidations(c.Request().Context(), questionIDs)
	if err != nil {
		c.Logger().Errorf("failed to get validations: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	validationMap := make(map[int]model.Validations, len(validations))
	for _, validation := range validations {
		validationMap[validation.QuestionID] = validation
	}

	options, err := q.IOption.GetOptions(c.Request().Context(), questionIDs)
	if err != nil {
		c.Logger().Errorf("failed to get options: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	optionMap := make(map[int][]model.Options, len(options))
	for _, option := range options {
		optionMap[option.QuestionID] = append(optionMap[option.QuestionID], option)
	}

	scaleLabels, err := q.IScaleLabel.GetScaleLabels(c.Request().Context(), questionIDs)
	if err != nil {
		c.Logger().Errorf("failed to get scale labels: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, scaleLabel := range scaleLabels {
		scaleLabelMap[scaleLabel.QuestionID] = scaleLabel
	}

	for _, responseMeta := range responseMetas {
		questionRequired[responseMeta.QuestionID] = false
		switch questionTypes[responseMeta.QuestionID] {
		case "Text", "TextArea":
			if !params.IsDraft {
				validation, ok := validationMap[responseMeta.QuestionID]
				if !ok {
					validation = model.Validations{}
				}
				err := q.IValidation.CheckTextValidation(validation, responseMeta.Data)
				if err != nil {
					if errors.Is(err, model.ErrTextMatching) {
						c.Logger().Errorf("invalid text: %+v", err)
						return res, echo.NewHTTPError(http.StatusBadRequest, err)
					}
					c.Logger().Errorf("invalid text: %+v", err)
					return res, echo.NewHTTPError(http.StatusBadRequest, err)
				}
			}
		case "Number":
			if !params.IsDraft {
				validation, ok := validationMap[responseMeta.QuestionID]
				if !ok {
					validation = model.Validations{}
				}
				err := q.IValidation.CheckNumberValidation(validation, responseMeta.Data)
				if err != nil {
					if errors.Is(err, model.ErrInvalidNumber) {
						c.Logger().Errorf("invalid number: %+v", err)
						return res, echo.NewHTTPError(http.StatusBadRequest, err)
					}
					c.Logger().Errorf("invalid number: %+v", err)
					return res, echo.NewHTTPError(http.StatusBadRequest, err)
				}
			}
		case "Checkbox", "MultipleChoice":
		case "LinearScale":
			if !params.IsDraft {
				label, ok := scaleLabelMap[responseMeta.QuestionID]
				if !ok {
					c.Logger().Errorf("scale label not found for question %d", responseMeta.QuestionID)
					return res, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("scale label not found for question %d", responseMeta.QuestionID))
				}
				err := q.IScaleLabel.CheckScaleLabel(label, responseMeta.Data)
				if err != nil {
					c.Logger().Errorf("invalid scale: %+v", err)
					return res, echo.NewHTTPError(http.StatusBadRequest, err)
				}
			}
		default:
			c.Logger().Errorf("invalid question id: %+v", responseMeta.QuestionID)
			return res, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid question id: %d", responseMeta.QuestionID))
		}
	}

	if !params.IsDraft {
		for _, question := range questions {
			if questionRequired[question.ID] {
				c.Logger().Errorf("required question is not answered: %+v", question.ID)
				return res, echo.NewHTTPError(http.StatusBadRequest, "required question is not answered")
			}
		}
	}

	var submittedAt time.Time
	//一時保存のときはnull
	if params.IsDraft {
		submittedAt = time.Time{}
	} else {
		submittedAt = time.Now()
	}

	var responseID int
	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		var err error
		responseID, err = q.InsertRespondent(ctx, userID, questionnaireID, null.NewTime(submittedAt, !params.IsDraft))
		if err != nil {
			c.Logger().Errorf("failed to insert respondant: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if len(responseMetas) > 0 {
			err = q.InsertResponses(ctx, responseID, responseMetas)
			if err != nil {
				c.Logger().Errorf("failed to insert responses: %+v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		return nil
	})
	if err != nil {
		c.Logger().Errorf("failed to insert response: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert response: %w", err))
	}

	response, err := q.GetResponse(c, responseID)
	if err != nil {
		c.Logger().Errorf("failed to get response: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response: %w", err))
	}

	return response, nil
}

func createQuestionnaireMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit null.Time, targets []string) []string {
	var resTimeLimitText string
	if resTimeLimit.Valid {
		resTimeLimitText = resTimeLimit.Time.Local().Format("2006/01/02 15:04")
	} else {
		resTimeLimitText = "なし"
	}

	prefix := fmt.Sprintf(
		"### アンケート『[%s](https://anke-to.trap.jp/questionnaires/%d)』が作成されました\n#### 管理者\n%s\n#### 説明\n%s\n#### 回答期限\n%s",
		title,
		questionnaireID,
		strings.Join(administrators, ","),
		description,
		resTimeLimitText,
	)
	suffix := fmt.Sprintf("\n#### 回答リンク\nhttps://anke-to.trap.jp/responses/new/%d", questionnaireID)

	return createMessagesFromTargets(prefix, suffix, targets, traq.MessageLimit)
}

func createReminderMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit time.Time, targets []string, leftTimeText string) []string {
	resTimeLimitText := resTimeLimit.Local().Format("2006/01/02 15:04")

	prefix := fmt.Sprintf(
		"### アンケート『[%s](https://anke-to.trap.jp/questionnaires/%d)』の回答期限が迫っています!\n==残り%sです!==\n#### 管理者\n%s\n#### 説明\n%s\n#### 回答期限\n%s",
		title,
		questionnaireID,
		leftTimeText,
		strings.Join(administrators, ","),
		description,
		resTimeLimitText,
	)
	suffix := fmt.Sprintf("\n#### 回答リンク\nhttps://anke-to.trap.jp/responses/new/%d", questionnaireID)

	return createMessagesFromTargets(prefix, suffix, targets, traq.MessageLimit)
}

// createMessagesFromTargets は対象者リストをlimit文字以内に収まるよう分割し、
// それぞれに完全なヘッダーとフッターを付けた複数のメッセージを返す。
func createMessagesFromTargets(prefix, suffix string, targets []string, limit int) []string {
	const targetsHeader = "\n#### 対象者\n"

	if len(targets) == 0 {
		return []string{prefix + targetsHeader + "なし" + suffix}
	}

	allTargetsText := "@" + strings.Join(targets, " @")
	full := prefix + targetsHeader + allTargetsText + suffix
	if len([]rune(full)) <= limit {
		return []string{full}
	}

	available := limit - len([]rune(prefix)) - len([]rune(targetsHeader)) - len([]rune(suffix))

	var messages []string
	var group []string
	groupLen := 0

	for _, target := range targets {
		mention := "@" + target
		addLen := len([]rune(mention))
		if groupLen > 0 {
			addLen++ // スペース区切り分
		}

		if groupLen+addLen > available && len(group) > 0 {
			messages = append(messages, prefix+targetsHeader+"@"+strings.Join(group, " @")+suffix)
			group = []string{target}
			groupLen = len([]rune(mention))
		} else {
			group = append(group, target)
			groupLen += addLen
		}
	}

	if len(group) > 0 {
		messages = append(messages, prefix+targetsHeader+"@"+strings.Join(group, " @")+suffix)
	}

	return messages
}
