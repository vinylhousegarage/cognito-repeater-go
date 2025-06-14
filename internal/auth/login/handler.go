package login

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Redirect to Cognito login
// @Description Initiates an OAuth2 login by redirecting to the Cognito authorization endpoint.
// @Description Sets a state cookie and is intended for browser access.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found: Redirect with Location header"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
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
			response.WriteErrorResponse(w, err, logger)
			return
		}

		url, err := BuildLoginURL(p, endpoint, state)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		logger.Info("redirecting to Cognito login", zap.String("url", url), zap.String("state", state))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
