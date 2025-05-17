package root_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"

	"github.com/stretchr/testify/assert"
)

func TestRouter_MetadataRoute_ReturnsExpectedJSON(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/metadata", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := map[string]string{
		"login_endpoint":               "/login",
		"logout_endpoint":              "/logout",
		"verify_access_token_endpoint": "/token",
		"verify_userinfo_endpoint":     "/user",
		"health_check_endpoint":        "/ping",
		"simulate_404_endpoint":        "/error/404",
		"openapi_url":                  "/openapi.json",
	}

	var actual map[string]string
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
