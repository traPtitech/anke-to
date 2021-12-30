package router

import (
	"fmt"
	"github.com/traPtitech/anke-to/router/session"
	"golang.org/x/oauth2"
	"os"
)

var (
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	baseURL = "https://q.trap.jp/api/v3"
)

type Oauth struct {
	config *oauth2.Config
	sessStore session.Store
}

func NewOauth(sessStore session.Store) *Oauth {
	return &Oauth{
		config:    &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     oauth2.Endpoint{
				AuthURL:   fmt.Sprintf("%s/%s",baseURL,"oauth2/authorize"),
				TokenURL:  fmt.Sprintf("%s/%s",baseURL,"oauth2/token"),
			},
			Scopes:       []string{"read"},
		},
		sessStore: sessStore,
	}
}

