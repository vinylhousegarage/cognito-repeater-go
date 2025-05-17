package auth

import (
	"net/http"
)

// LogoutHandler redirects the user to the configured logout URL.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://stab.com/logout", http.StatusFound)
}
