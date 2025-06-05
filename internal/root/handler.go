package root

import (
	"net/http"
)

// @Summary Redirect to login
// @Description Redirects the root path (/) to the login page.
// @Description This endpoint is intended to be accessed via a web browser.
// @Description It performs an HTTP redirect (status code 302) to /login.
// @Tags system
// @Produce plain
// @Success 302 {string} string "Found"
// @Router / [get]
func NewRootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
