package main

import (
	"log"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv := server.NewServer(cfg)

	log.Println("Listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
