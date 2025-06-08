package urlresolver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type Metadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	EndSessionEndpoint    string `json:"end_session_endpoint,omitempty"`
	JWKSURI               string `json:"jwks_uri,omitempty"`
	RevokeEndpoint        string `json:"revocation_endpoint,omitempty"`
	TokenEndpoint         string `json:"token_endpoint,omitempty"`
}

func FetchMetadata(metadataURL string, client httpclient.HTTPClient, logger *zap.Logger) (*Metadata, error) {
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		logger.Error("failed to create request", zap.String("url", metadataURL), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrFailedToCreateRequest, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to fetch metadata", zap.String("url", metadataURL), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrFailedToFetchMetadata, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Warn("failed to close response body", zap.Error(cerr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("metadata endpoint returned unexpected status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("%w: status=%d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var meta Metadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		logger.Error("failed to decode metadata JSON", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrFailedToDecodeMetadata, err)
	}

	if meta.AuthorizationEndpoint == "" {
		logger.Error("authorization_endpoint is missing in metadata")
		return nil, ErrMissingAuthorizationEndpoint
	}
	if meta.TokenEndpoint == "" {
		logger.Error("token_endpoint is missing in metadata")
		return nil, ErrMissingTokenEndpoint
	}
	if meta.EndSessionEndpoint == "" {
		logger.Error("end_session_endpoint is missing in metadata")
		return nil, ErrMissingEndSessionEndpoint
	}
	if meta.RevokeEndpoint == "" {
		logger.Error("revocation_endpoint is missing in metadata")
		return nil, ErrMissingRevocationEndpoint
	}
	if meta.JWKSURI == "" {
		logger.Error("jwks_uri is missing in metadata")
		return nil, ErrMissingJWKSURI
	}

	logger.Info("metadata fetched and validated successfully", zap.Any("metadata", meta))
	return &meta, nil
}
