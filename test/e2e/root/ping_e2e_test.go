package root_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/server"

	"github.com/stretchr/testify/assert"
)

func TestServer_Ping_E2E(t *testing.T) {
	s := server.NewServer()

	ts := httptest.NewServer(s.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
