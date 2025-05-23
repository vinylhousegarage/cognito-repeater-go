package callback

import (
	"net/url"
	"testing"

	"cognito-repeater-go/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildTokenRequestBody(t *testing.T) {
	cfg := &config.Config{
		UserPoolClientID: "abc123",
		RedirectURI:      "https://example.com/callback",
	}
	code := "xyz789"

	result := BuildTokenRequestBody(code, cfg)

	values, err := url.ParseQuery(result)
	require.NoError(t, err)

	assert.Equal(t, "authorization_code", values.Get("grant_type"))
	assert.Equal(t, "xyz789", values.Get("code"))
	assert.Equal(t, "abc123", values.Get("client_id"))
	assert.Equal(t, "https://example.com/callback", values.Get("redirect_uri"))
}
