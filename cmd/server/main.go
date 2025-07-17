package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"twitchvod/internal/twitch"
)

type Response struct {
	Success bool   `json:"success"`
	URL     string `json:"url,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get secret from environment variable
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable is required")
	}

	// Get client ID from environment variable
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	if clientID == "" {
		log.Fatal("TWITCH_CLIENT_ID environment variable is required")
	}

	// Create Twitch client
	twitchClient := twitch.New(clientID, "https://gql.twitch.tv/gql")

	// Define handler
	http.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests
		if r.Method != "GET" {
			respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get query parameters
		url := r.URL.Query().Get("url")
		providedSecret := r.URL.Query().Get("secret")

		// Validate secret
		if providedSecret != secret {
			respondWithError(w, "Invalid secret", http.StatusUnauthorized)
			return
		}

		// Validate URL parameter
		if url == "" {
			respondWithError(w, "URL parameter is required", http.StatusBadRequest)
			return
		}

		// Extract video ID
		videoID, err := twitchClient.ExtractVideoID(url)
		if err != nil {
			respondWithError(w, fmt.Sprintf("Error extracting video ID: %v", err), http.StatusBadRequest)
			return
		}

		// Get streaming URL
		streamingURL, err := twitchClient.GetVideoInfo(videoID)
		if err != nil {
			respondWithError(w, fmt.Sprintf("Error getting video info: %v", err), http.StatusInternalServerError)
			return
		}

		// Return success response
		response := Response{
			Success: true,
			URL:     streamingURL,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response := Response{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
