package model

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

// Options optionsテーブルの構造体
type Options struct {
	ID         int    `gorm:"type:int(11) AUTO_INCREMENT NOT NULL PRIMARY KEY;"`
	QuestionID int    `gorm:"type:int(11) NOT NULL;"`
	OptionNum  int    `gorm:"type:int(11) NOT NULL;"`
	Body       string `gorm:"type:text;default:NULL;"`
}

// InsertOption 選択肢の追加
func InsertOption(c echo.Context, lastID int, num int, body string) error {
	option := Options{
		QuestionID: lastID,
		OptionNum:  num,
		Body:       body,
	}
	err := db.Create(&option).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateOptions 選択肢の修正
func UpdateOptions(c echo.Context, options []string, questionID int) error {
	var err error
	for i, optionLabel := range options {
		option := Options{
			Body: optionLabel,
		}
		err = db.
			Model(Options{}).
			Where("question_id = ? AND option_num = ?", questionID, i).
			Update(&option).Error
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	err = db.Where("question_id = ? AND option_num > ?", questionID, len(options)).Delete(Options{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// DeleteOptions 選択肢の削除
func DeleteOptions(c echo.Context, questionID int) error {
	err := db.
		Where("question_id = ?", questionID).
		Delete(Options{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// GetOptions 質問の選択肢の取得
func GetOptions(c echo.Context, questionID int) ([]string, error) {
	bodies := []null.String{}

	err := db.
		Model(Options{}).
		Where("question_id = ?", questionID).
		Order("option_num").
		Pluck("body", &bodies).Error
	if err != nil {
		c.Logger().Error(err)
		return []string{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	options := make([]string, 0, len(bodies))
	for _, option := range bodies {
		if option.Valid {
			options = append(options, option.String)
		} else {
			options = append(options, "")
		}
	}

	return options, nil
}
