package root

import (
	"net/http"

	"cognito-repeater-go/internal/root/ping"
	"cognito-repeater-go/internal/root/rootroot"
)

func RegisterRootRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", ping.NewPingHandler)
	mux.HandleFunc("/", rootroot.NewRootHandler)
}
