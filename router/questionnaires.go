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
		return model.GetQuestionnaires(c, model.TargetType(model.Nontargeted))
	} else {
		return model.GetQuestionnaires(c, model.TargetType(model.All))
	}
}

func GetQuestionnaire(c echo.Context) error {

	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	questionnaire := model.Questionnaires{}
	if err := model.DB.Get(&questionnaire, "SELECT * FROM questionnaires WHERE id = ? AND deleted_at IS NULL", questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	targets, err := model.GetTargets(c, questionnaireID)
	if err != nil {
		return err
	}

	administrators, err := model.GetAdministrators(c, questionnaireID)
	if err != nil {
		return err
	}

	respondents := []string{}
	if err := model.DB.Select(&respondents, "SELECT user_traqid FROM respondents WHERE questionnaire_id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionnaireID": questionnaire.ID,
		"title":           questionnaire.Title,
		"description":     questionnaire.Description,
		"res_time_limit":  model.TimeConvert(questionnaire.ResTimeLimit),
		"created_at":      questionnaire.CreatedAt,
		"modified_at":     questionnaire.ModifiedAt,
		"res_shared_to":   questionnaire.ResSharedTo,
		"targets":         targets,
		"administrators":  administrators,
		"respondents":     respondents,
	})
}

func PostQuestionnaire(c echo.Context) error {

	// リクエストで投げられるJSONのスキーマ
	req := struct {
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		ResTimeLimit   time.Time `json:"res_time_limit"`
		ResSharedTo    string    `json:"res_shared_to"`
		Targets        []string  `json:"targets"`
		Administrators []string  `json:"administrators"`
	}{}

	// JSONを構造体につける
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var result sql.Result

	// アンケートの追加
	if req.ResTimeLimit.IsZero() {
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
		"created_at":      time.Now(),
		"modified_at":     time.Now(),
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
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		ResTimeLimit   time.Time `json:"res_time_limit"`
		ResSharedTo    string    `json:"res_shared_to"`
		Targets        []string  `json:"targets"`
		Administrators []string  `json:"administrators"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if req.ResSharedTo == "" {
		req.ResSharedTo = "administrators"
	}

	// アップデートする
	if req.ResTimeLimit.IsZero() {
		if _, err := model.DB.Exec(
			"UPDATE questionnaires SET title = ?, description = ?, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?",
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
	/*
		後で書く
		questionnaireID, err := GetAdminQuestionnaires(c, getUserID(c))
		if err != nil {
			return nil
		}*/
	return c.NoContent(http.StatusOK)
}

func GetTargetedQuestionnaire(c echo.Context) error {
	return model.GetQuestionnaires(c, model.TargetType(model.Targeted))
}
