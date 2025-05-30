package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type mockCallbackHandlerProvider struct{}

func (m *mockCallbackHandlerProvider) Audience() string {
	return m.UserPoolClientIDValue()
}

func (m *mockCallbackHandlerProvider) ClientSecretValue() string {
	return "mock-client-secret"
}

func (m *mockCallbackHandlerProvider) GetJWKSURI() string {
	return "https://example.com/jwks"
}

func (m *mockCallbackHandlerProvider) Issuer() string {
	region := "ap-northeast-1"
	userPoolID := "mock-user-pool-id"
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolID)
}

func (m *mockCallbackHandlerProvider) MetadataURL() string {
	return "https://mock.auth.ap-northeast-1.amazoncognito.com/.well-known/openid-configuration"
}

func (m *mockCallbackHandlerProvider) RedirectURIValue() string {
	return "https://localhost/callback"
}

func (m *mockCallbackHandlerProvider) ScopeValue() string {
	return "openid"
}

func (m *mockCallbackHandlerProvider) UserPoolClientIDValue() string {
	return "mock-client-id"
}

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

	provider := &mockCallbackHandlerProvider{}

	r := router.NewRouter(provider, provider, provider, provider, provider, client)

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc123&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&tokenResp)
	assert.NoError(t, err)
	assert.Equal(t, "ACCESS123", tokenResp["access_token"])
}
