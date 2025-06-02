package root

import (
	"net/http"

	"cognito-repeater-go/internal/root/ping"
)

func RegisterRootRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", ping.NewPingHandler)
}
