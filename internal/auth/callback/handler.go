package callback

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
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
			status := http.StatusBadRequest
			utils.WritePlainError(w, status, err, logger)
			return
		}

		metadataURL := p.MetadataURL()
		tokenEndpoint, err := GetCallbackURL(metadataURL, c, logger)
		if err != nil {
			switch {
			case errors.Is(err, ErrFailedToCreateRequest),
				errors.Is(err, ErrFailedToFetchMetadata),
				errors.Is(err, ErrFailedToDecodeMetadata):

				status := http.StatusInternalServerError
				utils.WritePlainError(w, status, err, logger)

			case errors.Is(err, ErrUnexpectedStatusCode),
				errors.Is(err, ErrInvalidMetadataNoEndpoint):

				status := http.StatusBadGateway
				utils.WritePlainError(w, status, err, logger)

			default:
				status := http.StatusInternalServerError
				utils.WritePlainError(w, status, err, logger)
			}
			return
		}

		bodyStr := BuildTokenRequestBody(code, p)
		req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(bodyStr))
		if err != nil {
			status := http.StatusInternalServerError
			utils.WritePlainError(w, status, err, logger)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", BuildBasicAuthHeader(p))

		resp, err := c.Do(req)
		if err != nil {
			status := http.StatusBadGateway
			utils.WritePlainError(w, status, err, logger)
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				logger.Warn("failed to close response body", zap.Error(cerr))
			}
		}()

		if resp.StatusCode != http.StatusOK {
			body, readErr := io.ReadAll(resp.Body)
			status := http.StatusBadGateway

			fields := []zap.Field{
				zap.Int("upstream_status", resp.StatusCode),
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
			}

			if readErr != nil {
				fields = append(fields, zap.Error(readErr))
				logger.Warn(http.StatusText(status)+": failed to read error response body", fields...)
			} else {
				fields = append(fields, zap.ByteString("body", body))
				logger.Warn(http.StatusText(status)+": token endpoint returned non-200", fields...)
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, http.StatusText(status), status)
			return
		}

		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			status := http.StatusInternalServerError
			utils.writePlainError(w, status, err, logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		logger.Info("token successfully issued",
			zap.String("client_id", p.UserPoolClientIDValue()),
			zap.String("token_type", tokenResp.TokenType),
		)

		if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
			status := http.StatusInternalServerError
			utils.WritePlainError(w, status, err, logger)
			return
		}
	}
}
