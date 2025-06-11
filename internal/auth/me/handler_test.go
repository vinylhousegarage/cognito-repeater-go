package me

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

type mockMeHandlerProvider struct {
	mockAudience    string
	mockGetJWKSURI  string
	mockIssuer      string
	mockMetadataURL string
}

func (m *mockMeHandlerProvider) Audience() string    { return m.mockAudience }
func (m *mockMeHandlerProvider) GetJWKSURI() string  { return m.mockGetJWKSURI }
func (m *mockMeHandlerProvider) Issuer() string      { return m.mockIssuer }
func (m *mockMeHandlerProvider) MetadataURL() string { return m.mockMetadataURL }

type mockHTTPClient struct {
	Responses map[string]*http.Response
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if resp, ok := m.Responses[req.URL.String()]; ok {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func generateNewMeHandlerTestToken(t *testing.T, claims jwt.RegisteredClaims, privKey *rsa.PrivateKey, kid string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(privKey)
	assert.NoError(t, err)
	return signed
}

func TestNewMeHandler_Success(t *testing.T) {
	t.Parallel()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	pubKey := &privKey.PublicKey

	pubKeyN := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())
	pubKeyE := base64.RawURLEncoding.EncodeToString([]byte{byte(pubKey.E)})
	jwkJSON := `{"alg":"RS256","e":"` + pubKeyE + `","kid":"test-kid","kty":"RSA","n":"` + pubKeyN + `"}`
	jwkSetJSON := `{"keys": [` + jwkJSON + `]}`

	mockJwksURL := "http://mock-cognito/jwks"
	mockMetadataURL := "http://mock-cognito/.well-known/openid-configuration"
	mockIssuer := "https://cognito-idp.us-west-2.amazonaws.com/mockpool"
	mockAudience := "mockclientid"

	metadataResponse := `{"jwks_uri": "` + mockJwksURL + `", "issuer": "` + mockIssuer + `"}`

	mockHTTPClient := &mockHTTPClient{
		Responses: map[string]*http.Response{
			mockMetadataURL: {
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(metadataResponse)),
			},
			mockJwksURL: {
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(jwkSetJSON)),
			},
		},
	}

	mockProvider := &mockMeHandlerProvider{
		mockMetadataURL: mockMetadataURL,
		mockIssuer:      mockIssuer,
		mockAudience:    mockAudience,
	}

	testClaims := jwt.RegisteredClaims{
		Issuer:    mockIssuer,
		Audience:  []string{mockAudience},
		Subject:   "test-user-123",
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	}

	validToken := generateNewMeHandlerTestToken(t, testClaims, privKey, "test-kid")

	formData := url.Values{}
	formData.Set("token", validToken)
	reqBody := strings.NewReader(formData.Encode())
	req := httptest.NewRequest(http.MethodPost, "/me", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := NewMeHandler(mockProvider, mockHTTPClient, zap.NewExample())
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Logf("Response body (on error): %s", rr.Body.String())
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var userResp UserResponse
	err = json.NewDecoder(rr.Body).Decode(&userResp)
	assert.NoError(t, err)
	assert.Equal(t, "test-user-123", userResp.Sub)
}
