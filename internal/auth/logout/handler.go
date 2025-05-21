package logout

import (
	"net/http"

	"cognito-repeater-go/internal/config"
)

func LogoutHandler(lup LogoutURLProvider, mup config.MetadataURLProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := lupp.GetLogoutURL(mup)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}
