package callback

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/config"

	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct{}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(nil, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
	assert.Equal(nil, "Basic Y2xpZW50MTIzOnNlY3JldDQ1Ng==", req.Header.Get("Authorization"))

	mockResp := callback.TokenResponse{
		AccessToken:  "ACCESS123",
		IDToken:      "ID123",
		RefreshToken: "REFRESH123",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}

	body, _ := json.Marshal(mockResp)

	rec := httptest.NewRecorder()
	rec.Write(body)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       httptest.NewBody(body),
		Header:     make(http.Header),
	}, nil
}

type mockURLProvider struct{}

func (m *mockURLProvider) GetCallbackURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/oauth2/token", nil
}

func TestCallbackHandlerSuccess(t *testing.T) {
	cfg := &config.Config{
		UserPoolClientID: "client123",
		ClientSecret:     "secret456",
		RedirectURI:      "https://example.com/callback",
	}

	deps := callback.CallbackHandlerDependencies{
		Config:      cfg,
		HTTPClient:  &mockHTTPClient{},
		URLProvider: &mockURLProvider{},
	}

	handler := callback.CallbackHandler(deps)

	req := httptest.NewRequest("GET", "/callback?code=testcode&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result callback.TokenResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "ACCESS123", result.AccessToken)
	assert.Equal(t, "ID123", result.IDToken)
	assert.Equal(t, "REFRESH123", result.RefreshToken)
	assert.Equal(t, 3600, result.ExpiresIn)
	assert.Equal(t, "Bearer", result.TokenType)
}
