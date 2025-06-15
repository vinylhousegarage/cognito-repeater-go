package openapi

import (
	"net/http"

	"go.uber.org/zap"
)

func RegisterOpenAPIRoutes(mux *http.ServeMux, logger *zap.Logger) {
	spec, err := LoadOpenAPISpec("openapi.yaml", logger)
	if err != nil {
		logger.Fatal("failed to load OpenAPI spec", zap.Error(err))
	}
	mux.HandleFunc("/openapi.json", NewOpenAPIHandler(spec, logger))
}
