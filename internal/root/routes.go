package root

import "net/http"

func RegisterRootRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", PingHandler)
}
