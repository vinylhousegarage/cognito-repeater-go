package health

import (
	"net/http"

	"cognito-repeater-go/internal/health/ping"
)

func RegisterHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", ping.NewPingHandler)
}
