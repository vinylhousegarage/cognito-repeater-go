package whoami

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cognito-repeater-go/internal/httpclient"
)

type UserinfoMetadata struct {
	UserinfoEndpoint string `json:"userinfo_endpoint"`
}

func GetUserinfoURL(metadataURL string, client httpclient.HTTPClient) (string, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("failed to close response body: %v\n", cerr)
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read metadata response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status code: %d\n", resp.StatusCode)
		fmt.Printf("response body: %s\n", string(bodyBytes))
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta UserinfoMetadata
	if err := json.Unmarshal(bodyBytes, &meta); err != nil {
		fmt.Printf("failed to decode metadata body: %s\n", string(bodyBytes))
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.UserinfoEndpoint == "" {
		return "", fmt.Errorf("metadata response missing userinfo_endpoint")
	}

	return meta.UserinfoEndpoint, nil
}

func FetchUserinfo(userinfoURL string, client httpclient.HTTPClient, token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", userinfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call userinfo endpoint: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("failed to close userinfo response body: %v\n", cerr)
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("userinfo endpoint returned status: %d\n", resp.StatusCode)
		fmt.Printf("userinfo response body: %s\n", string(bodyBytes))
		return nil, fmt.Errorf("userinfo endpoint returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("failed to decode userinfo body: %s\n", string(bodyBytes))
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return result, nil
}
