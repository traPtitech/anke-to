package model

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

//Targets targetsテーブルの構造体
type Targets struct {
	QuestionnaireID int    `gorm:"type:int(11);primary_key;NOT NULL;"`
	UserTraqID      string `gorm:"column:user_traqid;type:char(30);primary_key;NOT NULL;"`
}

// InsertTargets アンケートの対象を追加
func InsertTargets(c echo.Context, questionnaireID int, targets []string) error {
	rowTargets := make([]interface{}, 0, len(targets))
	for _, target := range targets {
		rowTargets = append(rowTargets, Targets{
			QuestionnaireID: questionnaireID,
			UserTraqID:      target,
		})
	}

	err := gormbulk.BulkInsert(gormDB, rowTargets, len(rowTargets))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to insert target: %w", err))
	}

	return nil
}

// DeleteTargets アンケートの対象を削除
func DeleteTargets(c echo.Context, questionnaireID int) error {
	err := gormDB.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&Targets{}).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete targets: %w", err))
	}

	return nil
}
