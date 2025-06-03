package root

import (
	"net/http"

	"cognito-repeater-go/internal/root/rootroot"
)

func RegisterRootRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", rootroot.NewRootHandler)
}
