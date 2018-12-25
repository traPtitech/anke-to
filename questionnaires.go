package main

import (
	"net/http"
	"strconv"
	"time"

	"database/sql"
	"github.com/labstack/echo"
)

// エラーが起きれば(nil, err)
// 起こらなければ(allquestions, nil)を返す
func getAllQuestionnaires(c echo.Context) ([]questionnaires, error) {
	// query parametar
	sort := c.QueryParam("sort")
	page := c.QueryParam("page")

	if page == "" {
		page = "1"
	}
	num, err := strconv.Atoi(page)
	if err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}

	var list = map[string]string{
		"":             "",
		"created_at":   "ORDER BY created_at",
		"-created_at":  "ORDER BY created_at DESC",
		"title":        "ORDER BY title",
		"-title":       "ORDER BY title DESC",
		"modified_at":  "ORDER BY modified_at",
		"-modified_at": "ORDER BY modified_at DESC",
	}
	// アンケート一覧の配列
	allquestionnaires := []questionnaires{}

	if err := db.Select(&allquestionnaires,
		"SELECT * FROM questionnaires WHERE deleted_at IS NULL "+list[sort]+" lIMIT 20 OFFSET "+strconv.Itoa(20*(num-1))); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return allquestionnaires, nil
}

// echo.Contextを引数にとってerrorを返り値とする
func getQuestionnaires(c echo.Context, targettype TargetType) error {
	allquestionnaires, err := getAllQuestionnaires(c)
	if err != nil {
		return err
	}

	userID := getUserID(c)

	targetedQuestionnaireID := []int{}
	if err := db.Select(&targetedQuestionnaireID,
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
		if (targettype == TargetType(Targeted) && !targeted) || (targettype == TargetType(Nontargeted) && targeted) {
			continue
		}
		ret = append(ret,
			questionnairesInfo{
				ID:           v.ID,
				Title:        v.Title,
				Description:  v.Description,
				ResTimeLimit: timeConvert(v.ResTimeLimit),
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

func getQuestionnaire(c echo.Context) error {

	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	questionnaire := questionnaires{}
	if err := db.Get(&questionnaire, "SELECT * FROM questionnaires WHERE id = ? AND deleted_at IS NULL", questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	targets, err := GetTargets(c, questionnaireID)
	if err != nil {
		return err
	}

	administrators, err := GetAdministrators(c, questionnaireID)
	if err != nil {
		return err
	}

	respondents := []string{}
	if err := db.Select(&respondents, "SELECT user_traqid FROM respondents WHERE questionnaire_id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionnaireID": questionnaire.ID,
		"title":           questionnaire.Title,
		"description":     questionnaire.Description,
		"res_time_limit":  timeConvert(questionnaire.ResTimeLimit),
		"created_at":      questionnaire.CreatedAt,
		"modified_at":     questionnaire.ModifiedAt,
		"res_shared_to":   questionnaire.ResSharedTo,
		"targets":         targets,
		"administrators":  administrators,
		"respondents":     respondents,
	})
}

func postQuestionnaire(c echo.Context) error {

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
		result, err = db.Exec(
			"INSERT INTO questionnaires (title, description, res_shared_to) VALUES (?, ?, ?)",
			req.Title, req.Description, req.ResSharedTo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		var err error
		result, err = db.Exec(
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

	if err := InsertTargets(c, int(lastID), req.Targets); err != nil {
		return err
	}

	if err := InsertAdministrators(c, int(lastID), req.Administrators); err != nil {
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

func editQuestionnaire(c echo.Context) error {

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
		if _, err := db.Exec(
			"UPDATE questionnaires SET title = ?, description = ?, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?",
			req.Title, req.Description, req.ResSharedTo, questionnaireID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := db.Exec(
			"UPDATE questionnaires SET title = ?, description = ?, res_time_limit = ?, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?",
			req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo, questionnaireID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	if err := DeleteTargets(c, questionnaireID); err != nil {
		return err
	}

	if err := InsertTargets(c, questionnaireID, req.Targets); err != nil {
		return err
	}

	if err := DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	if err := InsertAdministrators(c, questionnaireID, req.Administrators); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func deleteQuestionnaire(c echo.Context) error {

	questionnaireID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if _, err := db.Exec(
		"UPDATE questionnaires SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := DeleteTargets(c, questionnaireID); err != nil {
		return err
	}

	if err := DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func GetTitleAndLimit(c echo.Context, questionnaireID int) (string, string, error) {
	res := struct {
		Title        string `db:"title"`
		ResTimeLimit string `db:"res_time_limit"`
	}{}
	if err := db.Get(&res,
		"SELECT title, res_time_limit FROM questionnaires WHERE id = ? AND deleted_at IS NULL",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return "", "", echo.NewHTTPError(http.StatusNotFound)
		} else {
			return "", "", echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return res.Title, res.ResTimeLimit, nil
}

func getMyQuestionnaire(c echo.Context) error {
	/*
		後で書く
		questionnaireID, err := GetAdminQuestionnaires(c, getUserID(c))
		if err != nil {
			return nil
		}*/
	return c.NoContent(http.StatusOK)
}
