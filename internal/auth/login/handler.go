package login

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
)

// @Summary Redirect to Cognito login
// @Description Initiates OAuth2 login by redirecting the user to the Cognito authorization endpoint.
// @Description This endpoint is intended to be accessed via a web browser as it sets a state cookie and performs a redirect.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [get]
func NewLoginHandler(p deps.LoginHandlerProvider, c httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := GenerateState()
		http.SetCookie(w, BuildStateCookie(state))

		metadataURL := p.MetadataURL()
		endpoint, err := GetLoginURL(metadataURL, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url, err := BuildLoginURL(p, endpoint, state)
		if err != nil {
			http.Error(w, "invalid login URL", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
