package login

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type LoginMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

func GetLoginURL(
	metadataURL string,
	client httpclient.HTTPClient,
	logger *zap.Logger,
) (string, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		logger.Error("failed to create request", zap.String("url", metadataURL), zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch metadata", zap.String("url", metadataURL), zap.Error(err))
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close response body", zap.Error(cerr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error("unexpected status code from metadata endpoint",
			zap.Int("status", resp.StatusCode),
			zap.ByteString("body", body),
		)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta LoginMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.AuthorizationEndpoint == "" {
		logger.Error("metadata response missing authorization_endpoint")
		return "", fmt.Errorf("invalid metadata: authorization_endpoint is empty")
	}

	logger.Info("authorization_endpoint retrieved successfully", zap.String("authorization_endpoint", meta.AuthorizationEndpoint))
	return meta.AuthorizationEndpoint, nil
}
