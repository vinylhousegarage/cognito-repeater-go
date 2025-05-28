package whoami

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/httpclient"
)

type UserinfoMetadata struct {
	UserinfoEndpoint string `json:"userinfo_endpoint"`
}

func GetUserinfoURL(client httpclient.HTTPClient, metadataURL string) (string, error) {
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

	var meta UserinfoMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.UserinfoEndpoint == "" {
		return "", fmt.Errorf("metadata response missing userinfo_endpoint")
	}

	return meta.UserinfoEndpoint, nil
}

func FetchUserinfo(client httpclient.HTTPClient, userinfoURL, token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", userinfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call userinfo endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return result, nil
}
