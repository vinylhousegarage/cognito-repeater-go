package root_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "io"

    "cognito-repeater-go/internal/router"
    "github.com/stretchr/testify/assert"
)

func TestRouter_PingRouteRegistered(t *testing.T) {
    r := router.NewRouter()

    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    w := httptest.NewRecorder()

    r.ServeHTTP(w, req)

    resp := w.Result()
    body, _ := io.ReadAll(resp.Body)

    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Equal(t, "pong", string(body))
}
