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

	r := router.NewRouter(provider, client, testhelpers.MockLogger)

	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			testhelpers.MockLogger.Warn("failed to close response body", zap.Error(cerr))
		}
	}()

	assert.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Equal(t, "/", location, "unexpected redirect Location header")
}
