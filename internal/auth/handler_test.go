package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogoutHandler_StatusAndLocation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	LogoutHandler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://stab.com/logout", resp.Header.Get("Location"))
}

func TestMetadataHandler_ReturnsExpectedStatusAndJSONBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	LogoutRedirectHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	expected := map[string]string{
		"message": "Logout successful",
	}

	var actual map[string]string
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
