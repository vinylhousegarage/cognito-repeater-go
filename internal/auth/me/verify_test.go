package me

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestParseAndVerifyJWT_Success(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	result, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.NoError(t, err)
	assert.Equal(t, "user-123", result.Subject)
}

func TestParseAndVerifyJWT_InvalidAlg(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	})
	signed, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err)

	_, err = ParseAndVerifyJWT(signed, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.ErrorIs(t, err, ErrInvalidSigningAlg)
}

func TestParseAndVerifyJWT_InvalidToken(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	brokenToken := "this.is.not.a.valid.jwt"

	_, err := ParseAndVerifyJWT(brokenToken, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.ErrorIs(t, err, ErrJWTParseFailed)
}

func TestParseAndVerifyJWT_Expired(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestParseAndVerifyJWT_InvalidIssuer(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "wrong-issuer",
		Subject:   "user-123",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "expected-issuer", "test-aud")
	assert.ErrorIs(t, err, ErrInvalidIssuer)
}

func TestParseAndVerifyJWT_InvalidAudience(t *testing.T) {
	t.Parallel()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "user-123",
		Audience:  []string{"wrong-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "expected-aud")
	assert.ErrorIs(t, err, ErrInvalidAudience)
}

func TestParseAndVerifyJWT_MissingSub(t *testing.T) {
	t.Parallel()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Audience:  []string{"test-aud"},
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	}

	token := generateTestToken(t, claims, privKey)

	_, err = ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
	assert.ErrorIs(t, err, ErrMissingSubject)
}
