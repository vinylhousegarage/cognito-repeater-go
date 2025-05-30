package root_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestPingRouteReturnsPlainTextPong(t *testing.T) {
	t.Parallel()

	handlerDeps := deps.HandlerDependencies{
		Config:     testhelpers.MockCfg,
		HTTPClient: http.DefaultClient,
	}

	r := router.NewRouter(handlerDeps.Config, handlerDeps.Config, handlerDeps.Config, handlerDeps.Config, handlerDeps.Config, handlerDeps.HTTPClient)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
