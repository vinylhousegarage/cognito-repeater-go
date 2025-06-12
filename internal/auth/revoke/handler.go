package revoke

import (
	"fmt"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Revoke refresh token
// @Description Revokes a refresh token by calling the Cognito revocation endpoint.
// @Description This endpoint expects a form-encoded POST request and should be called from a secure backend environment.
// @Description To revoke a refresh token, send a POST request to this endpoint with the following body:
// @Description token=REFRESH_TOKEN_VALUE
// @Description Content-Type: application/x-www-form-urlencoded
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param token formData string true "Refresh token to be revoked"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /revoke [post]
func NewRevokeHandler(
	p deps.RevokeHandlerProvider,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := utils.ExtractFormValue(r, logger)
		if err != nil {
			logger.Warn("failed to extract token from form", zap.Error(err))
			response.WriteJSONError(w, http.StatusBadRequest, err.Error(), logger)
			return
		}

		revokeURL, err := GetRevokeURL(p.MetadataURL(), cli, logger)
		if err != nil {
			switch {
			case errors.Is(err, ErrFailedToFetchMetadata),
				errors.Is(err, ErrFailedToDecodeMetadata),
				errors.Is(err, ErrUnexpectedStatusCode),
				errors.Is(err, ErrMissingRevocationEndpoint):
				logger.Warn("GetRevokeURL returned an upstream error", zap.Error(err))
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			default:
				logger.Error("GetRevokeURL failed due to internal error", zap.Error(err))
				response.WriteJSONError(w, http.StatusInternalServerError, err.Error(), logger)
			}
			return
		}

		clientSecret := strings.TrimSpace(p.ClientSecretValue())
		clientID := strings.TrimSpace(p.UserPoolClientIDValue())

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientID, clientSecret)
		if err != nil {
			switch {
			case errors.Is(err, ErrMissingClientCredentials):
				logger.Error("missing client credentials", zap.Error(err))
				response.WriteJSONError(w, http.StatusUnauthorized, err.Error(), logger)
			case errors.Is(err, ErrFailedToSendRevokeRequest):
				logger.Error("failed to send revoke request", zap.Error(err))
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			default:
				logger.Error("unexpected response from revocation_endpoint", zap.Error(err))
				response.WriteJSONError(w, http.StatusInternalServerError, err.Error(), logger)
			}
			return
		}
		if resp != nil {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					logger.Warn("failed to close response body", zap.Error(err))
				}
			}()
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			response.WriteJSONError(w, http.StatusBadGateway, fmt.Sprintf("revocation failed with status: %s", resp.Status), logger)
			return
		}

		logger.Info("revocation succeeded",
			zap.String("status", resp.Status),
			zap.String("client_id", utils.Mask(clientID)),
		)
		w.WriteHeader(http.StatusNoContent)
	}
}
