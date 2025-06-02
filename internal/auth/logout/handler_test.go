package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type MockLogoutConfig struct{}

func (m *MockLogoutConfig) MetadataURL() string { return "https://mock.metadata.url" }

func TestNewLogoutHandlerRedirectsToLogoutEndpoint(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, _ = rec.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
			return rec.Result(), nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	handler := NewLogoutHandler(&MockLogoutConfig{}, mockClient)

	handler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
