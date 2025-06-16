package openapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestOpenAPIEndpoint_ReturnsJSON(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(
		testhelpers.NewMockRouteDependencies(),
		testhelpers.NewMockHTTPClientOK(),
		testhelpers.MockLogger,
	)

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			t.Errorf("failed to close response body: %v", cerr)
		}
	}()

	body, readErr := io.ReadAll(resp.Body)
	assert.NoError(t, readErr, "failed to read response body")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var spec map[string]any
	readErr = json.Unmarshal(body, &spec)
	assert.NoError(t, readErr, "failed to unmarshal OpenAPI JSON")

	openapiVersion, ok := spec["openapi"].(string)
	assert.True(t, ok, "openapi field is not a string")
	assert.Equal(t, "3.0.0", openapiVersion)
}
