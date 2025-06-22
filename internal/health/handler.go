package health

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// @Summary Health check endpoint
// @Description Returns "healthy" to indicate the server is alive
// @Tags health
// @Produce json
// @Success 200 {object} health.HealthResponse "healthy message"
// @Router /health [get]
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
