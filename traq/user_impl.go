package traq

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

type User struct {
}

type UserRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u *User) GetMyID(token *oauth2.Token) (string, error) {
	path := "https://q.trap.jp/api/v3/users/me"

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new req :%w", err)
	}

	token.SetAuthHeader(req)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get http res :%w", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to get http res :(status :%d): %w", res.StatusCode, res.Status)
	}

	user := &UserRes{}
	if err = json.NewDecoder(res.Body).Decode(user); err != nil {
		return "", fmt.Errorf("failed to decode res:%w", err)
	}

	return user.Name, nil
}

func NewUser() *User {
	return &User{}
}
