package me

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type JWKSMetadata struct {
	JWKSURI string `json:"jwks_uri"`
}

func GetJWKSURI(metadataURL string, client httpclient.HTTPClient, logger *zap.Logger) (string, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		logger.Error("failed to create metadata request", zap.String("url", metadataURL), zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch metadata", zap.String("url", metadataURL), zap.Error(err))
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close metadata response body", zap.Error(cerr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			logger.Warn("failed to read metadata response body", zap.Error(readErr))
		} else {
			logger.Error("unexpected status code from metadata endpoint",
				zap.Int("status", resp.StatusCode),
				zap.ByteString("body", body),
			)
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta JWKSMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.JWKSURI == "" {
		logger.Error("jwks_uri is missing in metadata")
		return "", fmt.Errorf("invalid metadata: jwks_uri is empty")
	}

	logger.Info("jwks_uri retrieved successfully", zap.String("jwks_uri", meta.JWKSURI))
	return meta.JWKSURI, nil
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

func FetchJWKSet(jwksURL string, client httpclient.HTTPClient, logger *zap.Logger) (*JWKSet, error) {
	req, err := http.NewRequest("GET", jwksURL, nil)
	if err != nil {
		logger.Error("failed to create jwks request", zap.String("url", jwksURL), zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch jwks", zap.String("url", jwksURL), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch jwks: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close jwks response body", zap.Error(cerr))
		}
	}()

	var set JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		logger.Error("failed to parse jwks JSON", zap.Error(err))
		return nil, fmt.Errorf("failed to parse jwks JSON: %w", err)
	}

	logger.Info("jwks successfully fetched", zap.Int("key_count", len(set.Keys)))
	return &set, nil
}
