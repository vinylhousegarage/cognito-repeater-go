package auth

import "net/http"

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/logout/redirect", LogoutRedirectHandler)
}
