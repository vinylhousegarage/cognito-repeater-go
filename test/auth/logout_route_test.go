package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/logout"
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

func TestLogoutRouteIsRegistered(t *testing.T) {
	t.Parallel()

	handler := logout.LogoutHandler(&mockEndpointProvider{}, &mockMetadataProvider{})

	router := http.NewServeMux()
	router.Handle("/logout", handler)

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
}
