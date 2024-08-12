package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	model.IAdministrator
	model.IAdministratorGroup
	model.IQuestion
	model.IOption
	model.IScaleLabel
	model.IValidation
	model.ITransaction
	traq.IWebhook
	Response
}

func NewQuestionnaire() *Questionnaire {
	return &Questionnaire{}
}

const MaxTitleLength = 50

func (q Questionnaire) GetQuestionnaires(ctx echo.Context, userID string, params openapi.GetQuestionnairesParams) (openapi.QuestionnaireList, error) {
	res := openapi.QuestionnaireList{}
	sort := string(*params.Sort)
	search := string(*params.Search)
	pageNum := int(*params.Page)
	if pageNum < 1 {
		pageNum = 1
	}

	questionnaireList, pageMax, err := q.IQuestionnaire.GetQuestionnaires(ctx.Request().Context(), userID, sort, search, pageNum, *params.OnlyTargetingMe, *params.OnlyAdministratedByMe)
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

func (q Questionnaire) PostQuestionnaire(c echo.Context, userID string, params openapi.PostQuestionnaireJSONRequestBody) (openapi.QuestionnaireDetail, error) {
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

	err := q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		questionnaireID, err := q.InsertQuestionnaire(ctx, params.Title, params.Description, responseDueDateTime, convertResponseViewableBy(params.ResponseViewableBy))
		if err != nil {
			c.Logger().Errorf("failed to insert questionnaire: %+v", err)
			return err
		}
		allTargetUsers, err := rollOutUsersAndGroups(params.Targets.Users, params.Targets.Groups)
		if err != nil {
			c.Logger().Errorf("failed to roll out users and groups: %+v", err)
			return err
		}
		targetGroupNames, err := uuid2GroupNames(params.Targets.Groups)
		if err != nil {
			c.Logger().Errorf("failed to get group names: %+v", err)
			return err
		}
		err = q.InsertTargets(ctx, questionnaireID, allTargetUsers)
		if err != nil {
			c.Logger().Errorf("failed to insert targets: %+v", err)
			return err
		}
		err = q.InsertTargetGroups(ctx, questionnaireID, params.Targets.Groups)
		if err != nil {
			c.Logger().Errorf("failed to insert target groups: %+v", err)
			return err
		}
		allAdminUsers, err := rollOutUsersAndGroups(params.Admins.Users, params.Admins.Groups)
		if err != nil {
			c.Logger().Errorf("failed to roll out administrators: %+v", err)
			return err
		}
		adminGroupNames, err := uuid2GroupNames(params.Admins.Groups)
		if err != nil {
			c.Logger().Errorf("failed to get group names: %+v", err)
			return err
		}
		err = q.InsertAdministrators(ctx, questionnaireID, allAdminUsers)
		if err != nil {
			c.Logger().Errorf("failed to insert administrators: %+v", err)
			return err
		}
		err = q.InsertAdministratorGroups(ctx, questionnaireID, params.Admins.Groups)
		if err != nil {
			c.Logger().Errorf("failed to insert administrator groups: %+v", err)
			return err
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

		return nil
	})
	if err != nil {
		c.Logger().Errorf("failed to create a questionnaire: %+v", err)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to create a questionnaire")
	}
	questionnaireInfo, targets, targetGroups, admins, adminGroups, respondents, err := q.GetQuestionnaireInfo(c.Request().Context(), questionnaireID)
	if err != nil {
		c.Logger().Errorf("failed to get questionnaire info: %+v", err)
		return openapi.QuestionnaireDetail{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to get questionnaire info")
	}

	questionnaireDetail := questionnaire2QuestionnaireDetail(*questionnaireInfo, admins, adminGroups, targets, targetGroups, respondents)
	return questionnaireDetail, nil
}
func (q Questionnaire) GetQuestionnaire(ctx echo.Context, questionnaireID int) (openapi.QuestionnaireDetail, error) {
	questionnaireInfo, targets, targetGroups, admins, adminGroups, respondents, err := q.GetQuestionnaireInfo(ctx.Request().Context(), questionnaireID)
	if err != nil {
		return openapi.QuestionnaireDetail{}, err
	}
	questionnaireDetail := questionnaire2QuestionnaireDetail(*questionnaireInfo, admins, adminGroups, targets, targetGroups, respondents)
	return questionnaireDetail, nil
}

func (q Questionnaire) EditQuestionnaire(c echo.Context, questionnaireID int, params openapi.EditQuestionnaireJSONRequestBody) error {
	responseDueDateTime := null.Time{}
	if params.ResponseDueDateTime != nil {
		responseDueDateTime.Valid = true
		responseDueDateTime.Time = *params.ResponseDueDateTime
	}
	err := q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		err := q.UpdateQuestionnaire(ctx, params.Title, params.Description, responseDueDateTime, string(params.ResponseViewableBy), questionnaireID)
		if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			c.Logger().Errorf("failed to update questionnaire: %+v", err)
			return err
		}
		err = q.DeleteTargets(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete targets: %+v", err)
			return err
		}
		err = q.DeleteTargetGroups(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete target groups: %+v", err)
			return err
		}
		allTargetUsers, err := rollOutUsersAndGroups(params.Targets.Users, params.Targets.Groups)
		if err != nil {
			c.Logger().Errorf("failed to roll out users and groups: %+v", err)
			return err
		}
		err = q.InsertTargets(ctx, questionnaireID, allTargetUsers)
		if err != nil {
			c.Logger().Errorf("failed to insert targets: %+v", err)
			return err
		}
		err = q.InsertTargetGroups(ctx, questionnaireID, params.Targets.Groups)
		if err != nil {
			c.Logger().Errorf("failed to insert target groups: %+v", err)
			return err
		}
		err = q.DeleteAdministrators(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrators: %+v", err)
			return err
		}
		err = q.DeleteAdministratorGroups(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrator groups: %+v", err)
			return err
		}
		allAdminUsers, err := rollOutUsersAndGroups(params.Admins.Users, params.Admins.Groups)
		if err != nil {
			c.Logger().Errorf("failed to roll out administrators: %+v", err)
			return err
		}
		err = q.InsertAdministrators(ctx, questionnaireID, allAdminUsers)
		if err != nil {
			c.Logger().Errorf("failed to insert administrators: %+v", err)
			return err
		}
		err = q.InsertAdministratorGroups(ctx, questionnaireID, params.Admins.Groups)
		if err != nil {
			c.Logger().Errorf("failed to insert administrator groups: %+v", err)
			return err
		}

		return nil
	})
	if err != nil {
		c.Logger().Errorf("failed to update a questionnaire: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update a questionnaire")
	}

	return nil
}

