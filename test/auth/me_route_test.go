package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func GenerateTestJWKS() (*JWKS, *rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.PublicKey.E)).Bytes())

	jwk := JWK{
		Kid: "test-kid",
		Kty: "RSA",
		Alg: "RS256",
		Use: "sig",
		N:   n,
		E:   e,
	}

	return &JWKS{Keys: []JWK{jwk}}, key, nil
}

func GenerateSignedToken(privateKey *rsa.PrivateKey, issuer, audience string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   "test-sub",
		ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-kid"

	return token.SignedString(privateKey)
}

func TestMeHandler_Integration_Success(t *testing.T) {
	t.Parallel()

	type mockAllProviders struct {
		*config.Config
	}

	func (m *mockAllProviders) GetJWKSURI() string {
		return "https://example.com/jwks"
	}

	cfg := &config.Config{
		Region:           "ap-northeast-1",
		UserPoolID:       "test-pool",
		UserPoolClientID: "test-client",
	}

	provider := &mockAllProviders{Config: cfg}

	jwks, privateKey, err := GenerateTestJWKS()
	assert.NoError(t, err)

	token, err := GenerateSignedToken(privateKey, cfg.Issuer(), cfg.Audience())
	assert.NoError(t, err)

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			switch {
			case strings.Contains(req.URL.String(), "/.well-known/openid-configuration"):
				rec.WriteHeader(http.StatusOK)
				_, _ = rec.WriteString(`{"jwks_uri": "https://example.com/jwks"}`)

			case strings.Contains(req.URL.String(), "/jwks"):
				rec.WriteHeader(http.StatusOK)
				data, err := json.Marshal(jwks)
				if err != nil {
					rec.WriteHeader(http.StatusInternalServerError)
				}
				_, _ = rec.Write(data)

			default:
				rec.WriteHeader(http.StatusNotFound)
			}
			return rec.Result(), nil
		},
	}

	router := router.NewRouter(provider, provider, provider, client)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "test-sub", body["sub"])
}
