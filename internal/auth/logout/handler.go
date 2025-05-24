package logout

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
)

func LogoutHandler(d deps.HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := GetLogoutURL(d.HTTPClient, d.Config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}
