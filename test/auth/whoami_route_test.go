package auth_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestWhoamiRoute_Integration(t *testing.T) {
	t.Parallel()

	var ts *httptest.Server

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, `{"userinfo_endpoint": "%s/oauth2/userinfo"}`, ts.URL)
			if err != nil {
				t.Errorf("failed to write openid config: %v", err)
			}
		case "/oauth2/userinfo":
			if r.Header.Get("Authorization") != "Bearer valid-token" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, `{"sub":"abc123"}`)
			if err != nil {
				t.Errorf("failed to write userinfo response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	})

	ts = httptest.NewServer(handler)
	defer ts.Close()

	mockDeps := testhelpers.NewMockRouteDependencies()
	mockDeps.WhoamiProvider = &testhelpers.MockMetadataURL{URL: ts.URL + "/.well-known/openid-configuration"}

	mockLogger := zap.NewNop()

	r := router.NewRouter(mockDeps, http.DefaultClient, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

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
