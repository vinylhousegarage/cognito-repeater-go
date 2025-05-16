package main

import (
	"log"

	"cognito-repeater-go/internal/server"
)

func main() {
	srv := server.NewServer()

	log.Println("Listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
