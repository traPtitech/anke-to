package model

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// ScaleLabel ScaleLabelRepositoryの実装
type ScaleLabel struct{}

// NewScaleLabel ScaleLabelのコンストラクター
func NewScaleLabel() *ScaleLabel {
	return new(ScaleLabel)
}

//ScaleLabels scale_labelsテーブルの構造体
type ScaleLabels struct {
	QuestionID      int    `json:"questionID"        gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	ScaleLabelRight string `json:"scale_label_right" gorm:"type:text;default:NULL;"`
	ScaleLabelLeft  string `json:"scale_label_left"  gorm:"type:text;default:NULL;"`
	ScaleMin        int    `json:"scale_min"         gorm:"type:int(11);default:NULL;"`
	ScaleMax        int    `json:"scale_max"         gorm:"type:int(11);default:NULL;"`
}

// InsertScaleLabel IDを指定してlabelを挿入する
func (*ScaleLabel) InsertScaleLabel(lastID int, label ScaleLabels) error {
	label.QuestionID = lastID
	err := db.
		Session(&gorm.Session{NewDB: true}).
		Create(&label).Error
	if err != nil {
		return fmt.Errorf("failed to insert the scale label (lastID: %d): %w", lastID, err)
	}
	return nil
}

// UpdateScaleLabel questionIDを指定してlabelを更新する
func (*ScaleLabel) UpdateScaleLabel(questionID int, label ScaleLabels) error {
	result := db.
		Session(&gorm.Session{NewDB: true}).
		Model(&ScaleLabels{}).
		Where("question_id = ?", questionID).
		Updates(map[string]interface{}{
			"question_id":       questionID,
			"scale_label_right": label.ScaleLabelRight,
			"scale_label_left":  label.ScaleLabelLeft,
			"scale_min":         label.ScaleMin,
			"scale_max":         label.ScaleMax,
		})
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to update the scale label (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update a scale label record: %w", ErrNoRecordUpdated)
	}
	return nil
}

// DeleteScaleLabel questionIDを指定してlabelを削除する
func (*ScaleLabel) DeleteScaleLabel(questionID int) error {
	result := db.
		Session(&gorm.Session{NewDB: true}).
		Where("question_id = ?", questionID).
		Delete(&ScaleLabels{})
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to delete the scale label (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete a scale label : %w", ErrNoRecordDeleted)
	}
	return nil
}

// GetScaleLabels 指定されたquestionIDの配列のlabelを取得する
func (*ScaleLabel) GetScaleLabels(questionIDs []int) ([]ScaleLabels, error) {
	labels := []ScaleLabels{}
	err := db.
		Session(&gorm.Session{NewDB: true}).
		Where("question_id IN (?)", questionIDs).
		Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get the scale label: %w", err)
	}

	return labels, nil
}

// CheckScaleLabel responseがScaleMin,ScaleMaxを満たしているか
func (*ScaleLabel) CheckScaleLabel(label ScaleLabels, response string) error {
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
