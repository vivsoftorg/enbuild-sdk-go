package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/manifests"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/types"
)

func printManifests(manifests []*types.Manifest) {
	fmt.Printf("Found %d manifests\n", len(manifests))
	for _, m := range manifests {
		fmt.Printf("ID: %v, Name: %v, Type: %v, Slug: %v, VCS: %v\n",
			m.ID, m.Name, m.Type, m.Slug, m.VCS)
	}
}

func main() {
	// Create client options
	options := []enbuild.ClientOption{
		enbuild.WithDebug(false), // Enable debug mode
	}

	// Get API token from environment variable if provided
	if token := os.Getenv("ENBUILD_API_TOKEN"); token != "" {
		options = append(options, enbuild.WithAuthToken(token))
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

	// Example: Get manifest by ID
	id := "6638a128d6852d0012a27491"
	fmt.Printf("\nGetting manifest with ID %s:\n", id)
	manifest, err := client.Manifests.Get(id, &manifests.ManifestListOptions{})
	if err != nil {
		log.Fatalf("Error getting manifest: %v", err)
	}
	printManifests([]*types.Manifest{manifest})
}
