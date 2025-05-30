package login

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
)

func LoginHandler(p deps.LoginHandlerProvider, c httpclient.HTTPClient) http.HandlerFunc {
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
