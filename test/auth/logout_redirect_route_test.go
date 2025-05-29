package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestRouterLogoutRedirectRouteReturnsExpectedJSON(t *testing.T) {
	t.Parallel()

	handlerDeps := deps.HandlerDependencies{
		Config:     testhelpers.MockCfg,
		HTTPClient: http.DefaultClient,
	}

	r := router.NewRouter(handlerDeps.Config, handlerDeps.Config, handlerDeps.Config, handlerDeps.Config, handlerDeps.HTTPClient)

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
