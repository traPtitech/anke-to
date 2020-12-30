package model

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

//ScaleLabels scale_labelsテーブルの構造体
type ScaleLabels struct {
	QuestionID      int    `json:"questionID"        gorm:"type:int(11) NOT NULL PRIMARY KEY;"`
	ScaleLabelRight string `json:"scale_label_right" gorm:"type:text;default:NULL;"`
	ScaleLabelLeft  string `json:"scale_label_left"  gorm:"type:text;default:NULL;"`
	ScaleMin        int    `json:"scale_min"         gorm:"type:int(11);default:NULL;"`
	ScaleMax        int    `json:"scale_max"         gorm:"type:int(11);default:NULL;"`
}

// InsertScaleLabel IDを指定してlabelを挿入する
func InsertScaleLabel(lastID int, label ScaleLabels) error {
	label.QuestionID = lastID
	if err := db.Create(&label).Error; err != nil {
		return fmt.Errorf("failed to insert the scale label (lastID: %d): %w", lastID, err)
	}
	return nil
}

// UpdateScaleLabel questionIDを指定してlabelを更新する
func UpdateScaleLabel(questionID int, label ScaleLabels) error {
	err := db.
		Model(&ScaleLabels{}).
		Where("question_id = ?", questionID).
		Update(map[string]interface{}{
			"question_id":       questionID,
			"scale_label_right": label.ScaleLabelRight,
			"scale_label_left":  label.ScaleLabelLeft,
			"scale_min":         label.ScaleMin,
			"scale_max":         label.ScaleMax}).
		Error
	if err != nil {
		return fmt.Errorf("failed to update the scale labell (questionID: %d): %w", questionID, err)
	}
	return nil
}

// DeleteScaleLabel questionIDを指定してlabelを削除する
func DeleteScaleLabel(questionID int) error {
	err := db.
		Where("question_id = ?", questionID).
		Delete(&ScaleLabels{}).
		Error
	if err != nil {
		return fmt.Errorf("failed to delete the scale labell (questionID: %d): %w", questionID, err)
	}
	return nil
}

// GetScaleLabel 指定されたquestionIDのlabelを取得する
func GetScaleLabel(questionID int) (ScaleLabels, error) {
	label := ScaleLabels{}
	err := db.
		Where("question_id = ?", questionID).
		First(&label).
		Error
	if gorm.IsRecordNotFoundError(err) {
		return ScaleLabels{}, nil
	} else if err != nil {
		return ScaleLabels{}, fmt.Errorf("failed to get the scale label (questionID: %d): %w", questionID, err)
	}
	return label, nil
}

// GetScaleLabels 指定されたquestionIDの配列のlabelを取得する
func GetScaleLabels(questionIDs ...int) (map[int]ScaleLabels, error) {
	labels := []ScaleLabels{}
	err := db.
		Where("question_id IN (?)", questionIDs).
		Find(&labels).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get the scale label: %w", err)
	}

	labelMap := make(map[int]ScaleLabels, len(labels))
	for _, label := range labels {
		labelMap[label.QuestionID] = label
	}

	return labelMap, nil
}

// CheckScaleLabel responseがScaleMin,ScaleMaxを満たしているか
func CheckScaleLabel(label ScaleLabels, response string) error {
	if response == "" {
		return nil
	}

	r, err := strconv.Atoi(response)
	if err != nil {
		return err
	}
	if r < label.ScaleMin {
		return fmt.Errorf("failed to meet the scale. the response must be greater than ScaleMin (number: %d, ScaleMin: %d)", r, label.ScaleMin)
	} else if r > label.ScaleMax {
		return fmt.Errorf("failed to meet the scale. the response must be less than ScaleMax (number: %d, ScaleMax: %d)", r, label.ScaleMax)
	}

	return nil
}
