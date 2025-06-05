package ping

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPingHandler_ReturnsPong(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	NewPingHandler(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res PingResponse
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)
	assert.Equal(t, "pong", res.Message)
}
