package router

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/router/session"
	"golang.org/x/oauth2"
)

var (
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	baseURL = "https://q.trap.jp/api/v3"
)

// OAuth2 oauthの構造体
type OAuth2 struct {
	conf *oauth2.Config
	sessStore session.ISessionStore
}

func newOAuth2(sessStore session.ISessionStore) *OAuth2 {
	return &OAuth2{
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("%s%s", baseURL, "/oauth2/auth"),
				TokenURL: fmt.Sprintf("%s%s", baseURL, "/oauth2/token"),
			},
    },
		sessStore: sessStore,
	}
}

// Callback GET /oauth2/callbackの処理部分
func (o *OAuth2) Callback(c echo.Context) error {
	sess, err := o.sessStore.GetSession(c)
	if errors.Is(err, session.ErrNoSession) {
		return echo.NewHTTPError(http.StatusUnauthorized, "no session")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session: %w", err))
	}

	reqState := c.QueryParam("state")

	sessState, err := sess.GetState()
	if errors.Is(err, session.ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "no state")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get state: %w", err))
	}

	if reqState != sessState {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid state")
	}

	code := c.QueryParam("code")

	codeVerifier, err := sess.GetCodeVerifier()
	if errors.Is(err, session.ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "no codeVerifier")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get codeVerifier: %w", err))
	}

	codeChallengeOption := oauth2.SetAuthURLParam("code_verifier", codeVerifier)
	token, err := o.conf.Exchange(c.Request().Context(), code, codeChallengeOption)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to exchange token: %w", err))
	}

	err = sess.SetToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to set token: %w", err))
	}

	err = sess.Save()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to save session: %w", err))
	}

	return c.NoContent(http.StatusOK)
}

// GetGeneratedCode POST /oauth2/generate/codeの処理部分
func (o *OAuth2) GetGeneratedCode(c echo.Context) error {
	sess, err := o.sessStore.GetSession(c)
	if errors.Is(err, session.ErrNoSession) {
		return echo.NewHTTPError(http.StatusUnauthorized, "no session")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session: %w", err))
	}

	state, err := randomString(60)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create a random string: %w", err))
	}

	codeVerifier, err := randomString(43)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create a random string: %w", err))
	}

	h := sha256.New()
	_, err = io.Copy(h, strings.NewReader(codeVerifier))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to copy codeVerifier: %w", err))
	}

	codeChallengeBuilder := strings.Builder{}

	bytesCodeChallenge := sha256.Sum256(h.Sum(nil))
	enc := base64.NewEncoder(base64.RawURLEncoding, &codeChallengeBuilder)
	defer enc.Close()

	_, err = enc.Write(bytesCodeChallenge[:])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to write codeChallenge: %w", err))
	}

	err = enc.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to encode codeChallenge: %w", err))
	}

	codeChallengeMethodOption := oauth2.SetAuthURLParam("code_challenge_method", "S256")
	codeChallengeOption := oauth2.SetAuthURLParam("code_challenge", codeChallengeBuilder.String())
	authURL := o.conf.AuthCodeURL(state, codeChallengeMethodOption, codeChallengeOption)

	err = sess.SetState(state)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to set state: %w", err))
	}

	err = sess.SetCodeVerifier(codeVerifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to set codeVerifier: %w", err))
	}

	err = sess.Save()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to save session: %w", err))
	}

	return c.String(http.StatusOK, authURL)
}

// PostLogout POST /oauth2/logoutの処理部分
func (o *OAuth2) PostLogout(c echo.Context) error {
	sess, err := o.sessStore.GetSession(c)
	if errors.Is(err, session.ErrNoSession) {
		return echo.NewHTTPError(http.StatusUnauthorized, "no session")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session: %w", err))
	}

	token, err := sess.GetToken()
	if errors.Is(err, session.ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "no token")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get token: %w", err))
	}

	path := fmt.Sprintf("%s%s", baseURL, "/oauth2/revoke")
	form := url.Values{}
	form.Set("token", token.AccessToken)
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed in HTTP request:%w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to revoke access token:(status:%d %s)", res.StatusCode, res.Status)
	}

	err = sess.Revoke()
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

func randomString(n int) (string, error) {
	bytesState := make([]byte,  n)
	_, err := rand.Read(bytesState)
	if err != nil {
    return "", fmt.Errorf("failed to read random bytes: %w", err)
  }

	sb := strings.Builder{}
  for _, v := range bytesState {
    // 制御文字が当たらないように調整
    _, err := sb.Write([]byte{v%byte(94) + 33})
		if err != nil {
			return "", fmt.Errorf("failed to write byte: %w", err)
		}
  }

  return sb.String(), nil
}
