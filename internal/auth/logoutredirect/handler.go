package logoutredirect

import (
	"net/http"

	"go.uber.org/zap"
)

func NewLogoutRedirectHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		n, err := w.Write([]byte("Logout successful"))
		if err != nil {
			logger.Info("write failed", zap.Error(err))
			return
		}
		logger.Info("response body written", zap.Int("bytes", n))
	}
}
