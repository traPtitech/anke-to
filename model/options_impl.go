package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

// Options optionsテーブルの構造体
type Options struct {
	ID         int    `gorm:"type:int(11) AUTO_INCREMENT NOT NULL PRIMARY KEY;"`
	QuestionID int    `gorm:"type:int(11) NOT NULL;"`
	OptionNum  int    `gorm:"type:int(11) NOT NULL;"`
	Body       string `gorm:"type:text;default:NULL;"`
}

// Option OptionRepositoryの実装
type Option struct{}

// InsertOption 選択肢の追加
func (*Option) InsertOption(lastID int, num int, body string) error {
	option := Options{
		QuestionID: lastID,
		OptionNum:  num,
		Body:       body,
	}
	err := db.Create(&option).Error
	if err != nil {
		return fmt.Errorf("failed to insert a option: %w", err)
	}
	return nil
}

// UpdateOptions 選択肢の修正
func (*Option) UpdateOptions(options []string, questionID int) error {
	var err error
	for i, optionLabel := range options {
		option := Options{
			Body: optionLabel,
		}
		query := db.
			Model(Options{}).
			Where("question_id = ? AND option_num = ?", questionID, i+1)
		err := query.First(&Options{}).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return fmt.Errorf("failed to get option: %w", err)
		}

		if gorm.IsRecordNotFoundError(err) {
			option.QuestionID = questionID
			option.OptionNum = i + 1
			err = db.Create(&option).Error
			if err != nil {
				return fmt.Errorf("failed to insert option: %w", err)
			}
		} else {
			result := query.Update(&option)
			err = result.Error
			if err != nil {
				return fmt.Errorf("failed to update option: %w", err)
			}
		}
	}
	err = db.Where("question_id = ? AND option_num > ?", questionID, len(options)).Delete(Options{}).Error
	if err != nil {
		return fmt.Errorf("failed to update option: %w", err)
	}
	return nil
}

// DeleteOptions 選択肢の削除
func (*Option) DeleteOptions(questionID int) error {
	err := db.
		Where("question_id = ?", questionID).
		Delete(Options{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete option: %w", err)
	}
	return nil
}

// GetOptions 質問の選択肢の取得
func (*Option) GetOptions(questionIDs []int) ([]Options, error) {
	type option struct {
		QuestionID int         `gorm:"type:int(11) NOT NULL;"`
		Body       null.String `gorm:"type:text;default:NULL;"`
	}
	options := []option{}

	err := db.
		Model(Options{}).
		Where("question_id IN (?)", questionIDs).
		Order("option_num").
		Select("question_id, body").
		Find(&options).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get option: %w", err)
	}

	optns := make([]Options, 0, len(options))
	for _, optn := range options {
		optns = append(optns, Options{
			QuestionID: optn.QuestionID,
			Body:       optn.Body.ValueOrZero(),
		})
	}

	return optns, nil
}
