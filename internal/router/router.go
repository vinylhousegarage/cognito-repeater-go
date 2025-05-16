package router

import (
    "net/http"

    "cognito-repeater-go/internal/root"
)

func NewRouter() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/ping", root.PingHandler)

    return mux
}
