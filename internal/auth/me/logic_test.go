package me

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestExtractKIDFromToken(t *testing.T) {
	t.Parallel()

	t.Run("valid token with kid", func(t *testing.T) {
		t.Parallel()

		header := `{"kid":"example-kid"}`
		encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
		token := encodedHeader + ".payload.signature"

		kid, err := ExtractKIDFromToken(token, testhelpers.MockLogger)
		assert.NoError(t, err)
		assert.Equal(t, "example-kid", kid)
	})

	t.Run("invalid JWT format", func(t *testing.T) {
		t.Parallel()

		token := "just.onepart"
		_, err := ExtractKIDFromToken(token, testhelpers.MockLogger)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidJWTFormat)
	})

	t.Run("invalid base64 in header", func(t *testing.T) {
		t.Parallel()

		token := "!!invalid!!base64.payload.signature"
		_, err := ExtractKIDFromToken(token, testhelpers.MockLogger)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrFailedToDecodeJWTHeader)
	})

	t.Run("invalid JSON in header", func(t *testing.T) {
		t.Parallel()

		invalidJSON := "{bad json}"
		encoded := base64.RawURLEncoding.EncodeToString([]byte(invalidJSON))
		token := encoded + ".payload.signature"

		_, err := ExtractKIDFromToken(token, testhelpers.MockLogger)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrFailedToParseJWTHeader)
	})

	t.Run("missing kid in header", func(t *testing.T) {
		t.Parallel()

		header := `{"alg":"RS256"}`
		encoded := base64.RawURLEncoding.EncodeToString([]byte(header))
		token := encoded + ".payload.signature"

		_, err := ExtractKIDFromToken(token, testhelpers.MockLogger)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingKID)
	})
}

func TestFindJWKByKID(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()
		set := &JWKSet{
			Keys: []JWK{{Kid: "key1"}, {Kid: "key2"}},
		}
		jwk, err := FindJWKByKID("key2", set, testhelpers.MockLogger)
		assert.NoError(t, err)
		assert.Equal(t, "key2", jwk.Kid)
	})

	t.Run("set is nil", func(t *testing.T) {
		t.Parallel()

		_, err := FindJWKByKID("any", nil, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrJWKSetNil)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		set := &JWKSet{
			Keys: []JWK{{Kid: "key1"}, {Kid: "key2"}},
		}
		_, err := FindJWKByKID("missing", set, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrJWKNotFound)
	})
}

func TestBase64URLToBigInt(t *testing.T) {
	t.Parallel()

	t.Run("valid input", func(t *testing.T) {
		t.Parallel()

		result, err := Base64URLToBigInt("AQAB")
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(65537), result)
	})

	t.Run("invalid base64", func(t *testing.T) {
		t.Parallel()

		_, err := Base64URLToBigInt("!!invalid!!")
		assert.ErrorIs(t, err, ErrInvalidBase64URL)
	})
}

func TestBuildRSAPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		n := big.NewInt(1234567890)
		e := big.NewInt(65537)

		pubKey, err := BuildRSAPublicKey(n, e)
		assert.NoError(t, err)
		assert.Equal(t, n, pubKey.N)
		assert.Equal(t, 65537, pubKey.E)
	})

	t.Run("exponent too large", func(t *testing.T) {
		t.Parallel()
		e := new(big.Int).Lsh(big.NewInt(1), 63)
		n := big.NewInt(987654321)

		_, err := BuildRSAPublicKey(n, e)
		assert.ErrorIs(t, err, ErrExponentTooLarge)
	})

	t.Run("invalid exponent values", func(t *testing.T) {
		t.Parallel()
		n := big.NewInt(1234567890)

		for _, exp := range []int64{0, 1} {
			exp := exp
			t.Run(fmt.Sprintf("e=%d", exp), func(t *testing.T) {
				t.Parallel()
				e := big.NewInt(exp)
				_, err := BuildRSAPublicKey(n, e)
				assert.ErrorIs(t, err, ErrInvalidExponent)
			})
		}
	})
}

func TestJWKToRSAPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		jwk := &JWK{
			N: "sXch4w", // base64url for []byte{0xb1, 0x77, 0x21, 0xe3}
			E: "AQAB",   // 65537
		}

		pubKey, err := JWKToRSAPublicKey(jwk, zap.NewNop())
		assert.NoError(t, err)
		assert.NotNil(t, pubKey)
		assert.Equal(t, 65537, pubKey.E)
		assert.Equal(t, new(big.Int).SetBytes([]byte{0xb1, 0x77, 0x21, 0xe3}), pubKey.N)
	})

	t.Run("invalid N", func(t *testing.T) {
		t.Parallel()

		jwk := &JWK{
			N: "!!invalid!!",
			E: "AQAB",
		}

		_, err := JWKToRSAPublicKey(jwk, zap.NewNop())
		assert.ErrorIs(t, err, ErrInvalidN)
	})

	t.Run("invalid E", func(t *testing.T) {
		t.Parallel()

		jwk := &JWK{
			N: "sXch4w",
			E: "!!invalid!!",
		}

		_, err := JWKToRSAPublicKey(jwk, zap.NewNop())
		assert.ErrorIs(t, err, ErrInvalidE)
	})
}
