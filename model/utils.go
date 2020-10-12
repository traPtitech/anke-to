package model

import (
	"time"

	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

// NullTimeToString null許容の時間のStringへの変換
func NullTimeToString(t null.Time) string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	}

	return "null"
}

// GetUserID ユーザーIDの取得
func GetUserID(c echo.Context) string {
	res := c.Request().Header.Get("X-Showcase-User")
	// test用
	if res == "" {
		return "mds_boy"
	}

	return res
}
