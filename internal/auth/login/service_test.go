package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

var _ LoginURLProvider = (*loginClient)(nil)

func TestGetLoginURLReturnsExpectedEndpoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	svc := NewLoginClient(testhelpers.MockClient)
	endpoint, err := svc.GetLoginURL(mock)

	expected := "https://example.com/oauth2/authorize"

	assert.NoError(t, err, "failed to fetch authorization endpoint")
	assert.Equal(t, expected, endpoint)
}

func TestGetLoginURLStatusCode500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	svc := NewLoginClient(http.DefaultClient)
	_, err := svc.GetLoginURL(mock)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
