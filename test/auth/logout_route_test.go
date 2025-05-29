package auth_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type MockLogoutConfig struct{}

func (m *MockLogoutConfig) Audience() string              { return "mock-audience" }
func (m *MockLogoutConfig) ClientSecretValue() string     { return "mock-secret" }
func (m *MockLogoutConfig) Issuer() string                { return "https://mock-issuer" }
func (m *MockLogoutConfig) MetadataURL() string           { return "https://mock.metadata.url" }
func (m *MockLogoutConfig) RedirectURIValue() string      { return "https://example.com/callback" }
func (m *MockLogoutConfig) ScopeValue() string            { return "openid" }
func (m *MockLogoutConfig) UserPoolClientIDValue() string { return "mock-client-id" }

func TestLogoutRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"end_session_endpoint": "https://example.com/logout"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		},
	}

	r := router.NewRouter(&MockLogoutConfig{}, &MockLogoutConfig{}, mockClient)

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Location"), "logout")
}
