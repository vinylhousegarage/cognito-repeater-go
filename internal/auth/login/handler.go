package login

import (
	"errors"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

// @Summary Redirect to Cognito login
// @Description Initiates OAuth2 login by redirecting the user to the Cognito authorization endpoint.
// @Description This endpoint is intended to be accessed via a web browser as it sets a state cookie and performs a redirect.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [get]
func NewLoginHandler(
	p deps.LoginHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := GenerateState()
		http.SetCookie(w, BuildStateCookie(state))

		metadataURL := p.MetadataURL()
		endpoint, err := GetLoginURL(metadataURL, c, logger)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, ErrUnexpectedStatusCode),
				errors.Is(err, ErrMissingAuthorizationEndpoint):
				status = http.StatusBadGateway
				logger.Warn("invalid metadata received from upstream", zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("failed to retrieve or parse metadata", zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		url, err := BuildLoginURL(p, endpoint, state)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, ErrFailedToParseLoginURL):
				status = http.StatusBadGateway
				logger.Warn("failed to parse login URL", zap.String("endpoint", endpoint), zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("unexpected error while building login URL", zap.String("endpoint", endpoint), zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		logger.Info("redirecting to Cognito login", zap.String("url", url), zap.String("state", state))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
