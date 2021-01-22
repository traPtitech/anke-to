package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/traPtitech/anke-to/model"
)

// GetUsersMe GET /users/me
func GetUsersMe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"traqID": model.GetUserID(c),
	})
}

// GetMyResponses GET /users/me/responses
func GetMyResponses(c echo.Context) error {
	userID := model.GetUserID(c)

	myResponses, err := model.GetRespondentInfos(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, myResponses)
}

// GetMyResponsesByID GET /users/me/responses/:questionnaireID
func GetMyResponsesByID(c echo.Context) error {
	userID := model.GetUserID(c)

	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	myresponses, err := model.GetRespondentInfos(userID, questionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, myresponses)
}

// GetTargetedQuestionnaire GET /users/me/targeted
func GetTargetedQuestionnaire(c echo.Context) error {
	userID := model.GetUserID(c)
	sort := c.QueryParam("sort")
	ret, err := model.GetTargettedQuestionnaires(userID, "", sort)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ret)
}

// GetMyQuestionnaire GET /users/me/administrates
func GetMyQuestionnaire(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	// 自分が管理者になっているアンケート一覧
	questionnaires, err := model.GetAdminQuestionnaires(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaires: %w", err))
	}

	questionnaireIDs := make([]int, 0, len(questionnaires))
	for _, questionnaire := range questionnaires {
		questionnaireIDs = append(questionnaireIDs, questionnaire.ID)
	}

	targets, err := model.GetTargets(questionnaireIDs)
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

	respondents, err := model.GetRespondentsUserIDs(questionnaireIDs)
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

	administrators, err := model.GetAdministrators(questionnaireIDs)
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
		ID             int      `json:"questionnaireID"`
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		ResTimeLimit   string   `json:"res_time_limit"`
		CreatedAt      string   `json:"created_at"`
		ModifiedAt     string   `json:"modified_at"`
		ResSharedTo    string   `json:"res_shared_to"`
		AllResponded   bool     `json:"all_responded"`
		Targets        []string `json:"targets"`
		Administrators []string `json:"administrators"`
		Respondents    []string `json:"respondents"`
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
			ResTimeLimit:   model.NullTimeToString(questionnaire.ResTimeLimit),
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
func GetTargettedQuestionnairesBytraQID(c echo.Context) error {
	traQID := c.Param("traQID")
	sort := c.QueryParam("sort")
	ret, err := model.GetTargettedQuestionnaires(traQID, "unanswered", sort)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ret)
}
