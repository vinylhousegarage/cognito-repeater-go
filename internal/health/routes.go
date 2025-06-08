package health

import (
	"net/http"

	"cognito-repeater-go/internal/health/ping"

	"go.uber.org/zap"
)

func RegisterHealthRoutes(mux *http.ServeMux, logger *zap.Logger) {
	mux.HandleFunc("/ping", ping.NewPingHandler(logger))
}
