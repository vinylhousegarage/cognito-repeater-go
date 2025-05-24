package callback

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestCallbackHandlerSuccess(t *testing.T) {
	t.Parallel()

	mockTokenResp := TokenResponse{
		AccessToken:  "ACCESS123",
		IDToken:      "ID123",
		RefreshToken: "REFRESH123",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}
	tokenBody, _ := json.Marshal(mockTokenResp)

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()

			switch req.Method {
			case http.MethodGet:
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.Write([]byte(`{"token_endpoint": "https://example.com/oauth2/token"}`))
			case http.MethodPost:
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.Write(tokenBody)
			default:
				rec.WriteHeader(http.StatusMethodNotAllowed)
			}

			return rec.Result(), nil
		},
	}

	cfg := &config.Config{
		UserPoolClientID: "client123",
		ClientSecret:     "secret456",
		RedirectURI:      "https://example.com/callback",
		Region:           "ap-northeast-1",
		UserPoolID:       "pool-id",
	}

	handlerDeps := deps.HandlerDependencies{
		Config:     cfg,
		HTTPClient: mockClient,
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
