package me

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKIDFromToken(t *testing.T) {
	t.Parallel()

	t.Run("valid token with kid", func(t *testing.T) {
		header := `{"kid":"example-kid"}`
		encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
		token := encodedHeader + ".payload.signature"

		kid, err := ExtractKIDFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "example-kid", kid)
	})

	t.Run("invalid JWT format", func(t *testing.T) {
		token := "just.onepart"
		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JWT format")
	})

	t.Run("invalid base64 in header", func(t *testing.T) {
		token := "!!invalid!!base64.payload.signature"
		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode JWT header")
	})

	t.Run("invalid JSON in header", func(t *testing.T) {
		invalidJSON := "{bad json}"
		encoded := base64.RawURLEncoding.EncodeToString([]byte(invalidJSON))
		token := encoded + ".payload.signature"

		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse JWT header JSON")
	})

	t.Run("missing kid in header", func(t *testing.T) {
		header := `{"alg":"RS256"}`
		encoded := base64.RawURLEncoding.EncodeToString([]byte(header))
		token := encoded + ".payload.signature"

		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kid not found in JWT header")
	})
}
