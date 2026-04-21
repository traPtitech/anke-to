package traq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageLimit(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 10000, MessageLimit)
}

func computeExpectedHMACSHA1(secret, message string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestCalcHMACSHA1(t *testing.T) {
	tests := []struct {
		description string
		secret      string
		message     string
	}{
		{
			description: "normal message with a known secret produces correct HMAC",
			secret:      "test-secret",
			message:     "Hello, traQ!",
		},
		{
			description: "empty message with a secret",
			secret:      "my-secret",
			message:     "",
		},
		{
			description: "empty secret with a message",
			secret:      "",
			message:     "some message",
		},
		{
			description: "both empty",
			secret:      "",
			message:     "",
		},
		{
			description: "long message",
			secret:      "s",
			message:     string(make([]byte, 65536)),
		},
		{
			description: "message with newlines and unicode characters",
			secret:      "secret-key",
			message:     "line1\nline2\nline3\u3042",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Setenv("TRAQ_WEBHOOK_SECRET", tt.secret)

			result, err := calcHMACSHA1(tt.message)
			require.NoError(t, err)

			assert.Len(t, result, 40, "HMAC-SHA1 hex output must be 40 characters")
			_, decodeErr := hex.DecodeString(result)
			assert.NoError(t, decodeErr, "result must be valid hex")

			expected := computeExpectedHMACSHA1(tt.secret, tt.message)
			assert.Equal(t, expected, result)
		})
	}
}

func TestCalcHMACSHA1_DifferentSecretsProduceDifferentMACs(t *testing.T) {
	const message = "important webhook payload"

	t.Setenv("TRAQ_WEBHOOK_SECRET", "secret-A")
	resultA, err := calcHMACSHA1(message)
	require.NoError(t, err)

	t.Setenv("TRAQ_WEBHOOK_SECRET", "secret-B")
	resultB, err := calcHMACSHA1(message)
	require.NoError(t, err)

	assert.NotEqual(t, resultA, resultB, "different secrets must produce different MACs")
}

func TestCalcHMACSHA1_DifferentMessagesProduceDifferentMACs(t *testing.T) {
	t.Setenv("TRAQ_WEBHOOK_SECRET", "fixed-secret")

	result1, err := calcHMACSHA1("message-1")
	require.NoError(t, err)

	result2, err := calcHMACSHA1("message-2")
	require.NoError(t, err)

	assert.NotEqual(t, result1, result2, "different messages must produce different MACs")
}

func TestCalcHMACSHA1_Deterministic(t *testing.T) {
	t.Setenv("TRAQ_WEBHOOK_SECRET", "stable-secret")

	r1, err := calcHMACSHA1("hello world")
	require.NoError(t, err)

	r2, err := calcHMACSHA1("hello world")
	require.NoError(t, err)

	assert.Equal(t, r1, r2, "same input must always produce the same HMAC")
}

func TestNewWebhook(t *testing.T) {
	t.Parallel()

	wh := NewWebhook()
	require.NotNil(t, wh)
}
