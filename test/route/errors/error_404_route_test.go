package errors_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "io"

    "cognito-repeater-go/internal/router"
    "github.com/stretchr/testify/assert"
)

func TestRouter_Error404RouteRegistered(t *testing.T) {
    r := router.NewRouter()

    req := httptest.NewRequest(http.MethodGet, "/error/404", nil)
    w := httptest.NewRecorder()

    r.ServeHTTP(w, req)

    resp := w.Result()
    body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

    assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    assert.Equal(t, "not found 404", string(body))
}
