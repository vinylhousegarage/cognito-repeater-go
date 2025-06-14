package notfound

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// @Summary Simulate 404 Not Found
// @Description Returns a simulated 404 Not Found error response
// @Tags error
// @Produce json
// @Failure 404 {object} notfound.ErrorSimulationResponse "Simulated 404 Not Found"
// @Router /error/404 [get]
func NewError404Handler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		resp := ErrorSimulationResponse{Message: "Simulated 404 Not Found"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write 404 simulation response", zap.Error(err))
		}
	}
}
