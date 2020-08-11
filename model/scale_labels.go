package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type ScaleLabels struct {
	ID              int    `json:"questionID"        db:"question_id"       gorm:"column:question_id"`
	ScaleLabelRight string `json:"scale_label_right" db:"scale_label_right" gorm:"column:scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"  db:"scale_label_left"  gorm:"column:scale_label_left"`
	ScaleMin        int    `json:"scale_min"         db:"scale_min"         gorm:"column:scale_min"`
	ScaleMax        int    `json:"scale_max"         db:"scale_max"         gorm:"column:scale_max"`
}

func GetScaleLabels(c echo.Context, questionID int) (ScaleLabels, error) {
	label := ScaleLabels{}
	if err := gormDB.Where("question_id = ?", questionID).First(&label).Error; err != nil {
		c.Logger().Error(err)
		return ScaleLabels{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return label, nil
}

func InsertScaleLabels(c echo.Context, lastID int, label ScaleLabels) error {
	label.ID = lastID
	if err := gormDB.Create(&label).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func UpdateScaleLabels(c echo.Context, questionID int, label ScaleLabels) error {
	labelBefore := ScaleLabels{}

	var err error
	if labelBefore, err = GetScaleLabels(c, questionID); gorm.IsRecordNotFoundError(err) {
		return nil
	} else if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := gormDB.Model(&labelBefore).Update(&label).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteScaleLabels(c echo.Context, questionID int) error {
	if err := gormDB.Where("question_id = ?", questionID).Delete(&ScaleLabels{}).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

//debug用　後で消す
func GetScaleLabelLists(c echo.Context) ([]ScaleLabels, error) {
	lists := []ScaleLabels{}
	if err := gormDB.Find(&lists).Error; err != nil {
		c.Logger().Error(err)
		return []ScaleLabels{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return lists, nil
}
