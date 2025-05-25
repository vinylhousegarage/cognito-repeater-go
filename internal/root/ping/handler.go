package ping

import (
	"fmt"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, "pong"); err != nil {
		fmt.Printf("failed to write response: %v\n", err)
	}
}
