package session

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
	"github.com/traPtitech/anke-to/model"
	"golang.org/x/oauth2"
)

type SessionStore struct {
	store *mysqlstore.MySQLStore
}

func NewSessionStore(sess model.ISession) (*SessionStore, error) {
	store, err := sess.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &SessionStore{
		store: store,
	}, nil
}

func (ss *SessionStore) GetMiddleware() echo.MiddlewareFunc {
	return session.Middleware(ss.store)
}

func (ss *SessionStore) GetSession(c echo.Context) (ISession, error) {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &Session{
		c: c,
		sess: sess,
	}, nil
}

type Session struct {
	c echo.Context
	store mysqlstore.MySQLStore
	sess *sessions.Session
}

func (s *Session) SetState(state string) error {
	s.sess.Values["state"] = state

	return nil
}

func (s *Session) GetState() (string, error) {
	iState, ok := s.sess.Values["state"]
	if !ok || iState == nil {
		return "", ErrNoValue
	}

	return iState.(string), nil
}

func (s *Session) SetCodeVerifier(codeVerifier string) error {
	s.sess.Values["codeVerifier"] = codeVerifier

	return nil
}

func (s *Session) GetCodeVerifier() (string, error) {
	iCodeVerifier, ok := s.sess.Values["codeVerifier"]
	if !ok || iCodeVerifier == nil {
		return "", ErrNoValue
	}

	return iCodeVerifier.(string), nil
}

func (s *Session) SetToken(token *oauth2.Token) error {
	s.sess.Values["token"] = token

	return nil
}

func (s *Session) GetToken() (*oauth2.Token, error) {
	iToken, ok := s.sess.Values["token"]
	if !ok || iToken == nil {
		return nil, ErrNoValue
	}

	return iToken.(*oauth2.Token), nil
}

func (s *Session) Save() error {
	err := s.sess.Save(s.c.Request(), s.c.Response())
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (s *Session) Revoke() error {
	err := s.store.Delete(s.c.Request(), s.c.Response(), s.sess)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
