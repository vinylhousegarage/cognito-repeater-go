package health_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/health/ping"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestPingRouteReturnsPlainTextPong(t *testing.T) {
	t.Parallel()

	deps := testhelpers.NewMockRouteDependencies()
	cli := testhelpers.NewMockHTTPClientOK()
	mockLogger := zap.NewNop()

	r := router.NewRouter(deps, cli, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res ping.PingResponse
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)
	assert.Equal(t, "pong", res.Message)
}
