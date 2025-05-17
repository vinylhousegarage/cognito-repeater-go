package root

import "net/http"

func RegisterRootRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metadata", MetadataHandler)
	mux.HandleFunc("/ping", PingHandler)
}
