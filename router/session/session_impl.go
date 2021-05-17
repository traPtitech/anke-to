package session

import (
	"fmt"
	"time"

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

func (s *Session) SetUserID(userID string) error {
	s.sess.Values["userID"] = userID

	return nil
}

func (s *Session) GetUserID() (string, error) {
	iUserID, ok := s.sess.Values["userID"]
	if !ok || iUserID == nil {
		return "", ErrNoValue
	}

	return iUserID.(string), nil
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
	s.sess.Values["access_token"] = token.AccessToken
	s.sess.Values["token_type"] = token.TokenType
	s.sess.Values["refresh_token"] = token.RefreshToken
	s.sess.Values["expiry"] = token.Expiry

	return nil
}

func (s *Session) GetToken() (*oauth2.Token, error) {
	iAccessToken, ok := s.sess.Values["access_token"]
	if !ok || iAccessToken == nil {
		return nil, ErrNoValue
	}

	iTokenType, ok := s.sess.Values["token_type"]
	if !ok || iTokenType == nil {
		return nil, ErrNoValue
	}

	iRefreshToken, ok := s.sess.Values["refresh_token"]
	if !ok || iRefreshToken == nil {
		return nil, ErrNoValue
	}

	iExpiry, ok := s.sess.Values["expiry"]
	if !ok || iExpiry == nil {
		return nil, ErrNoValue
	}

	return &oauth2.Token{
		AccessToken: iAccessToken.(string),
		TokenType: iTokenType.(string),
		RefreshToken: iRefreshToken.(string),
		Expiry: iExpiry.(time.Time),
	}, nil
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
