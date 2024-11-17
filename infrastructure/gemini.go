package infrastructure

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func GeminiGo() (*genai.Client, context.Context) {
	ctx := context.Background()

	// Load API key from .env file (ensure correct path)
	dir, _ := os.Getwd()
	currDir := fmt.Sprintf("%s/.env", dir)
	err := godotenv.Load(currDir)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Access your API key from environment variable
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Missing API_KEY environment variable. Set it up as described in the instructions.")
	}

	// Create client with API key
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal("Error creating client:", err)
	}

	return client, ctx
}
