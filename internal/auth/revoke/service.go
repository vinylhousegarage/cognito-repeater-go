package revoke

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cognito-repeater-go/internal/httpclient"
)

type RevokeMetadata struct {
	RevokeEndpoint string `json:"revocation_endpoint"`
}

func GetRevokeURL(metadataURL string, client httpclient.HTTPClient) (string, error) {
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta RevokeMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.RevokeEndpoint == "" {
		return "", fmt.Errorf("revocation_endpoint is missing in metadata")
	}

	return meta.RevokeEndpoint, nil
}

func SendRevokeRequest(
	revokeURL string,
	cli httpclient.HTTPClient,
	refreshToken string,
	clientID string,
	clientSecret string,
) (*http.Response, error) {
	form := url.Values{}
	form.Set("token", refreshToken)
	form.Set("token_type_hint", "refresh_token")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("client credentials must not be empty")
	}

	req, err := http.NewRequest(http.MethodPost, revokeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create revoke request: %w", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call revoke endpoint: %w", err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("failed to close response body: %v\n", cerr)
		}
	}()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %w", readErr)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("revoke request failed: status=%s, body=%s", resp.Status, string(body))
	}

	return resp, nil
}
