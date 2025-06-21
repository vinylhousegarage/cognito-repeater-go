package logout

import (
	"net/http"
	"net/url"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Redirect to Cognito logout
// @Description Logs the user out by redirecting to the Cognito end session endpoint.
// @Description This endpoint is intended to be accessed via a web browser as it performs a redirect.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found: Redirect with Location header"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
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
			response.WriteErrorResponse(w, err, logger)
			return
		}

		redirectURL := buildLogoutRedirectURL(
			endpoint,
			p.UserPoolClientIDValue(),
			url.QueryEscape(p.LogoutURIValue()),
		)

		logger.Info("redirecting to Cognito logout",
			zap.String("url", redirectURL),
			zap.String("client_id", p.UserPoolClientIDValue()),
			zap.String("post_logout_redirect_uri", url.QueryEscape(p.LogoutURIValue())),
		)

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
