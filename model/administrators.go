package model

import (
	"net/http"

	"github.com/labstack/echo"
)

type Administrator struct {
	QuestionnaireID int `gorm:"primary_key"`
	UserTraqid      string
}

func GetAdministrators(c echo.Context, questionnaireID int) ([]string, error) {
	userTraqids := []string{}
	administrators := []Administrator{}
	err := gormDB.Model(&Administrator{}).Where("questionnaire_id = ?", questionnaireID).Pluck("user_traqid", &userTraqids).Error
	if err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return userTraqids, nil
}

func InsertAdministrators(c echo.Context, questionnaireID int, administrators []string) error {
	var administrator Administrator
	var err error
	for _, v := range administrators {
		administrator = Administrator{
			QuestionnaireID: questionnaireID,
			UserTraqid:      v,
		}
		err = gormDB.Create(&administrator).Error
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

func DeleteAdministrators(c echo.Context, questionnaireID int) error {
	err := gormDB.Where("questionnaire_id = ?", questionnaireID).Delete(Administrator{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func GetAdminQuestionnaires(c echo.Context, user string) ([]int, error) {
	questionnaireIDs := []int{}
	administrators := []Administrator{}
	err := gormDB.Where("user_traqid = ?", user).Or("user_traqid = ?", "traP").Select("DISTINCT").Find(&administrators).Pluck("questionnaire_id", &questionnaireIDs).Error
	if err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return questionnaireIDs, nil
}

// 自分がadminなら(true, nil)
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
