package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV3_1(t *testing.T) {
	err := v3_1().Migrate(db)
	assert.NoError(t, err)
}
