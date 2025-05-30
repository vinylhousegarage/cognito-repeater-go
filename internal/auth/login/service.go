package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/httpclient"
)

type LoginMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

func GetLoginURL(metadataURL string, client httpclient.HTTPClient) (string, error) {
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

	var meta LoginMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.AuthorizationEndpoint == "" {
		return "", fmt.Errorf("invalid metadata: authorization_endpoint is empty")
	}

	return meta.AuthorizationEndpoint, nil
}
