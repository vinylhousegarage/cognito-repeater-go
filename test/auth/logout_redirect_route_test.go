package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"

	"github.com/stretchr/testify/assert"
)

var mockCfg = &config.Config{
	Region:           "ap-northeast-1",
	ClientSecret:     "client-secret",
	LogoutURI:        "https://example.com/logout",
	RedirectURI:      "https://localhost/callback",
	Scope:            "openid",
	UserPoolClientID: "client-id",
	UserPoolID:       "pool-id",
}

func TestRouter_LogoutRedirectRoute_ReturnsExpectedJSON(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(mockCfg)
	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := map[string]string{
		"message": "Logout successful",
	}

	var actual map[string]string
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err, "failed to unmarshal response JSON")

	assert.Equal(t, expected, actual, "unexpected JSON response")
}
