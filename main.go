package main

import (
	"log"

	"backend-raw-http/config"
	"backend-raw-http/database"
	"backend-raw-http/handlers"

	"github.com/codetesla51/raw-http/server"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Connected to database successfully")

	// Create server
	srv := server.NewServer(cfg.ServerPort)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler()

	// Health check
	srv.Register("GET", "/health", func(req *server.Request) ([]byte, string) {
		return server.CreateResponseBytes("200", "application/json", "OK", []byte(`{"status":"healthy"}`))
	})

	// Project routes
	srv.Register("GET", "/api/projects", projectHandler.GetAllProjects)
	srv.Register("GET", "/api/projects/:slug", projectHandler.GetProjectBySlug)

	// Start server
	log.Printf("Server starting on %s", cfg.ServerPort)
	log.Printf("Base URL: %s", cfg.BaseURL)
	log.Println("Routes:")
	log.Println("  GET  /health")
	log.Println("  GET  /api/projects")
	log.Println("  GET  /api/projects/:slug")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
