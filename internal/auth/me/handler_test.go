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

	"cognito-repeater-go/test/testhelpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(t *testing.T, claims jwt.RegisteredClaims, privateKey *rsa.PrivateKey) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	assert.NoError(t, err)
	return signed
}

type MockMeHandlerProvider struct {
	mockMetadataURL string
	mockIssuer      string
	mockAudience    string
}

func (m *MockMeHandlerProvider) MetadataURL() string {
	return m.mockMetadataURL
}

func (m *MockMeHandlerProvider) Issuer() string {
	return m.mockIssuer
}

func (m *MockMeHandlerProvider) Audience() string {
	return m.mockAudience
}

func TestNewMeHandler_Success(t *testing.T) {
	t.Parallel()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	pubKey := &privKey.PublicKey

	testTime := time.Date(2025, time.June, 10, 10, 0, 0, 0, time.UTC)

	pubKeyN := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())
	pubKeyE := base64.RawURLEncoding.EncodeToString([]byte{byte(pubKey.E)})
	jwkJSON := `{"alg":"RS256","e":"` + pubKeyE + `","kid":"test-kid","kty":"RSA","n":"` + pubKeyN + `"}`
	jwkSetJSON := `{"keys": [` + jwkJSON + `]}`

	mockMetadataURL := "http://mock-cognito/.well-known/openid-configuration"
	mockIssuer := "https://cognito-idp." + testhelpers.MockCfg.Region + ".amazonaws.com/" + testhelpers.MockCfg.UserPoolID
	mockAudience := testhelpers.MockCfg.UserPoolClientID

	metadataResponse := `{"jwks_uri": "` + mockJwksURL + `", "issuer": "` + mockIssuer + `"}`

	mockHTTPClient := &testhelpers.MockHTTPClient{
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

	mockProvider := &MockMeHandlerProvider{
		mockMetadataURL: mockMetadataURL,
		mockIssuer:      mockIssuer,
		mockAudience:    mockAudience,
	}

	testClaims := jwt.RegisteredClaims{
		Issuer:    mockIssuer,
		Audience:  []string{mockAudience},
		Subject:   "test-user-123",
		ExpiresAt: jwt.NewNumericDate(testTime.Add(1 * time.Hour)),
	}

	tokenWithKid := jwt.NewWithClaims(jwt.SigningMethodRS256, testClaims)
	tokenWithKid.Header["kid"] = "test-kid"
	validToken := generateTestToken(t, testClaims, privKey)

	formData := url.Values{}
	formData.Set("token", validToken)
	reqBody := strings.NewReader(formData.Encode())

	req := httptest.NewRequest(http.MethodPost, "/me", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := NewMeHandler(mockProvider, mockHTTPClient, testhelpers.MockLogger)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "expected status 200 OK")

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "expected Content-Type to be application/json")

	var userResp UserResponse
	err = json.NewDecoder(rr.Body).Decode(&userResp)
	assert.NoError(t, err, "failed to decode response body")
	assert.Equal(t, "test-user-123", userResp.Sub, "expected sub claim to match")
}
