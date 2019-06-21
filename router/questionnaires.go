package router

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

// GetQuestionnaires GET /questionnaires
func GetQuestionnaires(c echo.Context) error {

	questionnaires, pageMax, err := model.GetQuestionnaires(c, c.QueryParam("nontargeted") == "true")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"page_max":       pageMax,
		"questionnaires": questionnaires,
	})
}

// GetQuestionnaire GET /questionnaires/:id
func GetQuestionnaire(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	questionnaire, targets, administrators, respondents, err := model.GetQuestionnaireInfo(c, questionnaireID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionnaireID": questionnaire.ID,
		"title":           questionnaire.Title,
		"description":     questionnaire.Description,
		"res_time_limit":  model.NullTimeToString(questionnaire.ResTimeLimit),
		"created_at":      questionnaire.CreatedAt.Format(time.RFC3339),
		"modified_at":     questionnaire.ModifiedAt.Format(time.RFC3339),
		"res_shared_to":   questionnaire.ResSharedTo,
		"targets":         targets,
		"administrators":  administrators,
		"respondents":     respondents,
	})
}

// PostQuestionnaire POST /questionnaires
func PostQuestionnaire(c echo.Context) error {

	req := struct {
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		ResTimeLimit   string   `json:"res_time_limit"`
		ResSharedTo    string   `json:"res_shared_to"`
		Targets        []string `json:"targets"`
		Administrators []string `json:"administrators"`
	}{}

	// JSONを構造体につける
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	lastID, err := model.InsertQuestionnaire(c, req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo)
	if err != nil {
		return err
	}

	if err := model.InsertTargets(c, lastID, req.Targets); err != nil {
		return err
	}

	if err := model.InsertAdministrators(c, lastID, req.Administrators); err != nil {
		return err
	}

	time_limit := "なし"
	if req.ResTimeLimit != "NULL" {
		time_limit = req.ResTimeLimit
	}

	if err := model.PostMessage(c,
		"### 新しいアンケートが作成されました\n"+
			"#### タイトル\n"+req.Title+"\n"+
			"#### 管理者\n"+strings.Join(req.Administrators, ",")+"\n"+
			"#### 説明\n"+req.Description+"\n"+
			"#### 回答期限\n"+time_limit+"\n"+
			"http://anke-to.sysad.trap.show/responses/new/"+strconv.Itoa(lastID)); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": lastID,
		"title":           req.Title,
		"description":     req.Description,
		"res_time_limit":  req.ResTimeLimit,
		"deleted_at":      "NULL",
		"created_at":      time.Now().Format(time.RFC3339),
		"modified_at":     time.Now().Format(time.RFC3339),
		"res_shared_to":   req.ResSharedTo,
		"targets":         req.Targets,
		"administrators":  req.Administrators,
	})
}

// EditQuestionnaire PATCH /questonnaires/:id
func EditQuestionnaire(c echo.Context) error {

	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	req := struct {
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		ResTimeLimit   string   `json:"res_time_limit"`
		ResSharedTo    string   `json:"res_shared_to"`
		Targets        []string `json:"targets"`
		Administrators []string `json:"administrators"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if req.ResSharedTo == "" {
		req.ResSharedTo = "administrators"
	}

	if err := model.UpdateQuestionnaire(
		c, req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo, questionnaireID); err != nil {
		return err
	}

	if err := model.DeleteTargets(c, questionnaireID); err != nil {
		return err
	}

	if err := model.InsertTargets(c, questionnaireID, req.Targets); err != nil {
		return err
	}

	if err := model.DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	if err := model.InsertAdministrators(c, questionnaireID, req.Administrators); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestionnaire DELETE /questonnaires/:id
func DeleteQuestionnaire(c echo.Context) error {

	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := model.DeleteQuestionnaire(c, questionnaireID); err != nil {
		return err
	}

	if err := model.DeleteTargets(c, questionnaireID); err != nil {
		return err
	}

	if err := model.DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// GetMyQuestionnaire GET /users/me/administrates
func GetMyQuestionnaire(c echo.Context) error {
	// 自分が管理者になっているアンケート一覧
	questionnaireIDs, err := model.GetAdminQuestionnaires(c, model.GetUserID(c))
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

// GetTargetedQuestionnaire GET /users/me/targeted
func GetTargetedQuestionnaire(c echo.Context) error {
	ret, err := model.GetTargettedQuestionnaires(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
