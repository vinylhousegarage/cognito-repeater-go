package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"

	"github.com/stretchr/testify/assert"
)

type mockMetadataURL struct {
	URL string
}

func (m *mockMetadataURL) MetadataURL() string {
	return m.URL
}

func TestRouter_LogoutRedirectRoute_ReturnsExpectedJSON(t *testing.T) {
	r := router.NewRouter(p config.MetadataURLProvider)

	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := map[string]string{
		"message": "Logout successful",
	}

	var actual map[string]string
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestRouter_LogoutRoute_WithInjectedProvider_Returns302Redirect(t *testing.T) {
	mockJSON := `{
		"end_session_endpoint": "https://example.com/logout"
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	defer ts.Close()

	mockProvider := &mockMetadataURL{URL: ts.URL}

	r := mux.NewRouter(p config.MetadataURLProvider)
	r.HandleFunc("/logout", auth.LogoutHandler(mockProvider))

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com/logout", resp.Header.Get("Location"))
}
