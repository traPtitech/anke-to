package router

import (
	"net/http"
	"strconv"
	"time"

	"database/sql"
	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

func GetQuestionnaires(c echo.Context) error {
	if c.QueryParam("nontargeted") == "true" {
		questionnaires, err := model.GetQuestionnaires(c, model.TargetType(model.Nontargeted))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, questionnaires)
	} else {
		questionnaires, err := model.GetQuestionnaires(c, model.TargetType(model.All))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, questionnaires)
	}
}

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

func PostQuestionnaire(c echo.Context) error {

	// リクエストで投げられるJSONのスキーマ
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

	var result sql.Result

	// アンケートの追加
	if req.ResTimeLimit == "" || req.ResTimeLimit == "NULL" {
		req.ResTimeLimit = "NULL"
		var err error
		result, err = model.DB.Exec(
			"INSERT INTO questionnaires (title, description, res_shared_to) VALUES (?, ?, ?)",
			req.Title, req.Description, req.ResSharedTo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		var err error
		result, err = model.DB.Exec(
			"INSERT INTO questionnaires (title, description, res_time_limit, res_shared_to) VALUES (?, ?, ?, ?)",
			req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := model.InsertTargets(c, int(lastID), req.Targets); err != nil {
		return err
	}

	if err := model.InsertAdministrators(c, int(lastID), req.Administrators); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": int(lastID),
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

	// アップデートする
	if req.ResTimeLimit == "" || req.ResTimeLimit == "NULL" {
		req.ResTimeLimit = "NULL"
		if _, err := model.DB.Exec(
			"UPDATE questionnaires SET title = ?, description = ?, res_time_limit = NULL, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?",
			req.Title, req.Description, req.ResSharedTo, questionnaireID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := model.DB.Exec(
			"UPDATE questionnaires SET title = ?, description = ?, res_time_limit = ?, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?",
			req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo, questionnaireID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
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

func DeleteQuestionnaire(c echo.Context) error {

	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if _, err := model.DB.Exec(
		"UPDATE questionnaires SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := model.DeleteTargets(c, questionnaireID); err != nil {
		return err
	}

	if err := model.DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func GetMyQuestionnaire(c echo.Context) error {
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
	return c.JSON(http.StatusOK, ret)
}

func GetTargetedQuestionnaire(c echo.Context) error {
	questionnaires, err := model.GetQuestionnaires(c, model.TargetType(model.Targeted))
	if err != nil {
		return err
	}

	type QuestionnairesInfo struct {
		ID           int    `json:"questionnaireID"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		ResTimeLimit string `json:"res_time_limit"`
		ResSharedTo  string `json:"res_shared_to"`
		CreatedAt    string `json:"created_at"`
		ModifiedAt   string `json:"modified_at"`
		RespondedAt  string `json:"responded_at"`
	}
	questionnairesInfo := []QuestionnairesInfo{}

	for _, q := range questionnaires {
		respondedAt := sql.NullString{}
		if err := model.DB.Get(&respondedAt,
			`SELECT MAX(submitted_at) FROM respondents
			WHERE user_traqid = ? AND questionnaire_id = ? AND deleted_at IS NULL`,
			model.GetUserID(c), q.ID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		questionnairesInfo = append(questionnairesInfo,
			QuestionnairesInfo{
				ID:           q.ID,
				Title:        q.Title,
				Description:  q.Description,
				ResTimeLimit: q.ResTimeLimit,
				ResSharedTo:  q.ResSharedTo,
				CreatedAt:    q.CreatedAt,
				ModifiedAt:   q.ModifiedAt,
				RespondedAt:  model.NullStringConvert(respondedAt),
			})
	}
	return c.JSON(http.StatusOK, questionnairesInfo)
}
