package errors

import (
    "net/http"
    "io"
)

// Error404Handler simulates a 404 Not Found error response.
func Error404Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusNotFound)
    io.WriteString(w, "not found 404")
}
