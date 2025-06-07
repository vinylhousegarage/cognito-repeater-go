package simulatederrors

import (
	"net/http"

	"cognito-repeater-go/internal/simulatederrors/notfound"

	"go.uber.org/zap"
)

func RegisterErrorRoutes(mux *http.ServeMux, logger *zap.Logger) {
	mux.HandleFunc("/error/404", notfound.NewError404Handler(logger))
}
