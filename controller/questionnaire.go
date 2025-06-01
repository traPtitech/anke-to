package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
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

const MaxTitleLength = 50

func (q *Questionnaire) GetQuestionnaires(ctx echo.Context, userID string, params openapi.GetQuestionnairesParams) (openapi.QuestionnaireList, error) {
	res := openapi.QuestionnaireList{}
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

	var onlyTargetingMe, onlyAdministratedByMe bool
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
	questionnaireList, pageMax, err := q.IQuestionnaire.GetQuestionnaires(ctx.Request().Context(), userID, sort, search, pageNum, onlyTargetingMe, onlyAdministratedByMe)
	if err != nil {
		return res, err
	}

	for _, questionnaire := range questionnaireList {
		targets, err := q.ITarget.GetTargets(ctx.Request().Context(), []int{questionnaire.ID})
		if err != nil {
			return res, err
		}
		allRespondend := false
		if len(targets) == 0 {
			allRespondend = true
		} else {
			respondents, err := q.IRespondent.GetRespondentsUserIDs(ctx.Request().Context(), []int{questionnaire.ID})
			if err != nil {
				return res, err
			}
			allRespondend = isAllTargetsReponded(targets, respondents)
		}

		hasMyDraft := false
		hasMyResponse := false
		respondendDateTimeByMe := null.Time{}

		myRespondents, err := q.GetRespondentInfos(ctx.Request().Context(), userID, questionnaire.ID)
		if err != nil {
			return res, err
		}
		for _, respondent := range myRespondents {
			if !respondent.SubmittedAt.Valid {
				hasMyDraft = true
			}
			if respondent.SubmittedAt.Valid {
				if !respondendDateTimeByMe.Valid {
					respondendDateTimeByMe = respondent.SubmittedAt
				}
				hasMyResponse = true
			}
		}

		res.PageMax = pageMax
		res.Questionnaires = append(res.Questionnaires, *questionnaireInfo2questionnaireSummary(questionnaire, allRespondend, hasMyDraft, hasMyResponse, respondendDateTimeByMe))
	}
	return res, nil
}

