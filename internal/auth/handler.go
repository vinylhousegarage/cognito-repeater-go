package auth

import (
	"cognito-repeater-go/internal/config"
	"encoding/json"
	"log"
	"net/http"
)

// LogoutHandler redirects the user to the configured logout URL.
func LogoutHandler(p config.MetadataURLProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint, err := getLogoutEndpoint(p)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, endpoint, http.StatusFound)
	}
}

type LogoutResponse struct {
	Message string `json:"message"`
}

func LogoutRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := LogoutResponse{Message: "Logout successful"}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
