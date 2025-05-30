package login

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type mockLoginHandlerProvider struct {
	ClientID    string
	RedirectURI string
	Scope       string
	Metadata    string
}

func (m *mockLoginHandlerProvider) MetadataURL() string           { return m.Metadata }
func (m *mockLoginHandlerProvider) UserPoolClientIDValue() string { return m.ClientID }
func (m *mockLoginHandlerProvider) RedirectURIValue() string      { return m.RedirectURI }
func (m *mockLoginHandlerProvider) ScopeValue() string            { return m.Scope }

func newMockLoginHandlerProvider() *mockLoginHandlerProvider {
	return &mockLoginHandlerProvider{
		ClientID:    "example-client-id",
		RedirectURI: "https://example.com/callback",
		Scope:       "openid",
		Metadata:    "https://mock.example.com/.well-known/openid-configuration",
	}
}

func newMockHTTPClient() *testhelpers.MockHTTPClient {
	return &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, _ = rec.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
			return rec.Result(), nil
		},
	}
}

func TestLoginHandler_RedirectsToLoginEndpoint(t *testing.T) {
	t.Parallel()

	p := newMockLoginHandlerProvider()
	c := newMockHTTPClient()

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	LoginHandler(p, c)(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	parsed, err := url.Parse(location)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", parsed.Host)
	assert.Equal(t, "/oauth2/authorize", parsed.Path)

	state := parsed.Query().Get("state")
	assert.NotEmpty(t, state, "state param should be present in the redirect URL")
}

func TestLoginHandlerSetsStateCookie(t *testing.T) {
	t.Parallel()

	p := newMockLoginHandlerProvider()
	c := newMockHTTPClient()

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	LoginHandler(p, c)(w, req)

	resp := w.Result()

	cookies := resp.Cookies()
	var stateCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "oauth_state" {
			stateCookie = c
			break
		}
	}

	assert.NotNil(t, stateCookie, "oauth_state cookie should be set")
	assert.NotEmpty(t, stateCookie.Value)
	assert.True(t, stateCookie.HttpOnly)
	assert.True(t, stateCookie.Secure)
	assert.Equal(t, "/", stateCookie.Path)
	assert.Equal(t, http.SameSiteLaxMode, stateCookie.SameSite)
}
