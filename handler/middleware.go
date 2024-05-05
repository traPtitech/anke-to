package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	validatorKey       = "validator"
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

// SetUserIDMiddleware X-Showcase-UserからユーザーIDを取得しセットする
func SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Showcase-User")
		if userID == "" {
			userID = "mds_boy"
		}

		c.Set(userIDKey, userID)

		return next(c)
	}
}

// getUserID ユーザーIDを取得する
func getUserID(c echo.Context) (string, error) {
	rowUserID := c.Get(userIDKey)
	userID, ok := rowUserID.(string)
	if !ok {
		return "", errors.New("invalid context userID")
	}

	return userID, nil
}
