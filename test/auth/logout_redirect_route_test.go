package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"

	"github.com/stretchr/testify/assert"
)

type mockMetadataURL struct {
	URL string
}

func (m *mockMetadataURL) MetadataURL() string {
	return m.URL
}

func newMockProvider(t *testing.T) (*mockMetadataURL, func()) {
	mockJSON := `{
		"end_session_endpoint": "https://example.com/logout"
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	t.Cleanup(ts.Close)

	return &mockMetadataURL{URL: ts.URL}, ts.Close
}

func TestRouter_LogoutRedirectRoute_ReturnsExpectedJSON(t *testing.T) {
	mockProvider, teardown := newMockProvider(t)
	defer teardown()

	r := router.NewRouter(mockProvider)
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
