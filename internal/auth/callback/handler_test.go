package callback

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/config"
	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"

	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestCallbackHandlerSuccess(t *testing.T) {
	t.Parallel()

	mockResp := TokenResponse{
		AccessToken:  "ACCESS123",
		IDToken:      "ID123",
		RefreshToken: "REFRESH123",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}

	body, _ := json.Marshal(mockResp)

	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, _ = rec.Write(body)
			return rec.Result(), nil
		},
	}

	cfg := &config.Config{
		UserPoolClientID: "client123",
		ClientSecret:     "secret456",
		RedirectURI:      "https://example.com/callback",
		MetadataEndpoint: "https://example.com/.well-known/openid-configuration",
	}

	handlerDeps := deps.HandlerDependencies{
		Config:     cfg,
		HTTPClient: httpclient.HTTPClient(mockClient),
	}

	handler := CallbackHandler(handlerDeps)

	req := httptest.NewRequest(http.MethodGet, "/callback?code=testcode&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result TokenResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "ACCESS123", result.AccessToken)
	assert.Equal(t, "ID123", result.IDToken)
	assert.Equal(t, "REFRESH123", result.RefreshToken)
	assert.Equal(t, 3600, result.ExpiresIn)
	assert.Equal(t, "Bearer", result.TokenType)
}
