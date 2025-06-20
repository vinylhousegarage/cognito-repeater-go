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

	t.Run("GET returns 302 with /", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
		w := httptest.NewRecorder()

		NewLogoutRedirectHandler(testhelpers.MockLogger)(w, req)

		resp := w.Result()
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Fatalf("failed to close response body: %v", err)
			}
		}()

		assert.Equal(t, http.StatusFound, resp.StatusCode)
		assert.Equal(t, "/", resp.Header.Get("Location"))
	})

	t.Run("POST returns 405", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/logout/redirect", nil)
		w := httptest.NewRecorder()

		NewLogoutRedirectHandler(testhelpers.MockLogger)(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
