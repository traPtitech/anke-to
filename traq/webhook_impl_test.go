package traq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageLimit(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 2000, MessageLimit)
}
