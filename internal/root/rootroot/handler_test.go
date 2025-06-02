package rootroot

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootHandlerRedirectsToLogin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(NewRootHandler)
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/login", resp.Header.Get("Location"))
}
