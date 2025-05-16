package root

import (
	"encoding/json"
	"log"
	"net/http"
)

func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metadata := map[string]string{
		"login_endpoint":               "/login",
		"logout_endpoint":              "/logout",
		"verify_access_token_endpoint": "/token",
		"verify_userinfo_endpoint":     "/user",
		"health_check_endpoint":        "/ping",
		"simulate_404_endpoint":        "/error/404",
		"openapi_url":                  "/openapi.json",
	}

	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
