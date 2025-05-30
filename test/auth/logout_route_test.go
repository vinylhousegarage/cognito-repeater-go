package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLogoutRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	deps := testhelpers.NewMockRouteDependencies()
	cli := testhelpers.NewMockHTTPClientOK()

	r := router.NewRouter(deps, cli)

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Location"), "logout")
}
