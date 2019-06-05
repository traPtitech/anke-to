package router

import (
	"net/http"

	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

// GetUsersMe GET /users/me
func GetUsersMe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"traqID": model.GetUserID(c),
	})
}
