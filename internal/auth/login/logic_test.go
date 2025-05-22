package login

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	const state = "abc123"
	c := BuildStateCookie(state)

	assert.Equal(t, "oauth_state", c.Name)
	assert.Equal(t, state, c.Value)
	assert.Equal(t, "/", c.Path)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
	assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
}

func TestBuildLoginURLSuccess(t *testing.T) {
	t.Parallel()

	endpoint := "https://example.com/oauth2/authorize"
	const state = "xyz789"

	result, err := BuildLoginURL(endpoint, state)
	assert.NoError(t, err)

	parsed, err := url.Parse(result)
	assert.NoError(t, err)

	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/oauth2/authorize", parsed.Path)

	q := parsed.Query()
	assert.Equal(t, state, q.Get("state"))
}

func TestBuildLoginURLInvalidURL(t *testing.T) {
	t.Parallel()

	const (
		invalidEndpoint = ":::"
		state           = "abc"
	)

	_, err := BuildLoginURL(invalidEndpoint, state)
	assert.Error(t, err)
}
