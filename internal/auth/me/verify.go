package me

import (
	"crypto/rsa"
	"fmt"
	"time"

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

	now := time.Now()

	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(now) {
		return nil, fmt.Errorf("token is expired")
	}

	if claims.Issuer != expectedIss {
		return nil, fmt.Errorf("unexpected issuer: %s", claims.Issuer)
	}

	if len(claims.Audience) == 0 {
		return nil, fmt.Errorf("audience claim is missing")
	}

	validAud := false
	for _, aud := range claims.Audience {
		if aud == expectedAud {
			validAud = true
			break
		}
	}
	if !validAud {
		return nil, fmt.Errorf("invalid audience: expected %q", expectedAud)
	}

	if claims.Subject == "" {
		return nil, fmt.Errorf("missing subject (sub)")
	}

	return claims, nil
}
