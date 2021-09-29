package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v3"

	"github.com/traPtitech/anke-to/model"
)

// User Userの構造体
type User struct {
	model.IRespondent
	model.IQuestionnaire
	model.ITarget
	model.IAdministrator
}

type UserQueryparam struct {
	Sort string `validate:"omitempty,oneof=created_at -created_at title -title modified_at -modified_at"`
	Answered string `validate:"omitempty,oneof=answered unanswered"`
	TraQID string `validate:"required"`
}

// NewUser Userのコンストラクタ
func NewUser(respondent model.IRespondent, questionnaire model.IQuestionnaire, target model.ITarget, administrator model.IAdministrator) *User {
	return &User{
		IRespondent:    respondent,
		IQuestionnaire: questionnaire,
		ITarget:        target,
		IAdministrator: administrator,
	}
}

// GetUsersMe GET /users/me
func (*User) GetUsersMe(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"traqID": userID,
	})
}

// GetMyResponses GET /users/me/responses
func (u *User) GetMyResponses(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	myResponses, err := u.GetRespondentInfos(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, myResponses)
}

// GetMyResponsesByID GET /users/me/responses/:questionnaireID
func (u *User) GetMyResponsesByID(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	myresponses, err := u.GetRespondentInfos(c.Request().Context(), userID, questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, myresponses)
}

// GetTargetedQuestionnaire GET /users/me/targeted
func (u *User) GetTargetedQuestionnaire(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	sort := c.QueryParam("sort")
	ret, err := u.GetTargettedQuestionnaires(c.Request().Context(), userID, "", sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ret)
}

// GetMyQuestionnaire GET /users/me/administrates
func (u *User) GetMyQuestionnaire(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	// 自分が管理者になっているアンケート一覧
	questionnaires, err := u.GetAdminQuestionnaires(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaires: %w", err))
	}

	questionnaireIDs := make([]int, 0, len(questionnaires))
	for _, questionnaire := range questionnaires {
		questionnaireIDs = append(questionnaireIDs, questionnaire.ID)
	}

	targets, err := u.GetTargets(questionnaireIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get targets: %w", err))
	}
	targetMap := map[int][]string{}
	for _, target := range targets {
		tgts, ok := targetMap[target.QuestionnaireID]
		if !ok {
			targetMap[target.QuestionnaireID] = []string{target.UserTraqid}
		} else {
			targetMap[target.QuestionnaireID] = append(tgts, target.UserTraqid)
		}
	}

	respondents, err := u.GetRespondentsUserIDs(c.Request().Context(), questionnaireIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondents: %w", err))
	}
	respondentMap := map[int][]string{}
	for _, respondent := range respondents {
		rspdts, ok := respondentMap[respondent.QuestionnaireID]
		if !ok {
			respondentMap[respondent.QuestionnaireID] = []string{respondent.UserTraqid}
		} else {
			respondentMap[respondent.QuestionnaireID] = append(rspdts, respondent.UserTraqid)
		}
	}

	administrators, err := u.GetAdministrators(c.Request().Context(), questionnaireIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get administrators: %w", err))
	}
	administratorMap := map[int][]string{}
	for _, administrator := range administrators {
		admins, ok := administratorMap[administrator.QuestionnaireID]
		if !ok {
			administratorMap[administrator.QuestionnaireID] = []string{administrator.UserTraqid}
		} else {
			administratorMap[administrator.QuestionnaireID] = append(admins, administrator.UserTraqid)
		}
	}

	type QuestionnaireInfo struct {
		ID             int       `json:"questionnaireID"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		ResTimeLimit   null.Time `json:"res_time_limit"`
		CreatedAt      string    `json:"created_at"`
		ModifiedAt     string    `json:"modified_at"`
		ResSharedTo    string    `json:"res_shared_to"`
		AllResponded   bool      `json:"all_responded"`
		Targets        []string  `json:"targets"`
		Administrators []string  `json:"administrators"`
		Respondents    []string  `json:"respondents"`
	}
	ret := []QuestionnaireInfo{}

	for _, questionnaire := range questionnaires {
		targets, ok := targetMap[questionnaire.ID]
		if !ok {
			targets = []string{}
		}

		administrators, ok := administratorMap[questionnaire.ID]
		if !ok {
			administrators = []string{}
		}

		respondents, ok := respondentMap[questionnaire.ID]
		if !ok {
			respondents = []string{}
		}

		allresponded := true
		for _, t := range targets {
			found := false
			for _, r := range respondents {
				if t == r {
					found = true
					break
				}
			}
			if !found {
				allresponded = false
				break
			}
		}

		ret = append(ret, QuestionnaireInfo{
			ID:             questionnaire.ID,
			Title:          questionnaire.Title,
			Description:    questionnaire.Description,
			ResTimeLimit:   questionnaire.ResTimeLimit,
			CreatedAt:      questionnaire.CreatedAt.Format(time.RFC3339),
			ModifiedAt:     questionnaire.ModifiedAt.Format(time.RFC3339),
			ResSharedTo:    questionnaire.ResSharedTo,
			AllResponded:   allresponded,
			Targets:        targets,
			Administrators: administrators,
			Respondents:    respondents,
		})
	}

	return c.JSON(http.StatusOK, ret)
}

// GetTargettedQuestionnairesBytraQID GET /users/:traQID/targeted
func (u *User) GetTargettedQuestionnairesBytraQID(c echo.Context) error {
	traQID := c.Param("traQID")
	sort := c.QueryParam("sort")
	answered := c.QueryParam("answered")

	p := UserQueryparam{
		Sort:     sort,
		Answered: answered,
		TraQID:   traQID,
	}

	validate,err := getValidator(c)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get validator:%w",err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(),p)
	if err != nil {
		c.Logger().Info(fmt.Errorf("failed to validate:%w",err))
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}

	ret, err := u.GetTargettedQuestionnaires(c.Request().Context(), traQID, answered, sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ret)
}
