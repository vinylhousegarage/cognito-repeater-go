package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/test/test_helpers"

	"github.com/stretchr/testify/assert"
)

type mockLogoutEndpointProvider struct{}

func (m *mockLogoutEndpointProvider) GetLogoutURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/logout", nil
}

func TestLogoutRouteIsRegistered(t *testing.T) {
	t.Parallel()

	handler := logout.LogoutHandler(&mockLogoutEndpointProvider{}, &test_helpers.MockMetadataProvider{})

	router := http.NewServeMux()
	router.Handle("/logout", handler)

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
