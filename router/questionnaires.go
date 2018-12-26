package router

import (
	"net/http"
	"strconv"
	"time"

	"database/sql"
	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

// echo.Contextを引数にとってerrorを返り値とする
func GetQuestionnaires(c echo.Context, targettype model.TargetType) error {
	allquestionnaires, err := model.GetAllQuestionnaires(c)
	if err != nil {
		return err
	}

	userID := model.GetUserID(c)

	targetedQuestionnaireID := []int{}
	if err := model.DB.Select(&targetedQuestionnaireID,
		"SELECT questionnaire_id FROM targets WHERE user_traqid = ?", userID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	type questionnairesInfo struct {
		ID           int       `json:"questionnaireID"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		ResTimeLimit string    `json:"res_time_limit"`
		ResSharedTo  string    `json:"res_shared_to"`
		CreatedAt    time.Time `json:"created_at"`
		ModifiedAt   time.Time `json:"modified_at"`
		IsTargeted   bool      `json:"is_targeted"`
	}
	var ret []questionnairesInfo

	for _, v := range allquestionnaires {
		var targeted = false
		for _, w := range targetedQuestionnaireID {
			if w == v.ID {
				targeted = true
			}
		}
		if (targettype == model.TargetType(model.Targeted) && !targeted) || (targettype == model.TargetType(model.Nontargeted) && targeted) {
			continue
		}
		ret = append(ret,
			questionnairesInfo{
				ID:           v.ID,
				Title:        v.Title,
				Description:  v.Description,
				ResTimeLimit: model.TimeConvert(v.ResTimeLimit),
				ResSharedTo:  v.ResSharedTo,
				CreatedAt:    v.CreatedAt,
				ModifiedAt:   v.ModifiedAt,
				IsTargeted:   targeted})
	}

	if len(ret) == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// 構造体の定義で書いたJSONのキーで変換される
	return c.JSON(http.StatusOK, ret)
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
