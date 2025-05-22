package logout

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

var _ LogoutURLProvider = (*logoutClient)(nil)

func TestGetLogoutURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"end_session_endpoint":"https://example.com/logout"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		},
	}

	svc := NewLogoutClient(mockClient)
	endpoint, err := svc.GetLogoutURL(mock)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/logout", endpoint)
}

func TestGetLogoutURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	svc := NewLogoutClient(http.DefaultClient)
	_, err := svc.GetLogoutURL(mock)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
