package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
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
	IsRequired     bool           `json:"is_required"         gorm:"type:tinyint(4);NOT NULL;"`
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

//GetQuestionTypes 質問のIDと型の配列を取得
func GetQuestionTypes(c echo.Context, questionnaireID int) ([]QuestionIDType, error) {
	questionIDTypes := []QuestionIDType{}
	err := gormDB.
		Model(&Question{}).
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Select("id, type").
		Scan(&questionIDTypes).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get question`s ids and types: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusNotFound)
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return questionIDTypes, nil
}

//GetQuestions 質問のリストの取得
func GetQuestions(c echo.Context, questionnaireID int) ([]Question, error) {
	questions := []Question{}
	err := gormDB.
		Where("questionnaire_id = ?", questionnaireID).
		Find(&questions).Error
	// アンケートidの一致する質問を取る
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get questions: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusNotFound)
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return questions, nil
}

//InsertQuestion 質問の追加
func InsertQuestion(
	c echo.Context, questionnaireID int, pageNum int, questionNum int, questionType string,
	body string, isRequired bool) (int, error) {
	question := Question{
		QuestionnaireID: questionnaireID,
		PageNum: pageNum,
		QuestionNum: questionNum,
		Type: questionType,
		Body: body,
		IsRequired: isRequired,
	}

	err := gormDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&question).Error
		if err != nil {
			return fmt.Errorf("failed to insert a question record: %w", err)
		}

		err = tx.Select("id").Last(&question).Error
		if err != nil {
			return fmt.Errorf("failed to get the last question record: %w", err)
		}

		return nil
	})
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed in transaction: %w", err))
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return question.ID, nil
}

//UpdateQuestion 質問の修正
func UpdateQuestion(
	c echo.Context, questionnaireID int, pageNum int, questionNum int, questionType string,
	body string, isRequired bool, questionID int) error {
	question := Question{
		QuestionnaireID: questionnaireID,
		PageNum: pageNum,
		QuestionNum: questionNum,
		Type: questionType,
		Body: body,
		IsRequired: isRequired,
	}

	err := gormDB.Model(&Question{}).Update(&question).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to update a question record: %w", err))
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
