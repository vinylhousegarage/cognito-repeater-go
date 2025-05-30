package auth_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLoginRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	deps := testhelpers.MockLoginHandlerProvider{}
	cli := testhelpers.NewMockHTTPClientOK()

	r := router.NewRouter(deps, cli)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(t, location, "https://example.com/oauth2/authorize")
}
