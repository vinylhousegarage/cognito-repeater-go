package logoutredirect

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGetLogoutURL_StatusCode500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mock := &mockMetadataURL{URL: ts.URL}

	svc := NewLogoutService(http.DefaultClient)
	_, err := svc.GetLogoutURL(mock)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
