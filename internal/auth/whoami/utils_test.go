package whoami

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAuthHeaderToken_Success(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer test.token.value")

	token, err := ExtractAuthHeaderToken(req)
	assert.NoError(t, err)
	assert.Equal(t, "test.token.value", token)
}

func TestExtractAuthHeaderToken_MissingHeader(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/me", nil)

	_, err := ExtractAuthHeaderToken(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authorization header is missing")
}

func TestExtractAuthHeaderToken_InvalidFormat(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Basic abc123")

	_, err := ExtractAuthHeaderToken(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid authorization header format")
}
