package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

//Respondents respondentsテーブルの構造体
type Respondents struct {
	ResponseID int `gorm:"type:int(11);NOT NULL;PRIMARY_KEY;AUTO_INCREMENT"`
	QuestionnaireID int `gorm:"type:int(11);NOT NULL;"`
	UserTraqid string `gorm:"type:char(30);NOT NULL;"`
	ModifiedAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT CURRENT_TIMESTAMP;"`
	SubmittedAt null.Time `gorm:"type:timestamp;"`
	DeletedAT null.Time `gorm:"type:timestamp;"`
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (*Respondents) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedAt", time.Now())

	return nil
}

//InsertRespondent 回答者の追加
func InsertRespondent(c echo.Context, questionnaireID int, submitedAt null.Time) (int, error) {
	userID := GetUserID(c)

	var respondent Respondents
	if submitedAt.Valid {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid: userID,
			SubmittedAt: submitedAt,
		}
	} else {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid: userID,
		}
	}

	err := gormDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&respondent).Error
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to insert a respondent record: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		err = tx.Select("response_id").Last(&respondent).Error
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to get the last respondent record: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return nil
	})
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed in transaction: %w", err))
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return respondent.ResponseID, nil
}