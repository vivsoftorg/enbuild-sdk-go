package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vivsoftorg/enbuild-sdk-go"
)

func main() {
	// Get API token from environment variable
	token := os.Getenv("ENBUILD_API_TOKEN")
	if token == "" {
		log.Fatal("ENBUILD_API_TOKEN environment variable is required")
	}

	// Create client options
	options := []enbuild.ClientOption{
		enbuild.WithAuthToken(token),
	}

	// Get base URL from environment variable if provided
	baseURL := os.Getenv("ENBUILD_BASE_URL")
	if baseURL != "" {
		options = append(options, enbuild.WithBaseURL(baseURL))
	}

	// Create a new client
	client, err := enbuild.NewClient(options...)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// List manifests
	manifests, err := client.Manifests.List()
	if err != nil {
		log.Fatalf("Error listing manifests: %v", err)
	}
	fmt.Printf("Found %d manifests\n", len(manifests))
	for _, manifest := range manifests {
		fmt.Printf("- %s (%s)\n", manifest.Name, manifest.ID)
	}
}
