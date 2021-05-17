package session

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type IStore interface {
	GetMiddleware() echo.MiddlewareFunc
	GetSession(c echo.Context) (ISession, error)
}

type ISession interface {
	SetUserID(string) error
	GetUserID() (string, error)
	SetState(string) error
	GetState() (string, error)
	SetCodeVerifier(string) error
	GetCodeVerifier() (string, error)
	SetToken(*oauth2.Token) error
	GetToken() (*oauth2.Token, error)
	Save() error
	Revoke() error
}
