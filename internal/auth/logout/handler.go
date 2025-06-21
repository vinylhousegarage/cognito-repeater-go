package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Redirect to Cognito logout endpoint
// @Description Initiates a Cognito logout by redirecting the user's browser to the end_session_endpoint.
// @Description The request must include a valid ID token in the Authorization header as a Bearer token.
// @Description The token is extracted and passed to Cognito as `id_token_hint` according to the OIDC RP-Initiated Logout specification.
// @Description This endpoint is typically called from a frontend application via a browser to complete the logout flow.
// @Description
// @Description Example header:
// @Description   Authorization: Bearer ID_TOKEN_VALUE
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found: Redirect with Location header"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
// @Router /logout [get]
func NewLogoutHandler(
	p deps.LogoutHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idToken, err := utils.ExtractAuthHeaderToken(r)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		metadataURL := p.MetadataURL()
		endpoint, err := GetLogoutURL(metadataURL, c, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		redirectURL := buildLogoutRedirectURL(
			endpoint,
			p.UserPoolClientIDValue(),
			p.LogoutURIValue(),
			idToken,
		)

		logger.Info("redirecting to Cognito logout",
			zap.String("url", redirectURL),
			zap.String("client_id", p.UserPoolClientIDValue()),
			zap.String("post_logout_redirect_uri", p.LogoutURIValue()),
			zap.String("id_token_hint", utils.Mask(idToken)),
		)

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
