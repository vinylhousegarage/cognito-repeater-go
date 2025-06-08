package notfound

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestNewError404Handler_StatusAndBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/error/404", nil)
	w := httptest.NewRecorder()

	mockLogger := zap.NewNop()

	NewError404Handler(mockLogger)(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "not found 404", string(body))
}
