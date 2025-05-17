package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/config"
)

func fetchMetadata(p config.MetadataURLProvider) (map[string]interface{}, error) {
	url := cfg.MetadataURL()

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode metadata JSON: %w", err)
	}

	return result, nil
}
