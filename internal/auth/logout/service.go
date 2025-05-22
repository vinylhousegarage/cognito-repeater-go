package logout

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
)

type LogoutMetadata struct {
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

type LogoutURLProvider interface {
	GetLogoutURL(p config.MetadataURLProvider) (string, error)
}

type logoutClient struct {
	client httpclient.HTTPClient
}

func NewLogoutClient(client httpclient.HTTPClient) LogoutURLProvider {
	return &logoutClient{client: client}
}

func (s *logoutClient) GetLogoutURL(p config.MetadataURLProvider) (string, error) {
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
