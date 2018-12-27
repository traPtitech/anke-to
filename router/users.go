package router

import (
	"net/http"

	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

func GetUsersMe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"traqID": model.GetUserID(c),
	})
}
