package me

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
}

func ParseAndVerifyJWT(
	idToken string,
	pubKey *rsa.PublicKey,
	expectedIss, expectedAud string,
	now time.Time,
) (*MyClaims, error) {
  claims := &MyClaims{}

  token, err := jwt.ParseWithClaims(
    idToken,
    claims,
    func(token *jwt.Token) (interface{}, error) {
      if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
        return nil, ErrInvalidSigningAlg
      }
      return pubKey, nil
    },
    jwt.WithTime(now),
  )

	if err != nil || !token.Valid {
		return nil, ErrJWTParseFailed
	}

	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(now) {
		return nil, ErrTokenExpired
	}

	if claims.Issuer != expectedIss {
		return nil, ErrInvalidIssuer
	}

	if len(claims.Audience) == 0 {
		return nil, ErrMissingAudience
	}

	validAud := false
	for _, aud := range claims.Audience {
		if aud == expectedAud {
			validAud = true
			break
		}
	}
	if !validAud {
		return nil, ErrInvalidAudience
	}

	if claims.Subject == "" {
		return nil, ErrMissingSubject
	}

	return claims, nil
}
