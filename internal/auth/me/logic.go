package me

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidJWTFormat = errors.New("invalid JWT format")
	ErrMissingKID       = errors.New("kid not found in JWT header")
)

func ExtractKIDFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("%w", ErrInvalidJWTFormat)
	}

	headerSegment := parts[0]
	decodedHeader, err := base64.RawURLEncoding.DecodeString(headerSegment)
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT header: %w", err)
	}

	var header struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(decodedHeader, &header); err != nil {
		return "", fmt.Errorf("failed to parse JWT header JSON: %w", err)
	}

	if header.Kid == "" {
		return "", ErrMissingKID
	}

	return header.Kid, nil
}

var (
	ErrJWKNotFound = errors.New("JWK not found")
	ErrJWKSetNil   = errors.New("JWKSet is nil")
)

func FindJWKByKID(kid string, set *JWKSet) (*JWK, error) {
	if set == nil {
		return nil, ErrJWKSetNil
	}

	for _, key := range set.Keys {
		if key.Kid == kid {
			return &key, nil
		}
	}
	return nil, fmt.Errorf("JWK with kid %s not found: %w", kid, ErrJWKNotFound)
}
