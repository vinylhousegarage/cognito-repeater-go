package login

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
)

func LoginHandler(d deps.HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := GenerateState()
		http.SetCookie(w, BuildStateCookie(state))

		endpoint, err := GetLoginURL(d.HTTPClient, d.Config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url, err := BuildLoginURL(endpoint, state)
		if err != nil {
			http.Error(w, "invalid login URL", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
