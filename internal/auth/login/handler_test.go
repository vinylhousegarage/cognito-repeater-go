package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/test/testhelpers"
	"cognito-repeater-go/internal/config"

	"github.com/stretchr/testify/assert"
)

type mockLoginURLProvider struct{}

func (m *mockLoginURLProvider) GetLoginURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/oauth2/authorize", nil
}

func TestLoginHandlerRedirectsToLoginEndpoint(t *testing.T) {
	t.Parallel()

	handler := LoginHandler(&mockLoginURLProvider{}, &testhelpers.MockMetadataURLProvider{})

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/oauth2/authorize", resp.Header.Get("Location"))
}
