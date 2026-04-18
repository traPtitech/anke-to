package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestV3_1TitleColumnType(t *testing.T) {
	var charMaxLength int
	err := db.Raw(
		"SELECT CHARACTER_MAXIMUM_LENGTH FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'questionnaires' AND COLUMN_NAME = 'title'",
	).Scan(&charMaxLength).Error
	require.NoError(t, err)
	assert.Equal(t, 1024, charMaxLength)

	// verify a 1024-char title can be inserted
	q := Questionnaires{
		Title:       strings.Repeat("a", 1024),
		Description: "test",
		ResTimeLimit: null.Time{},
		ResSharedTo: "public",
		IsPublished: true,
	}
	err = db.Create(&q).Error
	assert.NoError(t, err)
	if err == nil {
		db.Delete(&q)
	}
}
