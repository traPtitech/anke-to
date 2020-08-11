package model

import (
	"database/sql"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

func NullTimeToString(t null.Time) string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	} else {
		return "NULL"
	}
}

func NullStringConvert(str sql.NullString) string {
	if str.Valid {
		return str.String
	} else {
		return "NULL"
	}
}

func GetUserID(c echo.Context) string {
	res := c.Request().Header.Get("X-Showcase-User")
	// test用
	if res == "" {
		return "mds_boy"
	}
	return res
}
