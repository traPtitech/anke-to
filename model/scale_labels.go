package model

import (
	"net/http"

	"github.com/labstack/echo"
)

type ScaleLabels struct {
	ID              int    `json:"questionID" db:"question_id"`
	ScaleLabelRight string `json:"scale_label_right" db:"scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"  db:"scale_label_left"`
	ScaleMin        int    `json:"scale_min" db:"scale_min"`
	ScaleMax        int    `json:"scale_max" db:"scale_max"`
}

func GetScaleLabels(c echo.Context, questionID int) (ScaleLabels, error) {
	scalelabel := ScaleLabels{}
	if err := db.Get(&scalelabel, "SELECT * FROM scale_labels WHERE question_id = ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return ScaleLabels{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return scalelabel, nil
}

func InsertScaleLabels(c echo.Context, lastID int, label ScaleLabels) error {
	if _, err := db.Exec(
		"INSERT INTO scale_labels (question_id, scale_label_left, scale_label_right, scale_min, scale_max) VALUES (?, ?, ?, ?, ?)",
		lastID, label.ScaleLabelLeft, label.ScaleLabelRight, label.ScaleMin, label.ScaleMax); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func UpdateScaleLabels(c echo.Context, questionID int, label ScaleLabels) error {
	if _, err := db.Exec(
		`INSERT INTO scale_labels (question_id, scale_label_right, scale_label_left, scale_min, scale_max) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE scale_label_right = ?, scale_label_left = ?, scale_min = ?, scale_max = ?`,
		questionID,
		label.ScaleLabelRight, label.ScaleLabelLeft, label.ScaleMin, label.ScaleMax,
		label.ScaleLabelRight, label.ScaleLabelLeft, label.ScaleMin, label.ScaleMax); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func DeleteScaleLabels(c echo.Context, questionID int) error {
	if _, err := db.Exec(
		"DELETE FROM scale_labels WHERE question_id= ?",
		questionID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
