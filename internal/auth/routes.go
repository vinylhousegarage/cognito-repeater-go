package auth

import "net/http"

func RegisterAuthRoutes(mux *http.ServeMux, p config.MetadataURLProvider) {
	mux.HandleFunc("/logout", LogoutHandler(p))
	mux.HandleFunc("/logout/redirect", LogoutRedirectHandler)
}
