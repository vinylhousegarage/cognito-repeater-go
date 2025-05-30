package simulatederrors_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestError404RouteReturnsPlainTextNotFound(t *testing.T) {
	t.Parallel()

	deps := auth.RouteDependencies{
		CallbackProvider: testhelpers.MockCallbackProvider{},
		LoginProvider: testhelpers.MockLoginProvider{},
		LogoutProvider: testhelpers.MockLogoutProvider{},
		MeProvider: testhelpers.MockMeProvider{},
	}

	cli := testhelpers.NewMockHTTPClientOK()

	r := router.NewRouter(deps, cli)

	req := httptest.NewRequest(http.MethodGet, "/error/404", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "not found 404", string(body))
}