func (q Questionnaire) DeleteQuestionnaire(c echo.Context, questionnaireID int) error {
	err := q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		err := q.IQuestionnaire.DeleteQuestionnaire(c.Request().Context(), questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete questionnaire: %+v", err)
			return err
		}

		err = q.DeleteTargets(c.Request().Context(), questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete targets: %+v", err)
			return err
		}

		err = q.DeleteAdministrators(c.Request().Context(), questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrators: %+v", err)
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

func (q Questionnaire) GetQuestionnaireMyRemindStatus(c echo.Context, questionnaireID int) (bool, error) {
	// todo: check remind status
	return false, nil
}

func (q Questionnaire) EditQuestionnaireMyRemindStatus(c echo.Context, questionnaireID int) error {
	// todo: edit remind status
	return nil
}

func (q Questionnaire) GetQuestionnaireResponses(c echo.Context, questionnaireID int, params openapi.GetQuestionnaireResponsesParams, userID string) (openapi.Responses, error) {
	res := []openapi.Response{}
	respondentDetails, err := q.GetRespondentDetails(c.Request().Context(), questionnaireID, string(*params.Sort), *params.OnlyMyResponse, userID)
	if err != nil {
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

func (q Questionnaire) PostQuestionnaireResponse(c echo.Context) error {
	// todo: PostQuestionnaireResponse
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