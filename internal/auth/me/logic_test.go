package me

import (
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKIDFromToken(t *testing.T) {
	t.Parallel()

	t.Run("valid token with kid", func(t *testing.T) {
		header := `{"kid":"example-kid"}`
		encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
		token := encodedHeader + ".payload.signature"

		kid, err := ExtractKIDFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "example-kid", kid)
	})

	t.Run("invalid JWT format", func(t *testing.T) {
		token := "just.onepart"
		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JWT format")
	})

	t.Run("invalid base64 in header", func(t *testing.T) {
		token := "!!invalid!!base64.payload.signature"
		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode JWT header")
	})

	t.Run("invalid JSON in header", func(t *testing.T) {
		invalidJSON := "{bad json}"
		encoded := base64.RawURLEncoding.EncodeToString([]byte(invalidJSON))
		token := encoded + ".payload.signature"

		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse JWT header JSON")
	})

	t.Run("missing kid in header", func(t *testing.T) {
		header := `{"alg":"RS256"}`
		encoded := base64.RawURLEncoding.EncodeToString([]byte(header))
		token := encoded + ".payload.signature"

		_, err := ExtractKIDFromToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kid not found in JWT header")
	})
}

func TestFindJWKByKID(t *testing.T) {
	set := &JWKSet{
		Keys: []JWK{
			{Kid: "key1"},
			{Kid: "key2"},
		},
	}

	jwk, err := FindJWKByKID("key2", set)
	assert.NoError(t, err)
	assert.Equal(t, "key2", jwk.Kid)

	_, err = FindJWKByKID("missing", set)
	assert.ErrorIs(t, err, ErrJWKNotFound)
}

func TestFindJWKByKID_NilSet(t *testing.T) {
	_, err := FindJWKByKID("any", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWKSet is nil")
}

func TestBase64URLToBigInt(t *testing.T) {
	result, err := Base64URLToBigInt("AQAB")
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(65537), result)

	_, err = Base64URLToBigInt("!!invalid!!")
	assert.ErrorIs(t, err, ErrInvalidBase64URL)
}

func TestBuildRSAPublicKey_Success(t *testing.T) {
	n := big.NewInt(1234567890)
	e := big.NewInt(65537)

	pubKey, err := BuildRSAPublicKey(n, e)
	assert.NoError(t, err)
	assert.Equal(t, n, pubKey.N)
	assert.Equal(t, 65537, pubKey.E)
}

func TestBuildRSAPublicKey_ExponentTooLarge(t *testing.T) {
	e := new(big.Int).Lsh(big.NewInt(1), 63)
	n := big.NewInt(987654321)

	_, err := BuildRSAPublicKey(n, e)
	assert.ErrorIs(t, err, ErrExponentTooLarge)
}

func TestBuildRSAPublicKey_InvalidExponentValue(t *testing.T) {
	n := big.NewInt(1234567890)

	for _, exp := range []int64{0, 1} {
		e := big.NewInt(exp)
		_, err := BuildRSAPublicKey(n, e)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid exponent value")
	}
}
