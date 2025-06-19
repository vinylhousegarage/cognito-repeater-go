package logoutredirect

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func NewLogoutRedirectHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp := LogoutResponse{Message: "Logout successful"}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to encode JSON", zap.Error(err))
			http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
			return
		}
	}
}
