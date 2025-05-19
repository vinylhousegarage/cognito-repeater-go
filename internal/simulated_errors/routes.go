package errors

import (
	"net/http"

	"cognito-repeater-go/internal/simulated_errors/404"
)

func RegisterErrorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/error/404", Error404Handler)
}
