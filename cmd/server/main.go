package main

import (
	"log"
	"net/http"
	"wifi-go-backend/config"
	"wifi-go-backend/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	cfg := config.Load()
	router := routes.SetupRouter(cfg)
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
