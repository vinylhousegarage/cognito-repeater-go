package openapi

import (
	"net/http"

	"go.uber.org/zap"
)

func RegisterOpenAPIRoutes(mux *http.ServeMux, logger *zap.Logger) {
	spec, err := LoadEmbeddedOpenAPISpec(logger)
	if err != nil {
		logger.Fatal("failed to load embedded OpenAPI spec", zap.Error(err))
	}
	mux.HandleFunc("/openapi.json", NewOpenAPIHandler(spec, logger))
}
