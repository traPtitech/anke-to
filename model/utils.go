package model

import (
	"errors"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

var (
	// ErrTooLargePageNum too large page number
	ErrTooLargePageNum = errors.New("too large page number")
	// ErrInvalidRegex invalid regexp
	ErrInvalidRegex = errors.New("invalid regexp")
	// ErrNoRecordUpdated no record updated
	ErrNoRecordUpdated = errors.New("no record updated")
	// ErrNoRecordDeleted no record deleted
	ErrNoRecordDeleted = errors.New("no record deleted")
	// ErrInvalidSortParam invalid sort param
	ErrInvalidSortParam = errors.New("invalid sort type")
	// ErrInvalidNumber MinBound,MaxBoundの指定が有効ではない
	ErrInvalidNumber = errors.New("invalid number")
	// ErrNumberBoundary MinBound <= value <= MaxBound でない
	ErrNumberBoundary = errors.New("the number is out of bounds")
	// ErrTextMatching RegexPatternにマッチしていない
	ErrTextMatching = errors.New("failed to match the pattern")
	// ErrInvalidAnsweredParam invalid sort param
	ErrInvalidAnsweredParam = errors.New("invalid answered param")
)

// NullTimeToString null許容の時間のStringへの変換
func NullTimeToString(t null.Time) string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	}

	return "null"
}

// GetUserID ユーザーIDの取得
func GetUserID(c echo.Context) string {
	res := c.Request().Header.Get("X-Showcase-User")
	// test用
	if res == "" {
		return "mds_boy"
	}

	return res
}
