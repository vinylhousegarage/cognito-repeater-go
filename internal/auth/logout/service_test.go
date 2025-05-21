package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth"

	"github.com/stretchr/testify/assert"
)

var _ LogoutEndpointProvider = (*logoutService)(nil)

func TestGetLogoutURLReturnsExpectedEndpoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
	}))
	defer ts.Close()

	mock := &auth.MockMetadataURL{URL: ts.URL}

	svc := NewLogoutService(http.DefaultClient)
	endpoint, err := svc.GetLogoutURL(mock)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/logout", endpoint)
}

func TestGetLogoutURLStatusCode500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mock := &auth.MockMetadataURL{URL: ts.URL}

	svc := NewLogoutService(http.DefaultClient)
	_, err := svc.GetLogoutURL(mock)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
