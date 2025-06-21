package logout

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type MockLogoutConfig struct{}

func (m *MockLogoutConfig) LogoutURIValue() string        { return "https://myapp.com/logout/redirect" }
func (m *MockLogoutConfig) MetadataURL() string           { return "https://mock.metadata.url" }
func (m *MockLogoutConfig) UserPoolClientIDValue() string { return "test-client-id" }

func TestNewLogoutHandlerRedirectsWithIDTokenHint(t *testing.T) {
	t.Parallel()

	const fakeIDToken = "eyJhbGciOi..."

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, _ = rec.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
			return rec.Result(), nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+fakeIDToken)

	w := httptest.NewRecorder()
	handler := NewLogoutHandler(&MockLogoutConfig{}, mockClient, testhelpers.MockLogger)
	handler.ServeHTTP(w, req)

	resp := w.Result()

	params := url.Values{}
	params.Set("client_id", "test-client-id")
	params.Set("post_logout_redirect_uri", "https://myapp.com/logout/redirect")
	params.Set("id_token_hint", fakeIDToken)

	expected := "https://example.com/logout?" + params.Encode()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, expected, resp.Header.Get("Location"))
}
