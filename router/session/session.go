//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package session

import "github.com/labstack/echo/v4"

type IStore interface {
	GetMiddleware() echo.MiddlewareFunc
	GetSession(c echo.Context) (*Session,error)
}
