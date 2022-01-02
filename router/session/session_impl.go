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

type Store struct {
	store *mysqlstore.MySQLStore
}

func (s *Store) GetMiddleware() echo.MiddlewareFunc {
	return session.Middleware(s.store)
}

func (s *Store) GetSession(c echo.Context) (*Session, error) {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return nil, fmt.Errorf("failed to get session:%w", err)
	}

	return &Session{
		c:    c,
		sess: sess,
	}, nil
}

func NewStore(sess model.Session) (*Store, error) {
	store, err := sess.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &Store{store: store}, nil
}

type Session struct {
	c    echo.Context
	sess *sessions.Session
}

func (s *Session) SetUserID(userID string) {
	s.sess.Values["userID"] = userID
}

func (s *Session) GetUserID() (string, error) {
	userID, ok := s.sess.Values["userID"].(string)
	if !ok || userID == "" {
		return "", ErrNoValue
	}

	return userID, nil
}

func (s *Session) SetVerifier(verifier string) {
	s.sess.Values["verifier"] = verifier
}

func (s *Session) GetVerifier() (string, error) {
	verifier, ok := s.sess.Values["verifier"].(string)
	if !ok || verifier == "" {
		return "", ErrNoValue
	}

	return verifier, nil
}

func (s *Session) SetToken(token *oauth2.Token) {
	s.sess.Values["access_token"] = token.AccessToken
	s.sess.Values["token_type"] = token.TokenType
	s.sess.Values["refresh_token"] = token.RefreshToken
	s.sess.Values["expiry"] = token.Expiry
}

func (s *Session) GetToken() (*oauth2.Token, error) {
	accessToken, ok := s.sess.Values["access_token"].(string)
	if !ok || accessToken == "" {
		return nil, ErrNoValue
	}

	tokenType, ok := s.sess.Values["token_type"].(string)
	if !ok || tokenType == "" {
		return nil, ErrNoValue
	}

	refreshToken, ok := s.sess.Values["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return nil, ErrNoValue
	}

	expiry, ok := s.sess.Values["expiry"].(time.Time)
	if !ok || expiry.IsZero() {
		return nil, ErrNoValue
	}

	return &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}, nil
}
