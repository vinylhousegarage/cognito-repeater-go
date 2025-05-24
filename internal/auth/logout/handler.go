package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
)

func LogoutHandler(d deps.HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadataURL := deps.Config.MetadataURL(d.Confg)
		endpoint, err := GetLogoutURL(d.HTTPClient, metadataURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}
