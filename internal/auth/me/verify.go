package me

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
}

func ParseAndVerifyJWT(tokenStr string, pubKey *rsa.PublicKey, expectedIss, expectedAud string) (*MyClaims, error) {
	claims := &MyClaims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("JWT parse error: %w", err)
	}

	if err := claims.Valid(); err != nil {
		return nil, fmt.Errorf("claims validation failed: %w", err)
	}
	if claims.Issuer != expectedIss {
		return nil, fmt.Errorf("unexpected issuer: %s", claims.Issuer)
	}
	if !claims.VerifyAudience(expectedAud, true) {
		return nil, fmt.Errorf("invalid audience")
	}
	if claims.Subject == "" {
		return nil, fmt.Errorf("missing subject (sub)")
	}

	return claims, nil
}
