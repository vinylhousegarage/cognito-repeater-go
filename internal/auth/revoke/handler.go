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
		refreshToken, err := utils.ExtractFormValue(r)
		if err != nil {
			logger.Warn("failed to extract token from form", zap.Error(err))
			response.WriteJSONError(w, http.StatusBadRequest, err.Error(), logger)
			return
		}

		revokeURL, err := GetRevokeURL(p.MetadataURL(), cli)
		if err != nil {
			logger.Error("failed to get revoke URL", zap.String("metadata_url", p.MetadataURL()), zap.Error(err))
			response.WriteJSONError(w, http.StatusInternalServerError, "failed to resolve revocation endpoint", logger)
			return
		}

		clientSecret := p.ClientSecretValue()
		clientID := p.UserPoolClientIDValue()

		clientSecret = strings.TrimSpace(clientSecret)
		clientID = strings.TrimSpace(clientID)

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientID, clientSecret)
		if err != nil {
			logger.Error("failed to revoke token", zap.String("revoke_url", revokeURL), zap.Error(err))
			response.WriteJSONError(w, http.StatusInternalServerError, "failed to revoke token", logger)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				logger.Warn("failed to close response body", zap.Error(err))
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			response.WriteJSONError(w, http.StatusBadGateway, fmt.Sprintf("revocation failed with status: %s", resp.Status), logger)
			return
		}

		logger.Info("revocation succeeded", zap.String("status", resp.Status))
		w.WriteHeader(http.StatusNoContent)
	}
}
