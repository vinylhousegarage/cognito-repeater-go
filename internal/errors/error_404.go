package errors

import (
	"io"
	"log"
	"net/http"
)

// Error404Handler simulates a 404 Not Found error response.
func Error404Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)

	if _, err := io.WriteString(w, "not found 404"); err != nil {
		log.Printf("failed to write 404 message: %v", err)
	}
}
