package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
)

func LogoutHandler(p deps.LogoutHandlerProvider, c httpclient.HTTPClient) http.HandlerFunc {
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
