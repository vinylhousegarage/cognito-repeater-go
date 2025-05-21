package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/config"
)

type LoginMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

type LoginEndpointProvider interface {
	GetLoginURL(p config.MetadataURLProvider) (string, error)
}

type loginService struct {
	client *http.Client
}

func NewLoginService(client *http.Client) LoginEndpointProvider {
	return &loginService{client: client}
}

func (s *loginService) GetLoginURL(p config.MetadataURLProvider) (string, error) {
	url := p.MetadataURL()

	resp, err := s.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta LoginMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.AuthorizationEndpoint == "" {
		return "", fmt.Errorf("invalid metadata: authorization endpoint is empty")
	}

	return meta.AuthorizationEndpoint, nil
}
