package logoutredirect

import (
	"net/http"

	"go.uber.org/zap"
)

func NewLogoutRedirectHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
