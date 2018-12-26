package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

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

	if _, err := db.Exec(
		"DELETE FROM options WHERE question_id= ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if _, err := db.Exec(
		"DELETE FROM options WHERE question_id= ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func getQuestionsType(c echo.Context, questionnaireID int) ([]questionIDType, error) {
	ret := []questionIDType{}
	if err := db.Select(&ret,
		`SELECT id, type FROM questions WHERE questionnaire_id = ? AND deleted_at IS NULL`,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ret, nil
}
