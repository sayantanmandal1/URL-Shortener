package main

import (
	"log"

	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Create and start server
	srv := server.New(cfg)

	log.Printf("Starting URL Shortener API in %s mode", cfg.Server.Environment)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
