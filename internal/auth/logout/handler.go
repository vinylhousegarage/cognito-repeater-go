package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

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
