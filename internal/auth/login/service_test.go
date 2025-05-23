package login

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

var _ LoginURLProvider = (*loginClient)(nil)

func TestGetLoginURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	client := NewLoginClient(http.DefaultClient)
	endpoint, err := client.GetLoginURL(mock)

	assert.NoError(t, err, "failed to fetch authorization endpoint")
	assert.Equal(t, "https://example.com/oauth2/authorize", endpoint)
}

func TestGetLoginURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	svc := NewLoginClient(http.DefaultClient)
	_, err := svc.GetLoginURL(mock)

	assert.Error(t, err)
}
