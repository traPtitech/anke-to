package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func timeConvert(time mysql.NullTime) string {
	if time.Valid {
		return time.Time.String()
	} else {
		return "NULL"
	}
}

func getUserID(c echo.Context) string {
	return c.Request().Header.Get("X-Showcase-User")
}
