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

func TestBuildLoginURL_Success(t *testing.T) {
	provider := newMockLoginHandlerProvider()

	endpoint := "https://auth.example.com/oauth2/authorize"
	state := "sample-state-value"

	result, err := BuildLoginURL(provider, endpoint, state)
	assert.NoError(t, err)

	parsed, err := url.Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, "auth.example.com", parsed.Host)
	assert.Equal(t, "/oauth2/authorize", parsed.Path)

	queries := parsed.Query()
	assert.Equal(t, "code", queries.Get("response_type"))
	assert.Equal(t, provider.ClientID, queries.Get("client_id"))
	assert.Equal(t, provider.RedirectURI, queries.Get("redirect_uri"))
	assert.Equal(t, provider.Scope, queries.Get("scope"))
	assert.Equal(t, state, queries.Get("state"))
}

func TestBuildLoginURL_InvalidEndpoint(t *testing.T) {
	provider := newMockLoginHandlerProvider()
	_, err := BuildLoginURL(provider, "://invalid-url", "state")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing protocol")
}
