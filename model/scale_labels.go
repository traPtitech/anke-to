package model

import (
	"fmt"
	"strconv"

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
	if err := gormDB.
		Where("question_id = ?", questionID).
		First(&label).
		Error; gorm.IsRecordNotFoundError(err) {
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
	if err := gormDB.
		Model(&ScaleLabels{}).
		Where("question_id = ?", questionID).
		Update(map[string]interface{}{
			"question_id":       questionID,
			"scale_label_right": label.ScaleLabelRight,
			"scale_label_left":  label.ScaleLabelLeft,
			"scale_min":         label.ScaleMin,
			"scale_max":         label.ScaleMax}).
		Error; err != nil {
		return fmt.Errorf("failed to update the scale labell (questionID: %d): %w", questionID, err)
	}
	return nil
}

// DeleteScaleLabels questionIDを指定してlabelを削除する
func DeleteScaleLabels(questionID int) error {
	if err := gormDB.
		Where("question_id = ?", questionID).
		Delete(&ScaleLabels{}).
		Error; err != nil {
		return fmt.Errorf("failed to delete the scale labell (questionID: %d): %w", questionID, err)
	}
	return nil
}

// CheckScaleLabels responceがScaleMin,ScaleMaxを満たしているか
func CheckScaleLabels(label ScaleLabels, responce string) error {
	if responce == "" {
		return nil
	}

	r, err := strconv.Atoi(responce)
	if err != nil {
		return err
	}
	if r < label.ScaleMin {
		return fmt.Errorf("failed to meet the scale. the responce must be greater than ScaleMin (number: %d, ScaleMin: %d)", r, label.ScaleMin)
	} else if r > label.ScaleMax {
		return fmt.Errorf("failed to meet the scale. the responce must be less than ScaleMax (number: %d, ScaleMax: %d)", r, label.ScaleMax)
	}

	return nil
}
