package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/test/test_helpers"

	"github.com/stretchr/testify/assert"
)

type mockLoginEndpointProvider struct{}

func (m *mockLoginEndpointProvider) GetLoginURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/oauth2/authorize", nil
}

func TestLoginRouteIsRegistered(t *testing.T) {
	t.Parallel()

	handler := login.LoginHandler(&mockLoginEndpointProvider{}, &test_helpers.MockMetadataProvider{})

	router := http.NewServeMux()
	router.Handle("/login", handler)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/oauth2/authorize", resp.Header.Get("Location"))
}
