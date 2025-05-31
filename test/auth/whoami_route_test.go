package auth_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/whoami"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestWhoamiHandler_ReturnsUserinfo(t *testing.T) {
	t.Parallel()

	var ts *httptest.Server

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			w.WriteHeader(http.StatusOK)
			if _, err := fmt.Fprintf(w, `{"userinfo_endpoint": "%s/oauth2/userinfo"}`, ts.URL); err != nil {
				t.Errorf("failed to write openid-configuration: %v", err)
			}
		case "/oauth2/userinfo":
			if r.Header.Get("Authorization") != "Bearer valid-token" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err := io.WriteString(w, `{"sub":"abc123"}`); err != nil {
				t.Errorf("failed to write userinfo response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	})

	ts = httptest.NewServer(handler)
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	provider := &testhelpers.MockMetadataURL{URL: ts.URL + "/.well-known/openid-configuration"}
	handlerFunc := whoami.WhoamiHandler(provider, testhelpers.NewMockHTTPClientOK())

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handlerFunc(rec, req)

	resp := rec.Result()
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)

	assert.Equal(t, "abc123", body["sub"])
}
