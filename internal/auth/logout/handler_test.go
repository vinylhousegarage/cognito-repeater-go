package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/authtesthelpers"
	"cognito-repeater-go/internal/config"

	"github.com/stretchr/testify/assert"
)

type mockLogoutURLProvider struct{}

func (m *mockLogoutURLProvider) GetLogoutURL(mup config.MetadataURLProvider) (string, error) {
	return "https://example.com/logout", nil
}

func TestLogoutHandlerRedirectsToLogoutEndpoint(t *testing.T) {
	t.Parallel()

	handler := LogoutHandler(&mockLogoutURLProvider{}, &authtesthelpers.MockMetadataURLProvider{})

	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