func (q *Questionnaire) PostQuestionnaire(c echo.Context, params openapi.PostQuestionnaireJSONRequestBody) (openapi.QuestionnaireDetail, error) {
	responseDueDateTime := null.Time{}
	if params.ResponseDueDateTime != nil {
		responseDueDateTime.Valid = true
		responseDueDateTime.Time = *params.ResponseDueDateTime
	}
	if responseDueDateTime.Valid {
		isBefore := responseDueDateTime.ValueOrZero().Before(time.Now())
		if isBefore {
			c.Logger().Infof("invalid resTimeLimit: %+v", responseDueDateTime)
			return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusBadRequest, "invalid resTimeLimit")
		}
	}

	questionnaireID := 0
	var err error

	if len(params.Title) == 0 || len(params.Title) > MaxTitleLength {
		c.Logger().Infof("invalid title: %+v", params.Title)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusBadRequest, "invalid title")
	}

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
			questionType := questionParsed["question_type"].(string)
			if questionType == "Text" {
				questionType = "Text"
			} else if questionType == "TextLong" {
				questionType = "TextArea"
			} else if questionType == "Number" {
				questionType = "Number"
			} else if questionType == "SingleChoice" {
				questionType = "MultipleChoice"
			} else if questionType == "MultipleChoice" {
				questionType = "Checkbox"
			} else if questionType == "Scale" {
				questionType = "LinearScale"
			} else {
				c.Logger().Errorf("invalid question type")
				return errors.New("invalid question type")
			}
			questionID, err := q.InsertQuestion(ctx, questionnaireID, 1, questoinNum+1, questionType, question.Body, question.IsRequired)
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
				err = q.IScaleLabel.InsertScaleLabel(ctx, questionID,
					model.ScaleLabels{
						ScaleLabelLeft:  *b.MinLabel,
						ScaleLabelRight: *b.MaxLabel,
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
						RegexPattern: "^.{0," + strconv.Itoa(*b.MaxLength) + "}$",
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
						RegexPattern: "^.{0," + fmt.Sprintf("%d", *b.MaxLength) + "}$",
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
				err = q.IValidation.CheckNumberValid(strconv.Itoa(*b.MinValue), strconv.Itoa(*b.MaxValue))
				if err != nil {
					c.Logger().Errorf("invalid number: %+v", err)
					return errors.New("invalid number")
				}
				err = q.IValidation.InsertValidation(ctx, questionID,
					model.Validations{
						MinBound: strconv.Itoa(*b.MinValue),
						MaxBound: strconv.Itoa(*b.MaxValue),
					})
				if err != nil {
					c.Logger().Errorf("failed to insert validation: %+v", err)
					return errors.New("failed to insert validation")
				}
			}
		}

		message := createQuestionnaireMessage(
			questionnaireID,
			params.Title,
			params.Description,
			append(allAdminUsers, adminGroupNames...),
			responseDueDateTime,
			append(allTargetUsers, targetGroupNames...),
		)
		err = q.PostMessage(message)
		if err != nil {
			c.Logger().Errorf("failed to post message: %+v", err)
			return err
		}

		if params.ResponseDueDateTime != nil {
			err = q.PushReminder(questionnaireID, params.ResponseDueDateTime)
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
	// unable to change the questionnaire from anoymous to non-anonymous
	isAnonymous, err := q.GetResponseIsAnonymousByQuestionnaireID(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get anonymous info: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get anonymous info")
	}
	if isAnonymous && !params.IsAnonymous {
		c.Logger().Info("unable to change the questionnaire from anoymous to non-anonymous")
		return echo.NewHTTPError(http.StatusBadRequest, "unable to change the questionnaire from anoymous to non-anonymous")
	}

	responseDueDateTime := null.Time{}
	if params.ResponseDueDateTime != nil {
		responseDueDateTime.Valid = true
		responseDueDateTime.Time = *params.ResponseDueDateTime
	}
	if responseDueDateTime.Valid {
		isBefore := responseDueDateTime.ValueOrZero().Before(time.Now())
		if isBefore {
			c.Logger().Infof("invalid resTimeLimit: %+v", responseDueDateTime)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid resTimeLimit")
		}
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
			questionType := questionParsed["question_type"].(string)
			if questionType == "Text" {
				questionType = "Text"
			} else if questionType == "TextLong" {
				questionType = "TextArea"
			} else if questionType == "Number" {
				questionType = "Number"
			} else if questionType == "SingleChoice" {
				questionType = "MultipleChoice"
			} else if questionType == "MultipleChoice" {
				questionType = "Checkbox"
			} else if questionType == "Scale" {
				questionType = "LinearScale"
			} else {
				c.Logger().Errorf("invalid question type")
				return errors.New("invalid question type")
			}
			if question.QuestionId == nil {
				questionID, err := q.InsertQuestion(ctx, questionnaireID, 1, questoinNum+1, questionType, question.Body, question.IsRequired)
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
					err = q.IScaleLabel.InsertScaleLabel(ctx, questionID,
						model.ScaleLabels{
							ScaleLabelLeft:  *b.MinLabel,
							ScaleLabelRight: *b.MaxLabel,
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
							RegexPattern: "^.{0," + strconv.Itoa(*b.MaxLength) + "}$",
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
							RegexPattern: "^.{0," + fmt.Sprintf("%d", *b.MaxLength) + "}$",
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
					err = q.IValidation.CheckNumberValid(strconv.Itoa(*b.MinValue), strconv.Itoa(*b.MaxValue))
					if err != nil {
						c.Logger().Errorf("invalid number: %+v", err)
						return errors.New("invalid number")
					}
					err = q.IValidation.InsertValidation(ctx, questionID,
						model.Validations{
							MinBound: strconv.Itoa(*b.MinValue),
							MaxBound: strconv.Itoa(*b.MaxValue),
						})
					if err != nil {
						c.Logger().Errorf("failed to insert validation: %+v", err)
						return errors.New("failed to insert validation")
					}
				}
			} else {
				ifQuestionExist[*question.QuestionId] = true
				err = q.UpdateQuestion(ctx, questionnaireID, 1, questoinNum+1, questionType, question.Body, question.IsRequired, *question.QuestionId)
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
					err = q.IScaleLabel.UpdateScaleLabel(ctx, *question.QuestionId,
						model.ScaleLabels{
							ScaleLabelLeft:  *b.MinLabel,
							ScaleLabelRight: *b.MaxLabel,
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
							RegexPattern: "^.{0," + strconv.Itoa(*b.MaxLength) + "}$",
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
							RegexPattern: "^.{0," + strconv.Itoa(*b.MaxLength) + "}$",
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
					err = q.IValidation.CheckNumberValid(strconv.Itoa(*b.MinValue), strconv.Itoa(*b.MaxValue))
					if err != nil {
						c.Logger().Errorf("invalid number: %+v", err)
						return errors.New("invalid number")
					}
					err = q.IValidation.UpdateValidation(ctx, *question.QuestionId,
						model.Validations{
							MinBound: strconv.Itoa(*b.MinValue),
							MaxBound: strconv.Itoa(*b.MaxValue),
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
		if params.ResponseDueDateTime != nil {
			err = q.PushReminder(questionnaireID, params.ResponseDueDateTime)
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
		err := q.IQuestionnaire.DeleteQuestionnaire(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete questionnaire: %+v", err)
			return err
		}

		err = q.DeleteTargets(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete targets: %+v", err)
			return err
		}

		err = q.DeleteAdministrators(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrators: %+v", err)
			return err
		}

		questions, err := q.GetQuestions(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to get questions: %+v", err)
			return err
		}
		for _, question := range questions {
			err = q.DeleteQuestion(ctx, question.ID)
			if err != nil {
				c.Logger().Errorf("failed to delete administrators: %+v", err)
				return err
			}
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
	status, err := q.GetTargetsCancelStatus(c.Request().Context(), questionnaireID, []string{userID})
	if err != nil {
		c.Logger().Errorf("failed to check remind status: %+v", err)
		return false, echo.NewHTTPError(http.StatusInternalServerError, "failed to check remind status")
	}

	return !status[0].IsCanceled, nil
}

func (q *Questionnaire) EditQuestionnaireMyRemindStatus(c echo.Context, questionnaireID int, userID string, isRemindEnabled bool) error {
	err := q.UpdateTargetsCancelStatus(c.Request().Context(), questionnaireID, []string{userID}, !isRemindEnabled)
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
	respondentDetails, err := q.GetRespondentDetails(c.Request().Context(), questionnaireID, sort, onlyMyResponse, userID)
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
		return res, echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	questions, err := q.IQuestion.GetQuestions(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get questions: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	responseMetas, err := responseBody2ResponseMetas(params.Body, questions)
	if err != nil {
		c.Logger().Errorf("failed to convert response body to response metas: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
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
		case "Number":
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
		case "Checkbox", "MultipleChoice":
			option, ok := optionMap[responseMeta.QuestionID]
			if !ok {
				option = []model.Options{}
			}
			var selectedOptions []int
			if questionTypes[responseMeta.QuestionID] == "MultipleChoice" {
				var selectedOption int
				err = json.Unmarshal([]byte(responseMeta.Data), &selectedOption)
				if err != nil {
					c.Logger().Errorf("invalid option: %+v", err)
					return res, echo.NewHTTPError(http.StatusBadRequest, err)
				}
				selectedOptions = append(selectedOptions, selectedOption)
			} else if questionTypes[responseMeta.QuestionID] == "Checkbox" {
				err = json.Unmarshal([]byte(responseMeta.Data), &selectedOptions)
				if err != nil {
					c.Logger().Errorf("invalid option: %+v", err)
					return res, echo.NewHTTPError(http.StatusBadRequest, err)
				}
			}
			ok = true
			if len(selectedOptions) == 0 {
				ok = false
			}
			sort.Slice(selectedOptions, func(i, j int) bool { return selectedOptions[i] < selectedOptions[j] })
			var preOption *int
			for _, selectedOption := range selectedOptions {
				if preOption != nil && *preOption == selectedOption {
					ok = false
					break
				}
				if selectedOption < 1 || selectedOption > len(option) {
					ok = false
					break
				}
				preOption = &selectedOption
			}
			if !ok {
				c.Logger().Errorf("invalid option: %+v", err)
				return res, echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "LinearScale":
			label, ok := scaleLabelMap[responseMeta.QuestionID]
			if !ok {
				label = model.ScaleLabels{}
			}
			err := q.IScaleLabel.CheckScaleLabel(label, responseMeta.Data)
			if err != nil {
				c.Logger().Errorf("invalid scale: %+v", err)
				return res, echo.NewHTTPError(http.StatusBadRequest, err)
			}
		default:
			c.Logger().Errorf("invalid question id: %+v", responseMeta.QuestionID)
			return res, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid question id: %d", responseMeta.QuestionID))
		}
	}

	for _, question := range questions {
		if questionRequired[question.ID] {
			c.Logger().Errorf("required question is not answered: %+v", question.ID)
			return res, echo.NewHTTPError(http.StatusBadRequest, "required question is not answered")
		}
	}

	var submittedAt, modifiedAt time.Time
	//一時保存のときはnull
	if params.IsDraft {
		submittedAt = time.Time{}
		modifiedAt = time.Time{}
	} else {
		submittedAt = time.Now()
		modifiedAt = time.Now()
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

	isAnonymous, err := q.GetResponseIsAnonymousByQuestionnaireID(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get response isanonymous by questionnaire id: %+v", err)
		return res, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res = openapi.Response{
		QuestionnaireId: questionnaireID,
		ResponseId:      responseID,
		Respondent:      &userID,
		SubmittedAt:     submittedAt,
		ModifiedAt:      modifiedAt,
		IsDraft:         params.IsDraft,
		IsAnonymous:     &isAnonymous,
		Body:            params.Body,
	}

	return res, nil
}

func createQuestionnaireMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit null.Time, targets []string) string {
	var resTimeLimitText string
	if resTimeLimit.Valid {
		resTimeLimitText = resTimeLimit.Time.Local().Format("2006/01/02 15:04")
	} else {
		resTimeLimitText = "なし"
	}

	var targetsMentionText string
	if len(targets) == 0 {
		targetsMentionText = "なし"
	} else {
		targetsMentionText = "@" + strings.Join(targets, " @")
	}

	return fmt.Sprintf(
		`### アンケート『[%s](https://anke-to.trap.jp/questionnaires/%d)』が作成されました
#### 管理者
%s
#### 説明
%s
#### 回答期限
%s
#### 対象者
%s
#### 回答リンク
https://anke-to.trap.jp/responses/new/%d`,
		title,
		questionnaireID,
		strings.Join(administrators, ","),
		description,
		resTimeLimitText,
		targetsMentionText,
		questionnaireID,
	)
}

func createReminderMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit time.Time, targets []string, leftTimeText string) string {
	resTimeLimitText := resTimeLimit.Local().Format("2006/01/02 15:04")
	targetsMentionText := "@" + strings.Join(targets, " @")

	return fmt.Sprintf(
		`### アンケート『[%s](https://anke-to.trap.jp/questionnaires/%d)』の回答期限が迫っています!
==残り%sです!==
#### 管理者
%s
#### 説明
%s
#### 回答期限
%s
#### 対象者
%s
#### 回答リンク
https://anke-to.trap.jp/responses/new/%d`,
		title,
		questionnaireID,
		leftTimeText,
		strings.Join(administrators, ","),
		description,
		resTimeLimitText,
		targetsMentionText,
		questionnaireID,
	)
}
