package auth

import (
	"cognito-repeater-go/internal/config"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, p config.MetadataURLProvider) {
	mux.HandleFunc("/logout", LogoutHandler(p))
	mux.HandleFunc("/logout/redirect", LogoutRedirectHandler)
}
