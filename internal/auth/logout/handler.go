package logout

import (
	"net/http"

	"cognito-repeater-go/internal/config"
)

func LogoutHandler(provider LogoutEndpointProvider) func(config.MetadataURLProvider) http.HandlerFunc {
	return func(p config.MetadataURLProvider) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			endpoint, err := provider.GetLogoutURL(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, endpoint, http.StatusFound)
		}
	}
}
