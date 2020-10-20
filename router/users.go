package router

import (
	"net/http"
	"sort"
	"strconv"
	"time"

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

	myResponses, err := model.GetRespondentInfos(c, userID)
	if err != nil {
		return err
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

	myresponses, err := model.GetRespondentInfos(c, userID, questionnaireID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myresponses)
}

// GetTargetedQuestionnaire GET /users/me/targeted
func GetTargetedQuestionnaire(c echo.Context) error {
	userID := model.GetUserID(c)

	ret, err := model.GetTargettedQuestionnaires(c, userID, "")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

// GetMyQuestionnaire GET /users/me/administrates
func GetMyQuestionnaire(c echo.Context) error {
	// 自分が管理者になっているアンケート一覧
	questionnaireIDs, err := model.GetAdminQuestionnaireIDs(c, model.GetUserID(c))
	if err != nil {
		return nil
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

	for _, questionnaireID := range questionnaireIDs {
		questionnaire, targets, administrators, respondents, err := model.GetQuestionnaireInfo(c, questionnaireID)
		if err != nil {
			return err
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

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ModifiedAt > ret[j].ModifiedAt
	})

	return c.JSON(http.StatusOK, ret)
}

// GetTargettedQuestionnairesBytraQID GET /users/:traQID/targeted
func GetTargettedQuestionnairesBytraQID(c echo.Context) error {
	traQID := c.Param("traQID")

	ret, err := model.GetTargettedQuestionnaires(c, traQID, "unanswered")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
