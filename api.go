package main

import (
	"net/http"
	"strconv"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func timeConvert(time mysql.NullTime) string {
	if time.Valid {
		return time.Time.String()
	} else {
		return "NULL"
	}
}

func getUserID(c echo.Context) string {
	return c.Request().Header.Get("X-Showcase-User")
}

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

type TargetType int

const (
	Targeted = iota
	Nontargeted
	All
)

// echo.Contextを引数にとってerrorを返り値とする
func getQuestionnaires(c echo.Context, targettype TargetType) error { /*
		// query parametar
		nontargeted := c.QueryParam("nontargeted") == "true"*/

	allquestionnaires, err := getAllQuestionnaires(c)
	if err != nil {
		return err
	}

	userID := getUserID(c)
	// test用
	if userID == "" {
		userID = "mds_boy"
	}

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

	questionnaireID := c.Param("id")

	questionnaire := questionnaires{}
	if err := db.Get(&questionnaire, "SELECT * FROM questionnaires WHERE id = ? AND deleted_at IS NULL", questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	targets := []string{}
	if err := db.Select(&targets, "SELECT user_traqid FROM targets WHERE questionnaire_id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	administrators := []string{}
	if err := db.Select(&administrators, "SELECT user_traqid FROM administrators WHERE questionnaire_id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
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

func getQuestions(c echo.Context) error {
	questionnaireID := c.Param("id")
	// 質問一覧の配列
	allquestions := []questions{}

	// アンケートidの一致する質問を取る
	if err := db.Select(
		&allquestions, "SELECT * FROM questions WHERE questionnaire_id = ? AND deleted_at IS NULL", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if len(allquestions) == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	type questionInfo struct {
		QuestionID      int       `json:"question_ID"`
		PageNum         int       `json:"page_num"`
		QuestionNum     int       `json:"question_num"`
		QuestionType    string    `json:"question_type"`
		Body            string    `json:"body"`
		IsRequrired     bool      `json:"is_required"`
		CreatedAt       time.Time `json:"created_at"`
		Options         []string  `json:"options"`
		ScaleLabelRight string    `json:"scale_label_right"`
		ScaleLabelLeft  string    `json:"scale_label_left"`
		ScaleMin        int       `json:"scale_min"`
		ScaleMax        int       `json:"sclae_max"`
	}
	var ret []questionInfo

	for _, v := range allquestions {
		options := []string{}
		scalelabel := scaleLabels{}

		switch v.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			if err := db.Select(
				&options, "SELECT body FROM options WHERE question_id = ? ORDER BY option_num",
				v.ID); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		case "LinearScale":
			if err := db.Get(&scalelabel, "SELECT * FROM scale_labels WHERE question_id = ?", v.ID); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}

		ret = append(ret,
			questionInfo{
				QuestionID:      v.ID,
				PageNum:         v.PageNum,
				QuestionNum:     v.QuestionNum,
				QuestionType:    v.Type,
				Body:            v.Body,
				IsRequrired:     v.IsRequrired,
				CreatedAt:       v.CreatedAt,
				Options:         options,
				ScaleLabelRight: scalelabel.ScaleLabelRight,
				ScaleLabelLeft:  scalelabel.ScaleLabelLeft,
				ScaleMin:        scalelabel.ScaleMin,
				ScaleMax:        scalelabel.ScaleMax,
			})
	}

	return c.JSON(http.StatusOK, ret)
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

	if req.ResSharedTo == "" {
		req.ResSharedTo = "administrators"
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

	// エラーチェック
	lastID, err := result.LastInsertId()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	for _, v := range req.Targets {
		if _, err := db.Exec(
			"INSERT INTO targets (questionnaire_id, user_traqid) VALUES (?, ?)",
			lastID, v); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	for _, v := range req.Administrators {
		if _, err := db.Exec(
			"INSERT INTO administrators (questionnaire_id, user_traqid) VALUES (?, ?)",
			lastID, v); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
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
	questionnaireID := c.Param("id")

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

	// TargetsとAdministratorsの変更はまだ

	return c.NoContent(http.StatusOK)
}

func deleteQuestionnaire(c echo.Context) error {
	questionnaireID := c.Param("id")

	if _, err := db.Exec(
		"UPDATE questionnaires SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func postQuestion(c echo.Context) error {
	req := struct {
		QuestionnaireID int      `json:"questionnaireID"`
		QuestionType    string   `json:"question_type"`
		QuestionNum     int      `json:"question_num"`
		PageNum         int      `json:"page_num"`
		Body            string   `json:"body"`
		IsRequrired     bool     `json:"is_required"`
		Options         []string `json:"options"`
		ScaleLabelRight string   `json:"scale_label_right"`
		ScaleLabelLeft  string   `json:"scale_label_left"`
		ScaleMin        int      `json:"scale_min"`
		ScaleMax        int      `json:"scale_max"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	result, err := db.Exec(
		"INSERT INTO questions (questionnaire_id, page_num, question_num, type, body, is_required) VALUES (?, ?, ?, ?, ?, ?)",
		req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body, req.IsRequrired)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	lastID, err2 := result.LastInsertId()
	if err2 != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		for i, v := range req.Options {
			if _, err := db.Exec(
				"INSERT INTO options (question_id, option_num, body) VALUES (?, ?, ?)",
				lastID, i+1, v); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}
	case "LinearScale":
		if _, err := db.Exec(
			"INSERT INTO scale_labels (question_id, scale_label_left, scale_label_right, scale_min, scale_max) VALUES (?, ?, ?, ?, ?)",
			lastID, req.ScaleLabelLeft, req.ScaleLabelRight, req.ScaleMin, req.ScaleMax); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionID":        int(lastID),
		"questionnaireID":   req.QuestionnaireID,
		"question_type":     req.QuestionType,
		"question_num":      req.QuestionNum,
		"page_num":          req.PageNum,
		"body":              req.Body,
		"is_required":       req.IsRequrired,
		"options":           req.Options,
		"scale_label_right": req.ScaleLabelRight,
		"scale_label_left":  req.ScaleLabelLeft,
		"scale_max":         req.ScaleMax,
		"scale_min":         req.ScaleMin,
	})
}

func editQuestion(c echo.Context) error {
	questionID := c.Param("id")

	req := struct {
		QuestionnaireID int      `json:"questionnaireID"`
		QuestionType    string   `json:"question_type"`
		QuestionNum     int      `json:"question_num"`
		PageNum         int      `json:"page_num"`
		Body            string   `json:"body"`
		IsRequrired     bool     `json:"is_required"`
		Options         []string `json:"options"`
		ScaleLabelRight string   `json:"scale_label_right"`
		ScaleLabelLeft  string   `json:"scale_label_left"`
		ScaleMax        int      `json:"scale_max"`
		ScaleMin        int      `json:"scale_min"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if _, err := db.Exec(
		"UPDATE questions SET questionnaire_id = ?, page_num = ?, question_num = ?, type = ?, body = ?, is_required = ? WHERE id = ?",
		req.QuestionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body, req.IsRequrired, questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		for i, v := range req.Options {
			if _, err := db.Exec(
				`INSERT INTO options (question_id, option_num, body) VALUES (?, ?, ?)
				ON DUPLICATE KEY UPDATE option_num = ?, body = ?`,
				questionID, i+1, v, i+1, v); err != nil {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}
		if _, err := db.Exec(
			"DELETE FROM options WHERE question_id= ? AND option_num > ?",
			questionID, len(req.Options)); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	case "LinearScale":
		if _, err := db.Exec(
			`INSERT INTO scale_labels (question_id, scale_label_right, scale_label_left, scale_min, scale_max) VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE scale_label_right = ?, scale_label_left = ?, scale_min = ?, scale_max = ?`,
			questionID,
			req.ScaleLabelRight, req.ScaleLabelLeft, req.ScaleMin, req.ScaleMax,
			req.ScaleLabelRight, req.ScaleLabelLeft, req.ScaleMin, req.ScaleMax); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.NoContent(http.StatusOK)
}

func deleteQuestion(c echo.Context) error {
	questionID := c.Param("id")

	if _, err := db.Exec(
		"UPDATE questions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
