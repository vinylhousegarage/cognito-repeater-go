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

	mockMetadataURL := "https://mock.metadata.url"

	d := deps.HandlerDependencies{
		Config: testhelpers.MockCfg,
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

	handler := func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := GetLogoutURL(d.HTTPClient, mockMetadataURL)
		if err != nil {
			http.Error(w, "failed to fetch logout url", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, endpoint, http.StatusFound)
	}

	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
