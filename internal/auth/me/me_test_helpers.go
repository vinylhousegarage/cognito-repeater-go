package me

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

type MeTestEnv struct {
	PrivKey         *rsa.PrivateKey
	PubKey          *rsa.PublicKey
	MockJwksURL     string
	MockMetadataURL string
	MockIssuer      string
	MockAudience    string
	Logger          *zap.Logger
	Provider        deps.MeHandlerProvider
}

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

func ServeMeHandler(req *http.Request, provider deps.MeHandlerProvider, client httpclient.HTTPClient, logger *zap.Logger) *httptest.ResponseRecorder {
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
