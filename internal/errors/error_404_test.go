package errors

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestError404Handler_StatusAndBody(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/error/404", nil)
    w := httptest.NewRecorder()

    Error404Handler(w, req)

    resp := w.Result()
    body, err := io.ReadAll(resp.Body)
    assert.NoError(t, err)

    assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    assert.Equal(t, "not found 404", string(body))
}
