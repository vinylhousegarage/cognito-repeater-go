package me

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(t *testing.T, claims jwt.RegisteredClaims, privateKey *rsa.PrivateKey) string {
	t.Helper()
	t.Parallel()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	assert.NoError(t, err)
	return signed
}

func TestParseAndVerifyJWT_Success(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	result, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", result.Subject)
}

func TestParseAndVerifyJWT_Expired(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 過去
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "claims validation failed")
}

func TestParseAndVerifyJWT_InvalidIssuer(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "wrong-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "expected-issuer", "test-aud")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected issuer")
}

func TestParseAndVerifyJWT_InvalidAudience(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"wrong-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "expected-aud")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid audience")
}

func TestParseAndVerifyJWT_MissingSub(t *testing.T) {
	t.Parallel()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Subject 欠落
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing subject")
}
