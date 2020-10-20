package model

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// Administrators administratorsテーブルの構造体
type Administrators struct {
	QuestionnaireID int `gorm:"primary_key"`
	UserTraqid      string
}

// InsertAdministrators アンケートの管理者を追加
func InsertAdministrators(c echo.Context, questionnaireID int, administrators []string) error {
	var administrator Administrators
	var err error
	for _, v := range administrators {
		administrator = Administrators{
			QuestionnaireID: questionnaireID,
			UserTraqid:      v,
		}
		err = db.Create(&administrator).Error
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

// DeleteAdministrators アンケートの管理者の削除
func DeleteAdministrators(c echo.Context, questionnaireID int) error {
	err := db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(Administrators{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// GetAdminQuestionnaireIDs 自分が管理者のアンケートの取得
func GetAdminQuestionnaireIDs(c echo.Context, user string) ([]int, error) {
	questionnaireIDs := []int{}
	err := db.
		Model(&Administrators{}).
		Where("user_traqid = ?", user).
		Or("user_traqid = ?", "traP").
		Pluck("DISTINCT questionnaire_id", &questionnaireIDs).Error
	if err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return questionnaireIDs, nil
}

// CheckAdmin 自分がアンケートの管理者か判定
func CheckAdmin(userID string, questionnaireID int) (bool, error) {
	err := db.
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		Find(&Administrators{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get a administrator: %w", err)
	}

	return true, nil
}
