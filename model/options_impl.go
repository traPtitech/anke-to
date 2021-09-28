package model

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/guregu/null.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Option OptionRepositoryの実装
type Option struct{}

// NewOption Optionのコンストラクター
func NewOption() *Option {
	return new(Option)
}

// Options optionsテーブルの構造体
type Options struct {
	ID         int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionID int    `gorm:"type:int(11);not null"`
	OptionNum  int    `gorm:"type:int(11);not null"`
	Body       string `gorm:"type:text;default:NULL;"`
}

// InsertOption 選択肢の追加
func (*Option) InsertOption(ctx context.Context, lastID int, num int, body string) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	option := Options{
		QuestionID: lastID,
		OptionNum:  num,
		Body:       body,
	}
	err = db.Create(&option).Error
	if err != nil {
		return fmt.Errorf("failed to insert a option: %w", err)
	}
	return nil
}

// UpdateOptions 選択肢の修正
func (*Option) UpdateOptions(ctx context.Context, options []string, questionID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	var previousOptions []Options
	err = db.
		Session(&gorm.Session{}).
		Where("question_id = ?", questionID).
		Select("OptionNum", "Body").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&previousOptions).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get option: %w", err)
	}

	isDelete := false
	optionMap := make(map[int]*Options, len(options))
	for i, option := range previousOptions {
		if option.OptionNum <= len(options) {
			optionMap[option.OptionNum] = &previousOptions[i]
		} else {
			isDelete = true
		}
	}

	createOptions := []Options{}
	for i, optionLabel := range options {
		optionNum := i + 1

		if option, ok := optionMap[optionNum]; ok {
			if option.Body != optionLabel {
				err := db.
					Session(&gorm.Session{}).
					Model(&Options{}).
					Where("option_num = ?", optionNum).
					Update("body", optionLabel).Error
				if err != nil {
					return fmt.Errorf("failed to update option: %w", err)
				}
			}
		} else {
			createOptions = append(createOptions, Options{
				QuestionID: questionID,
				OptionNum:  optionNum,
				Body:       optionLabel,
			})
		}
	}

	if len(createOptions) > 0 {
		err := db.
			Session(&gorm.Session{}).
			Create(&createOptions).Error
		if err != nil {
			return fmt.Errorf("failed to create option: %w", err)
		}
	}

	if isDelete {
		err = db.
			Where("question_id = ? AND option_num > ?", questionID, len(options)).
			Delete(Options{}).Error
		if err != nil {
			return fmt.Errorf("failed to update option: %w", err)
		}
	}

	return nil
}

// DeleteOptions 選択肢の削除
func (*Option) DeleteOptions(ctx context.Context, questionID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("question_id = ?", questionID).
		Delete(Options{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete option: %w", err)
	}

	return nil
}

// GetOptions 質問の選択肢の取得
func (*Option) GetOptions(ctx context.Context, questionIDs []int) ([]Options, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	type option struct {
		QuestionID int         `gorm:"type:int(11) NOT NULL;"`
		Body       null.String `gorm:"type:text;default:NULL;"`
	}
	options := []option{}

	err = db.
		Where("question_id IN (?)", questionIDs).
		Order("question_id, option_num").
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
