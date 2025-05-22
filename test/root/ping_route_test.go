package root_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestPingRouteReturnsPlainTextPong(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(testhelpers.MockCfg, http.DefaultClient)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
