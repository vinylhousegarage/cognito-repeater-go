package login

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_RedirectsToLoginEndpoint(t *testing.T) {
	t.Parallel()

	d := deps.HandlerDependencies{
		Config:     &testhelpers.MockCfg,
		HTTPClient: testhelpers.NewMockHTTPClientOKWithAuthEndpoint(),
	}

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	login.LoginHandler(d)(w, req)
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

	d := deps.HandlerDependencies{
		Config:     &testhelpers.MockCfg,
		HTTPClient: testhelpers.NewMockHTTPClientOKWithAuthEndpoint(),
	}

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	login.LoginHandler(d)(w, req)
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
