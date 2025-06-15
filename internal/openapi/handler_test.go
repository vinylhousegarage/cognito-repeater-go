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

	merr := os.MkdirAll("docs", 0755)
	if merr != nil {
		t.Errorf("failed to make docs directory: %v", merr)
	}

	require.NoError(t, os.WriteFile("docs/swagger.json", []byte(`{"openapi":"3.0.0"}`), 0644))

	t.Cleanup(func() {
		removeErr := os.RemoveAll("docs")
		if removeErr != nil {
			t.Errorf("failed to remove docs directory: %v", removeErr)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			t.Errorf("failed to close response body: %v", cerr)
		}
	}()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		t.Errorf("failed to read response body: %v", readErr)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"openapi":"3.0.0"}`, string(body))
}
