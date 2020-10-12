package model

import (
	"net/http"

	"github.com/labstack/echo"
)

// Administrators administratorsテーブルの構造体
type Administrators struct {
	QuestionnaireID int    `sql:"type:int(11);not null;primary_key;"`
	UserTraqid      string `sql:"type:char(32);not null;primary_key;"`
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

// GetAdministrators アンケートの管理者を取得
func GetAdministrators(c echo.Context, questionnaireID int) ([]string, error) {
	userTraqids := []string{}
	err := db.
		Model(&Administrators{}).
		Where("questionnaire_id = ?", questionnaireID).
		Pluck("user_traqid", &userTraqids).Error
	if err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return userTraqids, nil
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
func CheckAdmin(c echo.Context, questionnaireID int) (bool, error) {
	user := GetUserID(c)
	administrators, err := GetAdministrators(c, questionnaireID)
	if err != nil {
		c.Logger().Error(err)
		return false, err
	}

	found := false
	for _, admin := range administrators {
		if admin == user || admin == "traP" {
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}
	return true, nil
}
