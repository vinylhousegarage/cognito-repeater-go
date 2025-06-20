package logoutredirect

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLogoutRedirectHandler(t *testing.T) {
	t.Parallel()

	t.Run("GET returns 200 with logout message", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
		w := httptest.NewRecorder()

		NewLogoutRedirectHandler(testhelpers.MockLogger)(w, req)

		resp := w.Result()
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				_ = cerr
			}
		}()

		body := w.Body.String()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
		assert.Equal(t, "Logout successful", body)
	})

	t.Run("POST returns 405", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/logout/redirect", nil)
		w := httptest.NewRecorder()

		NewLogoutRedirectHandler(testhelpers.MockLogger)(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
