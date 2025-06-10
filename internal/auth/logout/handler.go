package logout

import (
	"errors"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

// @Summary Redirect to Cognito logout
// @Description Logs the user out by redirecting to the Cognito end session endpoint.
// @Description This endpoint is intended to be accessed via a web browser as it performs a redirect.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /logout [get]
func NewLogoutHandler(
	p deps.LogoutHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadataURL := p.MetadataURL()
		endpoint, err := GetLogoutURL(metadataURL, c, logger)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, ErrUnexpectedStatusCode),
				errors.Is(err, ErrMissingEndSessionEndpoint):
				status = http.StatusBadGateway
				logger.Warn("GetLogoutURL returned an upstream error", zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("GetLogoutURL failed due to internal error", zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}
