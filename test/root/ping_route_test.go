package root_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/internal/test/test_helpers"

	"github.com/stretchr/testify/assert"
)

func TestPingRouteReturnsPlainTextPong(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(test_helpers.MockCfg)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
