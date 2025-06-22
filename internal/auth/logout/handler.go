package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Log the user out (redirect)
// @Description Redirects the browser to Cognito’s `end_session_endpoint`, passing the provided
// @Description `id_token_hint` so Cognito can identify the session.
// @Description The user’s browser is expected to call this endpoint directly.
// @Description
// @Description Example:
// @Description   GET /logout?id_token_hint=ID_TOKEN_VALUE
// @Tags auth
// @Produce plain
// @Param id_token_hint query string true "ID token (JWT) to supply as id_token_hint"
// @Success 302 {string} string "Found: Location header points to Cognito"
// @Failure 400 {object} response.ErrorResponse "Missing or invalid query parameter"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway (Cognito unreachable)"
// @Router /logout [get]
func NewLogoutHandler(
	p deps.LogoutHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idToken := r.URL.Query().Get("id_token_hint")
		if idToken == "" {
			logger.Warn("missing id_token_hint in query")
			response.WriteErrorResponse(w, ErrMissingIDTokenHint, logger)
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
			zap.String("client_id", p.UserPoolClientIDValue()),
			zap.String("logout_uri", p.LogoutURIValue()),
			zap.String("id_token_hint", utils.Mask(idToken)),
		)

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
