package me

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io"
	"math/big"
	"net/http"
	"net/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func NewTestEnv(t *testing.T) *MeTestEnv {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	mockJwksURL := "http://mock-cognito/jwks"
	mockMetadataURL := "http://mock-cognito/.well-known/openid-configuration"
	mockIssuer := "https://cognito-idp.us-west-2.amazonaws.com/mockpool"
	mockAudience := "mockclientid"
	logger := zap.NewExample()

	return &MeTestEnv{
		PrivKey:         privKey,
		PubKey:          &privKey.PublicKey,
		MockJwksURL:     mockJwksURL,
		MockMetadataURL: mockMetadataURL,
		MockIssuer:      mockIssuer,
		MockAudience:    mockAudience,
		Logger:          logger,
		Provider: &mockMeHandlerProvider{
			mockMetadataURL: mockMetadataURL,
			mockIssuer:      mockIssuer,
			mockAudience:    mockAudience,
		},
	}
}

func generateJWKJSON(pubKey *rsa.PublicKey, kid string) string {
	eBytes := big.NewInt(int64(pubKey.E)).Bytes()
	pubKeyE := base64.RawURLEncoding.EncodeToString(eBytes)
	pubKeyN := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())
	return fmt.Sprintf(`{"alg":"RS256","e":"%s","kid":"%s","kty":"RSA","n":"%s"}`, pubKeyE, kid, pubKeyN)
}

func buildMockHTTPClient(metadataURL, jwksURL, issuer, jwkJSON string) *mockHTTPClient {
	jwkSetJSON := `{"keys": [` + jwkJSON + `]}`
	metadataResponse := `{"jwks_uri": "` + jwksURL + `", "issuer": "` + issuer + `"}`

	return &mockHTTPClient{
		Responses: map[string]*http.Response{
			metadataURL: {
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(metadataResponse)),
			},
			jwksURL: {
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(jwkSetJSON)),
			},
		},
	}
}

func NewTokenPostRequest(token string) *http.Request {
	form := url.Values{}
	form.Set("token", token)
	req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func ServeMeHandler(req *http.Request, provider Provider, client HTTPClient, logger *zap.Logger) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := NewMeHandler(provider, client, logger)
	handler.ServeHTTP(rr, req)
	return rr
}

func AssertUnauthorizedResponse(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Unauthorized")
}

func setupWithClaims(t *testing.T, claims jwt.RegisteredClaims) *httptest.ResponseRecorder {
	t.Helper()

	env := NewTestEnv(t)
	kid := "test-kid"
	jwkJSON := generateJWKJSON(env.PubKey, kid)
	client := buildMockHTTPClient(env.MockMetadataURL, env.MockJwksURL, env.MockIssuer, jwkJSON)
	token := generateNewMeHandlerTestToken(t, claims, env.PrivKey, kid)
	req := NewTokenPostRequest(token)
	return ServeMeHandler(req, env.Provider, client, env.Logger)
}
