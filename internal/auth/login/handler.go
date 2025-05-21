package login

import (
	"net/http"

	"cognito-repeater-go/internal/config"
)

func LoginHandler(provider LoginEndpointProvider, p config.MetadataURLProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := provider.GetLoginURL(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}
