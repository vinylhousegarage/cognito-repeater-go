package errors_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/server"

	"github.com/stretchr/testify/assert"
)

func TestServer_Error404_E2E(t *testing.T) {
	s := server.NewServer()

	ts := httptest.NewServer(s.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/error/404")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "not found 404", string(body))
}
