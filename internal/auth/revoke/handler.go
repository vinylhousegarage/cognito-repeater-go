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

func NewRevokeHandler(
	p deps.RevokeHandlerProvider,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		refreshToken, err := utils.ExtractFormValue(r, logger)
		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			response.WriteErrorResponse(w, err, logger)
			return
		}

		revokeURL, err := GetRevokeURL(p.MetadataURL(), cli, logger)
		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			response.WriteErrorResponse(w, err, logger)
			return
		}

		clientSecret := strings.TrimSpace(p.ClientSecretValue())
		clientID := strings.TrimSpace(p.UserPoolClientIDValue())

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientID, clientSecret, logger)
		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			response.WriteErrorResponse(w, err, logger)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				logger.Warn("failed to close response body", zap.Error(err))
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			response.WriteErrorResponse(w, ErrUnexpectedRevokeStatusCode, logger)
			return
		}

		logger.Info("revocation succeeded",
			zap.String("status", resp.Status),
			zap.String("client_id", utils.Mask(clientID)),
		)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
	}
}
