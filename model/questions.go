package model

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/go-sql-driver/mysql"
)

type Questions struct {
	ID              int            `json:"id"                  db:"id"`
	QuestionnaireId int            `json:"questionnaireID"     db:"questionnaire_id"`
	PageNum         int            `json:"page_num"            db:"page_num"`
	QuestionNum     int            `json:"question_num"        db:"question_num"`
	Type            string         `json:"type"                db:"type"`
	Body            string         `json:"body"                db:"body"`
	IsRequrired     bool           `json:"is_required"         db:"is_required"`
	DeletedAt       mysql.NullTime `json:"deleted_at"          db:"deleted_at"`
	CreatedAt       time.Time      `json:"created_at"          db:"created_at"`
}

type QuestionIDType struct {
	ID   int    `db:"id"`
	Type string `db:"type"`
}

func GetQuestionsType(c echo.Context, questionnaireID int) ([]QuestionIDType, error) {
	ret := []QuestionIDType{}
	if err := db.Select(&ret,
		`SELECT id, type FROM question WHERE questionnaire_id = ? AND deleted_at IS NULL ORDER BY question_num`,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ret, nil
}

func GetQuestions(c echo.Context, questionnaireID int) ([]Questions, error) {
	allquestions := []Questions{}

	// アンケートidの一致する質問を取る
	if err := db.Select(&allquestions,
		"SELECT * FROM question WHERE questionnaire_id = ? AND deleted_at IS NULL ORDER BY question_num",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return []Questions{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return allquestions, nil
}

func InsertQuestion(
	c echo.Context, questionnaireID int, pageNum int, questionNum int, questionType string,
	body string, isRequired bool) (int, error) {
	result, err := db.Exec(
		`INSERT INTO question (questionnaire_id, page_num, question_num, type, body, is_required, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		questionnaireID, pageNum, questionNum, questionType, body, isRequired, time.Now())
	if err != nil {
		c.Logger().Error(err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		c.Logger().Error(err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return int(lastID), nil
}

func UpdateQuestion(
	c echo.Context, questionnaireID int, pageNum int, questionNum int, questionType string,
	body string, isRequired bool, questionID int) error {
	if _, err := db.Exec(
		"UPDATE question SET questionnaire_id = ?, page_num = ?, question_num = ?, type = ?, body = ?, is_required = ? WHERE id = ?",
		questionnaireID, pageNum, questionNum, questionType, body, isRequired, questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteQuestion(c echo.Context, questionID int) error {
	if _, err := db.Exec(
		"UPDATE question SET deleted_at = ? WHERE id = ?",
		time.Now(), questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
