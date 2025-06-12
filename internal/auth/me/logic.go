package me

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"strings"

	"go.uber.org/zap"
)

func ExtractKIDFromToken(token string, logger *zap.Logger) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidJWTFormat
	}

	headerSegment := parts[0]
	decodedHeader, err := base64.RawURLEncoding.DecodeString(headerSegment)
	if err != nil {
		logger.Warn("failed to decode JWT header", zap.String("segment", headerSegment), zap.Error(err))
		return "", ErrFailedToDecodeJWTHeader
	}

	var header struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(decodedHeader, &header); err != nil {
		logger.Warn("failed to parse JWT header JSON",
			zap.ByteString("decodedHeader", decodedHeader),
			zap.Error(err),
		)
		return "", ErrFailedToParseJWTHeader
	}

	if header.Kid == "" {
		return "", ErrMissingKID
	}

	return header.Kid, nil
}

func FindJWKByKID(kid string, set *JWKSet, logger *zap.Logger) (*JWK, error) {
	if set == nil {
		return nil, ErrJWKSetNil
	}

	for _, key := range set.Keys {
		if key.Kid == kid {
			return &key, nil
		}
	}
	logger.Warn("kid not found in JWK set", zap.String("kid", kid))
	return nil, ErrJWKNotFound
}

func Base64URLToBigInt(b64 string) (*big.Int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(b64)
	if err != nil {
		return nil, ErrInvalidBase64URL
	}

	n := new(big.Int).SetBytes(decoded)
	return n, nil
}

func BuildRSAPublicKey(n *big.Int, e *big.Int) (*rsa.PublicKey, error) {
	if !e.IsInt64() {
		return nil, ErrExponentTooLarge
	}
	eInt := int(e.Int64())

	if eInt <= 1 {
		return nil, ErrInvalidExponent
	}

	return &rsa.PublicKey{
		N: n,
		E: eInt,
	}, nil
}

func JWKToRSAPublicKey(jwk *JWK, logger *zap.Logger) (*rsa.PublicKey, error) {
	n, err := Base64URLToBigInt(jwk.N)
	if err != nil {
		logger.Error("failed to convert modulus from base64", zap.String("n", jwk.N), zap.Error(err))
		return nil, ErrInvalidN
	}
	e, err := Base64URLToBigInt(jwk.E)
	if err != nil {
		logger.Error("failed to convert exponent from base64", zap.String("e", jwk.E), zap.Error(err))
		return nil, ErrInvalidE
	}
	return BuildRSAPublicKey(n, e)
}
