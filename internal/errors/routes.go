package errors

import "net/http"

func RegisterErrorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/error/404", Error404Handler)
}
