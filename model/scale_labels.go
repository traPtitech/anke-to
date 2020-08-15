package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//ScaleLabels scale_labelsテーブルの構造体
type ScaleLabels struct {
	ID              int    `json:"questionID"        gorm:"column:question_id"`
	ScaleLabelRight string `json:"scale_label_right" gorm:"column:scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"  gorm:"column:scale_label_left"`
	ScaleMin        int    `json:"scale_min"         gorm:"column:scale_min"`
	ScaleMax        int    `json:"scale_max"         gorm:"column:scale_max"`
}

// GetScaleLabels 指定されたquestionIDのlabelを取得する
func GetScaleLabels(c echo.Context, questionID int) (ScaleLabels, error) {
	label := ScaleLabels{}
	if err := gormDB.Where("question_id = ?", questionID).First(&label).Error; gorm.IsRecordNotFoundError(err) {
		return ScaleLabels{}, nil
	} else if err != nil {
		c.Logger().Error(err)
		return ScaleLabels{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return label, nil
}

// InsertScaleLabels IDを指定してlabelを挿入する
func InsertScaleLabels(c echo.Context, lastID int, label ScaleLabels) error {
	label.ID = lastID
	if err := gormDB.Create(&label).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateScaleLabels questionIDを指定してlabelを更新する
func UpdateScaleLabels(c echo.Context, questionID int, label ScaleLabels) error {
	if err := gormDB.Model(&ScaleLabels{}).Update(&label).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// DeleteScaleLabels questionIDを指定してlabelを削除する
func DeleteScaleLabels(c echo.Context, questionID int) error {
	if err := gormDB.Where("question_id = ?", questionID).Delete(&ScaleLabels{}).Error; err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
