package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/test_helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoginRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	router := router.NewRouter(test_helpers.MockCfg)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Location"), "oauth2/authorize")
}
