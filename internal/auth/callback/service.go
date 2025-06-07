package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

func ValidateCallbackRequest(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		return "", errors.New("missing code")
	}

	if state == "" {
		return "", errors.New("missing state")
	}

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return "", errors.New("missing oauth_state cookie")
	}

	if state != cookie.Value {
		return "", errors.New("invalid state")
	}

	return code, nil
}

type callbackMetadata struct {
	TokenEndpoint string `json:"token_endpoint"`
}

func GetCallbackURL(
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

	var meta callbackMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.TokenEndpoint == "" {
		logger.Error("metadata response missing token_endpoint")
		return "", fmt.Errorf("invalid metadata: token_endpoint is empty")
	}

	logger.Info("token_endpoint retrieved successfully", zap.String("token_endpoint", meta.TokenEndpoint))
	return meta.TokenEndpoint, nil
}
