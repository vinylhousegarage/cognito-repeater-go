package whoami

import (
	"encoding/json"
	"io"
	"net/http"

	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type UserinfoMetadata struct {
	UserinfoEndpoint string `json:"userinfo_endpoint"`
}

func GetUserinfoURL(
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
			logger.Warn("failed to close metadata response body", zap.Error(cerr))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			logger.Error("failed to read metadata response body", zap.Error(err))
			return nil, ErrFailedToReadMetadataResponse
	}

	if resp.StatusCode != http.StatusOK {
			logger.Error("unexpected response from metadata",
					zap.Int("status", resp.StatusCode),
					zap.ByteString("body", body),
			)
			return nil, ErrUnexpectedMetadataStatusCode
	}

	var meta UserinfoMetadata
	if err := json.Unmarshal(body, &meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return "", ErrFailedToDecodeMetadata
	}

	if meta.UserinfoEndpoint == "" {
		logger.Error("missing userinfo_endpoint")
		return "", ErrMissingUserinfoEndpoint
	}

	logger.Info("userinfo_endpoint retrieved successfully", zap.String("userinfo_endpoint", meta.UserinfoEndpoint))
	return meta.UserinfoEndpoint, nil
}

func FetchUserinfo(
	userinfoURL string,
	client httpclient.HTTPClient,
	token string,
	logger *zap.Logger,
) (*UserInfoResponse, error) {
	req, err := http.NewRequest("GET", userinfoURL, nil)
	if err != nil {
		logger.Error("failed to create userinfo request", zap.Error(err))
		return nil, ErrFailedToCreateUserinfoRequest
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch userinfo request", zap.String("url", userinfoURL), zap.Error(err))
		return "", ErrFailedToFetchUserinfoRequest
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close userinfo response body", zap.Error(cerr))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			logger.Error("failed to read userinfo response body", zap.Error(err))
			return nil, ErrFailedToReadUserinfoResponse
	}

	if resp.StatusCode != http.StatusOK {
			logger.Error("unexpected response from userinfo",
					zap.Int("status", resp.StatusCode),
					zap.ByteString("body", body),
			)
			return nil, ErrUnexpectedUserinfoStatusCode
	}

	var meta UserInfoResponse
	if err := json.Unmarshal(body, &meta); err != nil {
			logger.Warn("failed to decode userinfo body",
					zap.ByteString("body", body),
					zap.Error(err),
			)
			return nil, ErrFailedToParseUserinfoResponse
	}

	logger.Info("userinfo fetched successfully", zap.String("sub", meta.Sub))
	return &meta, nil
}
