package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestGetLogoutURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
	}))
	defer ts.Close()

	provider := &config.Config{MetadataEndpoint: ts.URL}

	endpoint, err := GetLogoutURL(http.DefaultClient, mock)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/logout", endpoint)
}

func TestGetLogoutURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	provider := &config.Config{MetadataEndpoint: ts.URL}

	endpoint, err := GetLogoutURL(http.DefaultClient, provider)

	assert.Error(t, err, "unexpected status code")
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
