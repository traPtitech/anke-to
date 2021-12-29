package session

import "github.com/labstack/echo/v4"

type IStore interface {
	GetMiddleware() echo.MiddlewareFunc
	GetSession(c echo.Context) (*Session,error)
}
