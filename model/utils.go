package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func TimeConvert(time mysql.NullTime) string {
	if time.Valid {
		return time.Time.String()
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
