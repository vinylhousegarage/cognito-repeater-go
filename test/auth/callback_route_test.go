package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestCallbackRouteReturnsTokenJSON(t *testing.T) {
	t.Parallel()

	mockResp := map[string]interface{}{
		"access_token":  "ACCESS123",
		"id_token":      "ID123",
		"refresh_token": "REFRESH123",
		"expires_in":    3600,
		"token_type":    "Bearer",
	}
	body, _ := json.Marshal(mockResp)

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()

			if req.Method == http.MethodGet {
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.Write([]byte(`{"token_endpoint":"https://example.com/token"}`))
			} else {
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.Write(body)
			}

			return rec.Result(), nil
		},
	}

	cfg := &config.Config{
		Region:           "ap-northeast-1",
		UserPoolID:       "test-pool",
		UserPoolClientID: "client123",
		ClientSecret:     "secret456",
		RedirectURI:      "https://localhost/callback",
	}

	r := router.NewRouter(cfg, client)

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc123&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&tokenResp)
	assert.NoError(t, err)

	assert.Equal(t, "ACCESS123", tokenResp["access_token"])
}
