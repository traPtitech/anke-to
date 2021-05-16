package router

import (
	"crypt/rand"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
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
}

func newOAuth2() *OAuth2 {
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
	}
}

// Callback GET /oauth2/callbackの処理部分
func (o *OAuth2) Callback(code string, c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return fmt.Errorf("Failed In Getting Session: %w", err)
	}

	interfaceCodeVerifier, ok := sess.Values["codeVerifier"]
	if !ok || interfaceCodeVerifier == nil {
		return errors.New("CodeVerifier IS NULL")
	}
	codeVerifier := interfaceCodeVerifier.(string)

	res, err := o.getAccessToken(code, codeVerifier)
	if err != nil {
		return fmt.Errorf("Failed In Getting AccessToken:%w", err)
	}

	sess.Values["accessToken"] = res.AccessToken
	sess.Values["refreshToken"] = res.RefreshToken

	user, err := o.oauth.GetMe(res.AccessToken)
	if err != nil {
		return fmt.Errorf("Failed In Getting Me: %w", err)
	}

	sess.Values["userID"] = user.Id
	sess.Values["userName"] = user.Name

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return fmt.Errorf("Failed In Save Session: %w", err)
	}

	return nil
}

// GetGeneratedCode POST /oauth2/generate/codeの処理部分
func (o *OAuth2) GetGeneratedCode(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return nil, fmt.Errorf("Failed In Getting Session: %w", err)
	}

	state, err := randomString(60)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create a random string: %w", err))
	}

	codeVerifier, err := randomString(43)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create a random string: %w", err))
	}

	bytesCodeChallenge := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64url.Encode(bytesCodeChallenge[:])

	codeChallengeMethodOption := oauth2.SetAuthURLParam("code_challenge_method", "S256")
	codeChallengeOption := oauth2.SetAuthURLParam("code_challenge", codeChallenge)
	authURL := o.conf.AuthCodeURL(state)

	sess.Values["codeVerifier"] = codeVerifier

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return nil, fmt.Errorf("Failed In Save Session: %w", err)
	}

	return pkceParams, nil
}

// PostLogout POST /oauth2/logoutの処理部分
func (o *OAuth2) PostLogout(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return fmt.Errorf("Failed In Getting Session: %w", err)
	}

	interfaceAccessToken, ok := sess.Values["accessToken"]
	if !ok || interfaceAccessToken == nil {
		log.Printf("error: Unexpected No Access Token")
		return errors.New("No Access Token")
	}
	accessToken, ok := interfaceAccessToken.(string)
	if !ok {
		return errors.New("Invalid Access Token")
	}

	path := o.oauth.BaseURL()
	path.Path += "/oauth2/revoke"
	form := url.Values{}
	form.Set("token", accessToken)
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", path.String(), reqBody)
	if err != nil {
		return fmt.Errorf("Failed In Making HTTP Request:%w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed In HTTP Request:%w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Failed In Getting Access Token:(Status:%d %s)", res.StatusCode, res.Status)
	}

	err = o.session.RevokeSession(c)
	if err != nil {
		return fmt.Errorf("Failed In Revoke Session: %w", err)
	}

	return nil
}

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), letterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		idx := int(cache & letterIdxMask)
		if idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

func (o *OAuth2) getAccessToken(code string, codeVerifier string) (*authResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", o.clientID)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	reqBody := strings.NewReader(form.Encode())
	path := o.oauth.BaseURL()
	path.Path += "/oauth2/token"
	req, err := http.NewRequest("POST", path.String(), reqBody)
	if err != nil {
		return &authResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return &authResponse{}, err
	}
	if res.StatusCode != 200 {
		return &authResponse{}, fmt.Errorf("Failed In Getting Access Token:(Status:%d %s)", res.StatusCode, res.Status)
	}

	authRes := &authResponse{}
	err = json.NewDecoder(res.Body).Decode(authRes)
	if err != nil {
		return &authResponse{}, fmt.Errorf("Failed In Parsing Json: %w", err)
	}
	return authRes, nil
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
