package revoke

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type RevokeMetadata struct {
	RevocationEndpoint string `json:"revocation_endpoint"`
}

func GetRevokeURL(
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
		logger.Error("unexpected response from metadata",
			zap.Int("status", resp.StatusCode),
			zap.ByteString("body", body),
		)
		return "", ErrUnexpectedMetadataStatusCode
	}

	var meta RevokeMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", ErrFailedToDecodeMetadata
	}

	if meta.RevocationEndpoint == "" {
		logger.Error("missing revocation_endpoint")
		return "", ErrMissingRevocationEndpoint
	}

	logger.Info("revocation_endpoint retrieved successfully", zap.String("revocation_endpoint", meta.RevocationEndpoint))
	return meta.RevocationEndpoint, nil
}

func SendRevokeRequest(
	revokeURL string,
	cli httpclient.HTTPClient,
	refreshToken string,
	clientID string,
	clientSecret string,
	logger *zap.Logger,
) (*http.Response, error) {
	form := url.Values{}
	form.Set("token", refreshToken)
	form.Set("token_type_hint", "refresh_token")

	if clientID == "" || clientSecret == "" {
		logger.Error("missing client credentials",
			zap.String("client_id", utils.Mask(clientID)),
			zap.String("client_secret", utils.Mask(clientSecret)),
		)
		return nil, ErrMissingClientCredentials
	}

	req, err := http.NewRequest(http.MethodPost, revokeURL, strings.NewReader(form.Encode()))
	if err != nil {
		logger.Error("failed to create revoke request", zap.Error(err))
		return nil, ErrFailedToCreateRevokeRequest
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cli.Do(req)
	if err != nil {
		logger.Error("failed to send revoke request", zap.Error(err))
		return nil, ErrFailedToSendRevokeRequest
	}

	return resp, nil
}
