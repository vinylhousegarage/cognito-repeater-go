package login

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type mockLoginURLProvider struct{}

func (m *mockLoginURLProvider) GetLoginURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/oauth2/authorize", nil
}

func TestLoginHandlerRedirectsToLoginEndpoint(t *testing.T) {
	t.Parallel()

	handler := LoginHandler(&mockLoginURLProvider{}, &testhelpers.MockMetadataURLProvider{})

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler(w, req)
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

	handler := LoginHandler(&mockLoginURLProvider{}, &testhelpers.MockMetadataURLProvider{})
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler(w, req)
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
