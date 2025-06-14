package simulatederrors_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/internal/simulatederrors/notfound"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestError404RouteReturnsJSONNotFound(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(
		testhelpers.NewMockRouteDependencies(),
		testhelpers.NewMockHTTPClientOK(),
		testhelpers.MockLogger,
	)

	req := httptest.NewRequest(http.MethodGet, "/error/404", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()

	t.Cleanup(func() {
		err := resp.Body.Close()
		if err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var got notfound.ErrorSimulationResponse
	err := json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)

	expected := notfound.ErrorSimulationResponse{Message: "Simulated 404 Not Found"}
	assert.Equal(t, expected, got)
}
