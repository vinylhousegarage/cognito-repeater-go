package revoke

import (
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
// @Description This will invalidate the refresh token and prevent further token refresh or reuse,
// @Description effectively ending the user session.
// @Description This endpoint accepts a form-encoded POST request and should be called from a secure backend environment.
// @Description
// @Description Example form body:
// @Description   token=REFRESH_TOKEN_VALUE
// @Description   Content-Type: application/x-www-form-urlencoded
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param token formData string true "Refresh token to be revoked"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
// @Router /revoke [post]
func NewRevokeHandler(
	p deps.RevokeHandlerProvider,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := utils.ExtractFormValue(r, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		revokeURL, err := GetRevokeURL(p.MetadataURL(), cli, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		clientSecret := strings.TrimSpace(p.ClientSecretValue())
		clientID := strings.TrimSpace(p.UserPoolClientIDValue())

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientID, clientSecret, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				logger.Warn("failed to close response body", zap.Error(err))
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			response.WriteErrorResponse(w, ErrUnexpectedRevokeStatusCode, logger)
			return
		}

		logger.Info("revocation succeeded",
			zap.String("status", resp.Status),
			zap.String("client_id", utils.Mask(clientID)),
		)
		w.WriteHeader(http.StatusNoContent)
	}
}
