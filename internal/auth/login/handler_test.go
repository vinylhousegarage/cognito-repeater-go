package login

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func NewMockHTTPClientOKWithAuthEndpoint() httpclient.HTTPClient {
	return &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(
					`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`,
				)),
			}, nil
		},
	}
}

var mockProvider = &testhelpers.MockCfg
var mockClient = NewMockHTTPClientOKWithAuthEndpoint()
var mockLogger = zap.NewNop()

func TestNewLoginHandler_RedirectsToLoginEndpoint(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	NewLoginHandler(mockProvider, mockClient, mockLogger)(w, req)

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

func TestNewLoginHandlerSetsStateCookie(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	NewLoginHandler(mockProvider, mockClient, mockLogger)(w, req)

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

func TestNewLoginHandler_MetadataFetchFails(t *testing.T) {
	t.Parallel()

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	NewLoginHandler(&testhelpers.MockCfg, client, mockLogger)(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
