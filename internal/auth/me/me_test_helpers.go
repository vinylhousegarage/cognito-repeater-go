package me

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
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
	"github.com/stretchr/testify/require"

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
	t.Helper()

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

func generateNewMeHandlerTestToken(t *testing.T, claims jwt.RegisteredClaims, privKey *rsa.PrivateKey, kid string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(privKey)
	assert.NoError(t, err)

	return signed
}

func NewAuthHeaderRequest(token string) *http.Request {
    req := httptest.NewRequest(http.MethodGet, "/me", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    return req
}

func ServeMeHandler(req *http.Request, provider deps.MeHandlerProvider, client httpclient.HTTPClient, logger *zap.Logger) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := NewMeHandler(provider, client, logger)
	handler.ServeHTTP(rr, req)
	return rr
}

// 400 Bad Request Response
func AssertBadRequestResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var body map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, expectedMessage, body["error"])
}

// 401 Unauthorized Response
func AssertUnauthorizedResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var body map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, expectedMessage, body["error"])
}

// 500 Internal Server Error Response
func AssertInternalServerErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var body map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, expectedMessage, body["error"])
}

// 502 Bad Gateway Response
func AssertBadGatewayResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()
	assert.Equal(t, http.StatusBadGateway, rr.Code)

	var body map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, expectedMessage, body["error"])
}

// test_helpers build function
func setupWithClaims(t *testing.T, claims jwt.RegisteredClaims) *httptest.ResponseRecorder {
	t.Helper()

	env := NewTestEnv(t)
	kid := "test-kid"
	jwkJSON := generateJWKJSON(env.PubKey, kid)
	client := buildMockHTTPClient(env.MockMetadataURL, env.MockJwksURL, env.MockIssuer, jwkJSON)
	token := generateNewMeHandlerTestToken(t, claims, env.PrivKey, kid)
	req := NewAuthHeaderRequest(token)
	return ServeMeHandler(req, env.Provider, client, env.Logger)
}
