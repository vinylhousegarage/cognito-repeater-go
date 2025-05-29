package me

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/httpclient"
)

type JWKSMetadata struct {
	JWKSURI string `json:"jwks_uri"`
}

func GetJWKSURI(metadataURL string, client httpclient.HTTPClient) (string, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta JWKSMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.JWKSURI == "" {
		return "", fmt.Errorf("invalid metadata: jwks_uri is empty")
	}

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

func FetchJWKSet(jwksURL string, client httpclient.HTTPClient) (*JWKSet, error) {
	req, err := http.NewRequest("GET", jwksURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jwks: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	var set JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, fmt.Errorf("failed to parse jwks JSON: %w", err)
	}
	return &set, nil
}
