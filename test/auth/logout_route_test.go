package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLogoutRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()

			switch {
			case strings.Contains(req.URL.String(), "/.well-known/openid-configuration"):
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.WriteString(`{"end_session_endpoint": "https://example.com/logout"}`)

			default:
				rec.WriteHeader(http.StatusNotFound)
			}

			return rec.Result(), nil
		},
	}

	r := router.NewRouter(testhelpers.NewMockRouteDependencies(), mockClient, testhelpers.MockLogger)
	req := httptest.NewRequest(http.MethodGet, "/logout?id_token_hint=dummy-id-token", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Location"), "https://example.com/logout")
}
