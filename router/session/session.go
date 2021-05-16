package session

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type ISessionStore interface {
	SetMiddleware() echo.MiddlewareFunc
	GetSession(c echo.Context) (ISession, error)
}

type ISession interface {
	SetState(string) error
	GetState() (string, error)
	SetCodeVerifier(string) error
	GetCodeVerifier() (string, error)
	SetToken(*oauth2.Token) error
	GetToken() (*oauth2.Token, error)
	Save() error
	Revoke() error
}
