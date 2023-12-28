package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"prompt-generator/internal/api" // Replace with your actual import path
	"prompt-generator/internal/db"  // Replace with your actual import path
)

func main() {
	// Initialize database connection
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.DB.Close()

	// Set up the router
	r := mux.NewRouter()

	// Define route handlers
	r.HandleFunc("/api/prompts", api.GetPromptsHandler).Methods("GET")
	r.HandleFunc("/api/prompts", api.PostPromptHandler).Methods("POST")
	r.HandleFunc("/api/prompts/{id}", api.GetPromptHandler).Methods("GET")
	r.HandleFunc("/api/prompts/{id}", api.DeletePromptHandler).Methods("DELETE")

	// Determine the port to listen on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
