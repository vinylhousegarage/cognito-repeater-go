package callback

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func NewCallbackHandler(
	p deps.CallbackHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := ValidateCallbackRequest(r)
		if err != nil {
			logger.Warn("invalid callback request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metadataURL := p.MetadataURL()
		tokenEndpoint, err := GetCallbackURL(metadataURL, c, logger)
		if err != nil {
			logger.Error("failed to get token endpoint", zap.String("metadataURL", metadataURL), zap.Error(err))
			http.Error(w, "failed to get token endpoint", http.StatusInternalServerError)
			return
		}

		bodyStr := BuildTokenRequestBody(code, p)
		req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(bodyStr))
		if err != nil {
			logger.Error("failed to create request to token endpoint", zap.String("tokenEndpoint", tokenEndpoint), zap.Error(err))
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", BuildBasicAuthHeader(p))

		resp, err := c.Do(req)
		if err != nil {
			logger.Error("token request failed", zap.Error(err))
			http.Error(w, "failed to send token request", http.StatusBadGateway)
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				logger.Warn("failed to close response body", zap.Error(cerr))
			}
		}()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			logger.Warn("token endpoint returned non-200",
				zap.Int("status", resp.StatusCode),
				zap.ByteString("body", body),
			)
			http.Error(w, "token endpoint returned error", http.StatusBadGateway)
			return
		}

		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			logger.Error("failed to decode token response", zap.Error(err))
			http.Error(w, "failed to decode token response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		logger.Info("token successfully issued",
			zap.String("client_id", p.UserPoolClientIDValue()),
			zap.String("token_type", tokenResp.TokenType),
		)

		if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
			logger.Error("failed to write token response", zap.Error(err))
			http.Error(w, "failed to encode token response", http.StatusInternalServerError)
			return
		}
	}
}
