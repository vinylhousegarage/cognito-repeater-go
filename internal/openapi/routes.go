package openapi

import (
	"net/http"

	"go.uber.org/zap"
)

func RegisterOpenAPIRoutes(mux *http.ServeMux, logger *zap.Logger) {
	mux.HandleFunc("/openapi.json", NewOpenAPIHandler(logger))
}
