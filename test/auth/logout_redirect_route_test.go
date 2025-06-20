package auth_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRouterLogoutRedirectRouteReturns200OK(t *testing.T) {
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

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
	assert.Equal(t, "Logout successful", string(body))
}
