package notfound

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestNewError404Handler_StatusAndBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/error/404", nil)
	w := httptest.NewRecorder()

	NewError404Handler(testhelpers.MockLogger)(w, req)

	resp := w.Result()

	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var got ErrorSimulationResponse
	err := json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, "Simulated 404 Not Found", got.Message)
}
