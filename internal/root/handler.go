package root

import (
	"net/http"

	"go.uber.org/zap"
)

// @Summary Redirect to login
// @Description Redirects the root path (/) to the login page.
// @Description This endpoint is intended to be accessed via a web browser.
// @Description It performs an HTTP redirect (status code 302) to /login.
// @Tags system
// @Produce plain
// @Success 302 {string} string "Found: Redirect with Location header"
// @Router / [get]
func NewRootHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("redirecting to /login", zap.String("path", r.URL.Path))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
