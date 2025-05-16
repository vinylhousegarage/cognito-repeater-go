package e2e_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute_E2E(t *testing.T) {
	ts := httptest.NewServer(router.NewRouter())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}
