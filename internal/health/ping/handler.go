package ping

import (
	// "encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// @Summary Health check endpoint
// @Description Returns "pong" to indicate the server is alive
// @Tags health
// @Produce json
// @Success 200 {object} ping.PingResponse "pong message"
// @Router /ping [get]
// func NewPingHandler(logger *zap.Logger) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger.Info("responding to /ping")

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		resp := PingResponse{Message: "pong"}
// 		if err := json.NewEncoder(w).Encode(resp); err != nil {
// 			logger.Error("failed to write ping response", zap.Error(err))
// 		}
// 	}
// }

func NewPingHandler(_ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("simulated crash")
	}
}
