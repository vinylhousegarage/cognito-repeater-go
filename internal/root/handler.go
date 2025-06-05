package root

import (
	"net/http"
)

// @Summary Redirect to login
// @Description Redirects the root path (/) to the login page
// @Tags root
// @Produce plain
// @Success 302 {string} string "Found"
// @Router / [get]
func NewRootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
