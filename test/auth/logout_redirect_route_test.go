package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRouterLogoutRedirectRouteReturns302Redirect(t *testing.T) {
	t.Parallel()

	provider := testhelpers.NewMockRouteDependencies()
	client := testhelpers.NewMockHTTPClientOK()
	mockLogger := zap.NewNop()

	r := router.NewRouter(provider, client, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Equal(t, "/", location, "unexpected redirect Location header")
}
