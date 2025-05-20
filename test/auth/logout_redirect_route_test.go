package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/internal/test/test_helpers"

	"github.com/stretchr/testify/assert"
)

func TestRouter_LogoutRedirectRoute_ReturnsExpectedJSON(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(test_helpers.MockCfg)

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
