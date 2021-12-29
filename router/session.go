package router

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/thanhpk/randstr"
	"net/http"
	"net/url"
)

type OAuth struct {
	ClientID string
	Secret   string
}

func (o *OAuth) GetCallback(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	verifier := randstr.String(64)
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	sess.Values["verifier"] = verifier
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30,
		HttpOnly: true,
	}


	if err = sess.Save(c.Request(), c.Response());err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	oauthUrl,_ := url.Parse("https://q.trap.jp/api/v3/oauth2/authorize")
	q := oauthUrl.Query()
	oauthUrl.Query().Set("response_type","code")
	oauthUrl.Query().Set("client_id",o.ClientID)
	oauthUrl.Query().Set("code_challenge_method","s256")
	oauthUrl.Query().Set("code_challenge",challenge)
	oauthUrl.RawQuery = q.Encode()

	return c.JSON(http.StatusOK,oauthUrl.String())
}
