package health_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/health"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestHealthRouteReturnsJSONHealthy(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(
		testhelpers.NewMockRouteDependencies(),
		testhelpers.NewMockHTTPClientOK(),
		testhelpers.MockLogger,
	)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res health.HealthResponse
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", res.Status)
}
