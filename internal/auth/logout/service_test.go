package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ LogoutEndpointProvider = (*logoutService)(nil)

type mockMetadataURL struct {
	URL string
}

func (m *mockMetadataURL) MetadataURL() string {
	return m.URL
}

func TestGetLogoutURL_ReturnsExpectedEndpoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
	}))
	defer ts.Close()

	mock := &mockMetadataURL{URL: ts.URL}

	svc := NewLogoutService(http.DefaultClient)
	endpoint, err := svc.GetLogoutURL(mock)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/logout", endpoint)
}
