package openapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestNewOpenAPIHandler(t *testing.T) {
	t.Parallel()

	spec := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Mock API",
			Version: "1.0.0",
		},
		Paths: &openapi3.Paths{},
	}

	logger := zap.NewExample()

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()

	handler := NewOpenAPIHandler(spec, testhelpers.MockLogger)
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close response body", zap.Error(cerr))
		}
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	expected, err := spec.MarshalJSON()
	require.NoError(t, err)

	actual := w.Body.Bytes()
	assert.JSONEq(t, string(expected), string(actual))
}
