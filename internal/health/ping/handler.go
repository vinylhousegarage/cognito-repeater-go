package ping

import (
	"fmt"
	"net/http"
)

// @Summary Health check endpoint
// @Description Returns "pong" to indicate the server is alive
// @Tags health
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func NewPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, "pong"); err != nil {
		fmt.Printf("failed to write response: %v\n", err)
	}
}
