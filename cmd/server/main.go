package main

import (
	"log"
	"net/http"

	"github.com/fsipic/go-rivo/internal/generator"
	"github.com/fsipic/go-rivo/pkg/api"
)

func main() {
	log.Println("Starting the application...")
	mux := api.SetupRoutes()

	go generator.StartPriceFluctuation()

	log.Println("Server is about to start...")
	err := http.ListenAndServeTLS(":8443", "certs/server.crt", "certs/server.key", mux)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
