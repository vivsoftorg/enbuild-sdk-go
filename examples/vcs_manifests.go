package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/manifests"
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
		enbuild.WithDebug(false), // Disable debug mode
	}

	// Get base URL from environment variable if provided
	if baseURL := os.Getenv("ENBUILD_BASE_URL"); baseURL != "" {
		options = append(options, enbuild.WithBaseURL(baseURL))
	}

	// Create a new client
	client, err := enbuild.NewClient(options...)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Fetch and print GitHub manifests
	fmt.Println("Fetching GitHub manifests...")
	fetchAndPrintManifests(client, "GITHUB")

	// Fetch and print GitLab manifests
	fmt.Println("\nFetching GitLab manifests...")
	fetchAndPrintManifests(client, "GITLAB")
}

func fetchAndPrintManifests(client *enbuild.Client, vcsType string) {
	manifests, err := client.Manifests.List(&manifests.ManifestListOptions{
		VCS: vcsType,
	})
	if err != nil {
		log.Printf("Error listing %s manifests: %v", vcsType, err)
		return // Return early on error
	}

	fmt.Printf("Found %d %s manifests\n", len(manifests), vcsType)
	for _, manifest := range manifests {
		fmt.Printf("ID: %v, Name: %v, Type: %v, Slug: %v\n",
			manifest.ID, manifest.Name, manifest.Type, manifest.Slug)
	}
}
