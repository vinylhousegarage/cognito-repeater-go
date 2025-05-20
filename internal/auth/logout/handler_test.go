package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/config"

	"github.com/stretchr/testify/assert"
)

type mockEndpointProvider struct{}

func (m *mockEndpointProvider) GetLogoutURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/logout", nil
}

type mockMetadataProvider struct{}

func (m *mockMetadataProvider) MetadataURL() string {
	return "https://mock.metadata.url"
}

func TestLogoutHandlerRedirectsToLogoutEndpoint(t *testing.T) {
	t.Parallel()

	handler := LogoutHandler(&mockEndpointProvider{})(&mockMetadataProvider{})

	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
