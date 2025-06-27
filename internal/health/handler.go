package health

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func NewHealthHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("responding to /health")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := HealthResponse{Status: "healthy"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write health response", zap.Error(err))
		}
	}
}
