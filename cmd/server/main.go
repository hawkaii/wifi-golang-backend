package main

import (
	"log"
	"net/http"
	"wifi-go-backend/config"
	"wifi-go-backend/internal/routes"
)

func main() {
	cfg := config.Load()
	router := routes.SetupRouter(cfg)
	log.Fatal(http.ListenAndServe(":8080", router))
}
