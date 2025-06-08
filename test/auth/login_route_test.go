package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestLoginRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()

			switch {
			case strings.Contains(req.URL.String(), "/.well-known/openid-configuration"):
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.WriteString(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`)

			default:
				rec.WriteHeader(http.StatusNotFound)
			}

			return rec.Result(), nil
		},
	}

	mockLogger := zap.NewNop()

	r := router.NewRouter(testhelpers.NewMockRouteDependencies(), mockClient, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(t, location, "https://example.com/oauth2/authorize")
}
