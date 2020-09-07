package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
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
func GetScaleLabels(questionID int) (ScaleLabels, error) {
	label := ScaleLabels{}
	if err := gormDB.Where("question_id = ?", questionID).First(&label).Error; gorm.IsRecordNotFoundError(err) {
		return ScaleLabels{}, nil
	} else if err != nil {
		return ScaleLabels{}, fmt.Errorf("failed to get the scale label (questionID: %d): %w", questionID, err)
	}
	return label, nil
}

// InsertScaleLabels IDを指定してlabelを挿入する
func InsertScaleLabels(lastID int, label ScaleLabels) error {
	label.ID = lastID
	if err := gormDB.Create(&label).Error; err != nil {
		return fmt.Errorf("failed to insert the scale label (lastID: %d): %w", lastID, err)
	}
	return nil
}

// UpdateScaleLabels questionIDを指定してlabelを更新する
func UpdateScaleLabels(questionID int, label ScaleLabels) error {
	if err := gormDB.Model(&ScaleLabels{}).Update(&label).Error; err != nil {
		return fmt.Errorf("failed to update the scale labell (questionID: %d): %w", questionID, err)
	}
	return nil
}

// DeleteScaleLabels questionIDを指定してlabelを削除する
func DeleteScaleLabels(questionID int) error {
	if err := gormDB.Where("question_id = ?", questionID).Delete(&ScaleLabels{}).Error; err != nil {
		return fmt.Errorf("failed to delete the scale labell (questionID: %d): %w", questionID, err)
	}
	return nil
}
