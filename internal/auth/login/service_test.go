package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestGetLoginURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
	}))
	defer ts.Close()

	mockLogger := zap.NewNop()

	endpoint, err := GetLoginURL(ts.URL, http.DefaultClient, mockLogger)

	assert.NoError(t, err, "failed to fetch authorization endpoint")
	assert.Equal(t, "https://example.com/oauth2/authorize", endpoint)
}

func TestGetLoginURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mockLogger := zap.NewNop()

	_, err := GetLoginURL(ts.URL, http.DefaultClient, mockLogger)

	assert.Error(t, err, "unexpected status code")
}
