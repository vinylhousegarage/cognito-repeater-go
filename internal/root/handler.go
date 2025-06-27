package root

import (
	"net/http"

	"go.uber.org/zap"
)

func NewRootHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("redirecting to /login", zap.String("path", r.URL.Path))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
