package whoami

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserinfoURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"userinfo_endpoint": "https://example.com/oauth2/userinfo"}`))
	}))
	defer ts.Close()

	endpoint, err := GetUserinfoURL(http.DefaultClient, ts.URL)

	assert.NoError(t, err, "failed to fetch userinfo_endpoint")
	assert.Equal(t, "https://example.com/oauth2/userinfo", endpoint)
}

func TestGetUserinfoURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := GetUserinfoURL(http.DefaultClient, ts.URL)

	assert.Error(t, err, "unexpected status code")
}
