package logout

import (
	"encoding/json"
	"io"
	"net/http"

	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type LogoutMetadata struct {
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

func GetLogoutURL(
	metadataURL string,
	client httpclient.HTTPClient,
	logger *zap.Logger,
) (string, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		logger.Error("failed to create request", zap.String("url", metadataURL), zap.Error(err))
		return "", ErrFailedToCreateRequest
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch metadata", zap.String("url", metadataURL), zap.Error(err))
		return "", ErrFailedToFetchMetadata
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
		return "", ErrUnexpectedStatusCode
	}

	var meta LogoutMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", ErrFailedToDecodeMetadata
	}

	if meta.EndSessionEndpoint == "" {
		logger.Error("missing end_session_endpoint in metadata response")
		return "", ErrMissingEndSessionEndpoint
	}

	logger.Info("end_session_endpoint retrieved successfully", zap.String("end_session_endpoint", meta.EndSessionEndpoint))
	return meta.EndSessionEndpoint, nil
}
