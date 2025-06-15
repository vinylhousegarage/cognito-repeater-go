package openapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenAPIHandler_FileExists(t *testing.T) {
	t.Parallel()

	handler := NewOpenAPIHandler(testhelpers.MockLogger)

	_ = os.MkdirAll("docs", 0755)
	require.NoError(t, os.WriteFile("docs/swagger.json", []byte(`{"openapi":"3.0.0"}`), 0644))

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"openapi":"3.0.0"}`, string(body))
}
