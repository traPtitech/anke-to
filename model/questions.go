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
	if err := DB.Select(&ret,
		`SELECT id, type FROM question WHERE questionnaire_id = ? AND deleted_at IS NULL`,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ret, nil
}
