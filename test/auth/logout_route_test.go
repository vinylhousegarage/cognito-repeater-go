import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/logout"

	"github.com/stretchr/testify/assert"
)

func TestLogoutRouteIsRegistered(t *testing.T) {
	router := http.NewServeMux()
	router.Handle("/logout", logout.LogoutHandler(&mockEndpointProvider{})(&mockMetadataProvider{}))

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
}
