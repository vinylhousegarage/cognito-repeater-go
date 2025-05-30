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

type mockLoginHandlerProvider struct {
	ClientID    string
	RedirectURI string
	Scope       string
}

func (m *mockLoginHandlerProvider) UserPoolClientIDValue() string { return m.ClientID }
func (m *mockLoginHandlerProvider) RedirectURIValue() string      { return m.RedirectURI }
func (m *mockLoginHandlerProvider) ScopeValue() string            { return m.Scope }

func TestBuildLoginURL_Success(t *testing.T) {
	provider := &mockLoginHandlerProvider{
		ClientID:    "example-client-id",
		RedirectURI: "https://example.com/callback",
		Scope:       "openid",
	}

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
	provider := &mockLoginHandlerProvider{}
	_, err := BuildLoginURL(provider, "://invalid-url", "state")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing protocol")
}
