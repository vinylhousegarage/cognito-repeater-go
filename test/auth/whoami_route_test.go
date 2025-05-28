package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/whoami"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	assert.Equal(t, "abc123", body["sub"])
	assert.Equal(t, "user@example.com", body["email"])
}
