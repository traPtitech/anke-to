package router

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/router/session"
	"golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	baseURL      = "https://q.trap.jp/api/v3"
)

type Oauth struct {
	config    *oauth2.Config
	sessStore session.Store
}

func NewOauth(sessStore session.Store) *Oauth {
	return &Oauth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("%s/%s", baseURL, "oauth2/authorize"),
				TokenURL: fmt.Sprintf("%s/%s", baseURL, "oauth2/token"),
			},
			Scopes: []string{"read"},
		},
		sessStore: sessStore,
	}
}

func (o *Oauth) GetCode(c echo.Context) error {
	sess, err := o.sessStore.GetSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session:%w", err))
	}

	verifier := RandomString(90)
	sess.SetVerifier(verifier)

	state := RandomString(32)
	sess.SetState(state)

	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	challengeOption := oauth2.SetAuthURLParam("code_challenge", challenge)
	methodOption := oauth2.SetAuthURLParam("code_challenge_method", "s256")

	authURL := o.config.AuthCodeURL(state, challengeOption, methodOption)

	if err = sess.Save(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to save session:%w", err))
	}

	return c.String(http.StatusOK, authURL)
}

func (o *Oauth) Callback(c echo.Context) error {
	sess, err := o.sessStore.GetSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session :%w", err))
	}

	state := c.QueryParam("state")
	sessState, err := sess.GetState()
	if err != nil {
		if errors.Is(err, session.ErrNoValue) {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("failed to get state for no value  : %w", err))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get state :%w", err))
	}

	if state != sessState {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to match state")
	}

	code := c.QueryParam("code")
	verifier, err := sess.GetVerifier()
	if err != nil {
		if errors.Is(err, session.ErrNoValue) {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("failed to get verifier for no value:%w", err))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get verifier :%w", err))
	}

	codeChallengeOption := oauth2.SetAuthURLParam("code_verifier", verifier)

	token, err := o.config.Exchange(c.Request().Context(), code, codeChallengeOption)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to exchange token: %w", err))
	}

	sess.SetToken(token)

	if err = sess.Save(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to save :%w ", err))
	}

	return c.NoContent(http.StatusOK)
}
