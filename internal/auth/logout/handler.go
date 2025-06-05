package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
)

// @Summary Redirect to Cognito logout
// @Description Logs the user out by redirecting to the Cognito end session endpoint.
// @Description This endpoint is intended to be accessed via a web browser as it performs a redirect.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /logout [get]
func NewLogoutHandler(p deps.LogoutHandlerProvider, c httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadataURL := p.MetadataURL()
		endpoint, err := GetLogoutURL(metadataURL, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}
