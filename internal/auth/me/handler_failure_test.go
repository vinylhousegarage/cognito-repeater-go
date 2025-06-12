package me

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewMeHandler_Errors(t *testing.T) {
	t.Parallel()

	t.Run("missing-subject", func(t *testing.T) {
		t.Parallel()

		rr := setupWithClaims(t, jwt.RegisteredClaims{
			Issuer:    "https://cognito-idp.us-west-2.amazonaws.com/mockpool",
			Audience:  []string{"mockclientid"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("expired-token", func(t *testing.T) {
		t.Parallel()

		rr := setupWithClaims(t, jwt.RegisteredClaims{
			Issuer:    "https://cognito-idp.us-west-2.amazonaws.com/mockpool",
			Audience:  []string{"mockclientid"},
			Subject:   "expired-user",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		})

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("invalid-issuer", func(t *testing.T) {
		t.Parallel()

		rr := setupWithClaims(t, jwt.RegisteredClaims{
			Issuer:    "https://malicious.example.com",
			Audience:  []string{"mockclientid"},
			Subject:   "hacker",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("invalid-audience", func(t *testing.T) {
		t.Parallel()

		rr := setupWithClaims(t, jwt.RegisteredClaims{
			Issuer:    "https://cognito-idp.us-west-2.amazonaws.com/mockpool",
			Audience:  []string{"invalid-client-id"},
			Subject:   "invalid-audience-user",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("broken-jwt-format", func(t *testing.T) {
		t.Parallel()

		env := NewTestEnv(t)
		jwkJSON := generateJWKJSON(env.PubKey, "test-kid")
		mockHTTPClient := buildMockHTTPClient(env.MockMetadataURL, env.MockJwksURL, env.MockIssuer, jwkJSON)
		brokenToken := "this.is.not.a.valid.jwt"
		req := NewTokenPostRequest(brokenToken)
		rr := ServeMeHandler(req, env.Provider, mockHTTPClient, env.Logger)

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("unsupported-algorithm", func(t *testing.T) {
		t.Parallel()

		env := NewTestEnv(t)
		jwkJSON := generateJWKJSON(env.PubKey, "test-kid")
		mockHTTPClient := buildMockHTTPClient(env.MockMetadataURL, env.MockJwksURL, env.MockIssuer, jwkJSON)

		hmacSecret := []byte("supersecret")

		claims := jwt.RegisteredClaims{
			Issuer:    env.MockIssuer,
			Audience:  []string{env.MockAudience},
			Subject:   "hmac-user",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}

		tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(hmacSecret)
		assert.NoError(t, err)

		req := NewTokenPostRequest(tokenStr)
		rr := ServeMeHandler(req, env.Provider, mockHTTPClient, env.Logger)

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("missing-kid-in-header", func(t *testing.T) {
		t.Parallel()

		env := NewTestEnv(t)
		jwkJSON := generateJWKJSON(env.PubKey, "test-kid")
		mockHTTPClient := buildMockHTTPClient(env.MockMetadataURL, env.MockJwksURL, env.MockIssuer, jwkJSON)

		claims := jwt.RegisteredClaims{
			Issuer:    env.MockIssuer,
			Audience:  []string{env.MockAudience},
			Subject:   "user-without-kid",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		delete(token.Header, "kid")

		tokenStr, err := token.SignedString(env.PrivKey)
		assert.NoError(t, err)

		req := NewTokenPostRequest(tokenStr)
		rr := ServeMeHandler(req, env.Provider, mockHTTPClient, env.Logger)

		AssertUnauthorizedResponse(t, rr)
	})

	t.Run("kid-not-found-in-jwk-set", func(t *testing.T) {
		t.Parallel()

		env := NewTestEnv(t)

		knownKid := "known-kid"
		jwkJSON := generateJWKJSON(env.PubKey, knownKid)
		mockHTTPClient := buildMockHTTPClient(env.MockMetadataURL, env.MockJwksURL, env.MockIssuer, jwkJSON)

		unknownKid := "unknown-kid"
		claims := jwt.RegisteredClaims{
			Issuer:    env.MockIssuer,
			Audience:  []string{env.MockAudience},
			Subject:   "user-with-unknown-kid",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = unknownKid

		tokenStr, err := token.SignedString(env.PrivKey)
		assert.NoError(t, err)

		req := NewTokenPostRequest(tokenStr)
		rr := ServeMeHandler(req, env.Provider, mockHTTPClient, env.Logger)

		AssertUnauthorizedResponse(t, rr)
	})
}
