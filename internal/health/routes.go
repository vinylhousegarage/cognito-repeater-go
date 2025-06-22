package health

import (
	"net/http"

	"go.uber.org/zap"
)

func RegisterHealthRoutes(mux *http.ServeMux, logger *zap.Logger) {
	mux.HandleFunc("/health", NewHealthHandler(logger))
}
