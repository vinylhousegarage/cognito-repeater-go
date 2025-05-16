package root

import (
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestPingHandler_ReturnsPong(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    w := httptest.NewRecorder()

    PingHandler(w, req)

    resp := w.Result()
    body, _ := io.ReadAll(resp.Body)

    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Equal(t, "pong", string(body))
}
