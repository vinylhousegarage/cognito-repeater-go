package root

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataHandler_StatusAndBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metadata", nil)
	w := httptest.NewRecorder()

	MetadataHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

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
