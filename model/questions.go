package model

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/go-sql-driver/mysql"
)

//Question questionテーブルの構造体
type Question struct {
	ID              int            `json:"id"                  gorm:"type:int(11);PRIMARY_KEY;NOT NULL;AUTO_INCREMENT;"`
	QuestionnaireID int            `json:"questionnaireID"     gorm:"type:int(11);DEFAULT:NULL;"`
	PageNum         int            `json:"page_num"            gorm:"type:int(11);NOT NULL;"`
	QuestionNum     int            `json:"question_num"        gorm:"type:int(11);NOT NULL;"`
	Type            string         `json:"type"                gorm:"type:char(20);NOT NULL;"`
	Body            string         `json:"body"                gorm:"type:text;"`
	IsRequrired     bool           `json:"is_required"         gorm:"type:tinyint(4);NOT NULL;"`
	DeletedAt       mysql.NullTime `json:"deleted_at"          gorm:"type:timestamp;"`
	CreatedAt       time.Time      `json:"created_at"          gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;"`
}

//TableName テーブル名が単数形なのでその対応
func (*Question) TableName() string {
	return "question"
}

//QuestionIDType 質問のIDと種類の構造体
type QuestionIDType struct {
	ID   int
	Type string
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

func GetQuestions(c echo.Context, questionnaireID int) ([]Question, error) {
	allquestions := []Question{}

	// アンケートidの一致する質問を取る
	if err := db.Select(&allquestions,
		"SELECT * FROM question WHERE questionnaire_id = ? AND deleted_at IS NULL ORDER BY question_num",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return []Question{}, echo.NewHTTPError(http.StatusInternalServerError)
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
