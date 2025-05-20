package simulatederrors

import (
	"net/http"

	"cognito-repeater-go/internal/simulatederrors/notfound"
)

func RegisterErrorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/error/404", notfound.Error404Handler)
}
