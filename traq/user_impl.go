package traq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type User struct {}

func NewUser() *User {
	return &User{}
}

type UserRes struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func (*User) GetMyID(token *oauth2.Token) (string, error) {
	path := fmt.Sprintf("%s%s", baseURL, "/users/me")

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new request: %w", err)
	}

	token.SetAuthHeader(req)
	httpClient := http.DefaultClient

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to http request: %w", err)
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	user := &UserRes{}
	err = json.NewDecoder(res.Body).Decode(user)
	if err != nil {
		return "", err
	}

	return user.Name, nil
}
