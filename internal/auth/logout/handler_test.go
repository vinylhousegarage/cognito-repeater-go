package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLogoutHandlerRedirectsToLogoutEndpoint(t *testing.T) {
	t.Parallel()

	d := deps.HandlerDependencies{
		Config: &testhelpers.MockMetadataURL{URL: "https://mock.metadata.url"},
		HTTPClient: &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				rec := httptest.NewRecorder()
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
				return rec.Result(), nil
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	handler := LogoutHandler(d)
	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
