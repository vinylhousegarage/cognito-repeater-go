package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
)

type LoginMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

type LoginURLProvider interface {
	GetLoginURL(p config.MetadataURLProvider) (string, error)
}

type loginClient struct {
	client httpclient.HTTPClient
}

func NewLoginClient(client httpclient.HTTPClient) LoginURLProvider {
	return &loginClient{client: client}
}

func (s *loginClient) GetLoginURL(p config.MetadataURLProvider) (string, error) {
	url := p.MetadataURL()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	var meta LoginMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.AuthorizationEndpoint == "" {
		return "", fmt.Errorf("invalid metadata: authorization endpoint is empty")
	}

	return meta.AuthorizationEndpoint, nil
}
