package auth_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/whoami"
	"cognito-repeater-go/test/testhelpers"
)

func TestWhoamiRoute_Success(t *testing.T) {
	t.Parallel()

	userinfoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"sub":"abc123","email":"user@example.com"}`))
	}))
	defer userinfoServer.Close()

	metadataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"userinfo_endpoint":"` + userinfoServer.URL + `"}`))
	}))
	defer metadataServer.Close()

	deps := deps.HandlerDependencies{
		Config:     testhelpers.MockCfg,
		HTTPClient: &testhelpers.MockHTTPClient{DoFunc: http.DefaultClient.Do},
	}

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer mocktoken")
	w := httptest.NewRecorder()

	handler := whoami.WhoamiHandler(deps)
	handler(w, req)

	resp := w.Result()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("DEBUG: whoami response body = %s\n", string(bodyBytes))
}
