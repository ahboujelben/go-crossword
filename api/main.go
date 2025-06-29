package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ahboujelben/go-crossword/api/server"
)

func main() {
	// Read configuration from environment variables with defaults
	port := getEnvInt("PORT", 8080)
	ollamaURL := getEnvString("OLLAMA_URL", "http://localhost:11434")
	ollamaModel := getEnvString("OLLAMA_MODEL", "llama3:8b")

	// Start the API server
	fmt.Printf("Starting API server on port %d...\n", port)
	apiServer := server.NewServer(server.ServerConfig{
		Port:        port,
		OllamaURL:   ollamaURL,
		OllamaModel: ollamaModel,
	})

	fmt.Printf("API server running on http://localhost:%d/api/crossword\n", port)
	fmt.Printf("Using Ollama model: %s at %s\n", ollamaModel, ollamaURL)
	fmt.Println("Press Ctrl+C to stop the server")

	// Start the server and block until it exits
	if err := apiServer.Start(); err != nil {
		log.Fatalf("Error starting API server: %v", err)
	}
}

// getEnvString gets a string value from an environment variable with a fallback
func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvInt gets an integer value from an environment variable with a fallback
func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
