package model

import (
	"net/http"

	"github.com/labstack/echo"
)

type Option struct {
	Id        int
	QuestionId int
	OptionNum  int
	Body       string
}

func GetOptions(c echo.Context, questionID int) ([]string, error) {
	bodies := []string{}

	err := db.
		Model(Option{}).
		Where("question_id = ?", questionID).
		Order("option_num").
		Pluck("body", &bodies).Error
	if err != nil {
		c.Logger().Error(err)
		return []string{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return bodies, nil
}

func InsertOption(c echo.Context, lastID int, num int, body string) error {
	option := Option{
		QuestionId: lastID,
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

func UpdateOptions(c echo.Context, options []string, questionID int) error {
	var err error
	option := Option{}
	for i, v := range options {
		err = db.Where(Option{QuestionId: questionID}).Assign(Option{OptionNum: i + 1, Body: v}).FirstOrCreate(&option).Error
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	err = db.Where("question_id = ? AND option_num > ?", questionID, len(options)).Delete(Option{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteOptions(c echo.Context, questionID int) error {
	err := db.Where("question_id = ?", questionID).Delete(Option{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
