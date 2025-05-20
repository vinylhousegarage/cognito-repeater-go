package logout

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/config"
)

type LogoutMetadata struct {
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

type LogoutEndpointProvider interface {
	GetLogoutURL(p MetadataURLProvider) (string, error)
}

type logoutService struct {
	client *http.Client
}

func NewLogoutService(client *http.Client) LogoutEndpointProvider {
	return &logoutService{client: client}
}

func (s *logoutService) GetLogoutURL(p config.MetadataURLProvider) (string, error) {
	url := p.MetadataURL()

	resp, err := s.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta LogoutMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.EndSessionEndpoint == "" {
		return "", fmt.Errorf("invalid metadata: end_session_endpoint is empty")
	}

	return meta.EndSessionEndpoint, nil
}
