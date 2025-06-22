package health

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestNewHealthHandler_ReturnsHealthy(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	NewHealthHandler(testhelpers.MockLogger)(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res HealthResponse
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", res.Status)
}
