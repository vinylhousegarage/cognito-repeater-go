package root

import (
	"net/http"

	"go.uber.org/zap"
)

func RegisterRootRoutes(mux *http.ServeMux, logger *zap.Logger) {
	mux.HandleFunc("/", NewRootHandler(logger))
}
