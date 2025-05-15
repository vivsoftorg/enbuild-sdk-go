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

// all_manifest_ListOptions := &manifests.ManifestListOptions{}
// github_manifest_ListOptions := &manifests.ManifestListOptions{VCS: "GITHUB"}
// gitlab_manifest_ListOptions := &manifests.ManifestListOptions{VCS: "GITLAB"}
get_by_id_ListOptions := &manifests.ManifestListOptions{ID: "6638a0a4d6852d0012a27490"}

	// List all manifests
	// fmt.Println("Listing all manifests:")
	// fetchAndPrintManifests(client, all_manifest_ListOptions)

	// Fetch and print GitHub manifests
	// fmt.Println("Fetching GitHub manifests...")
	// fetchAndPrintManifests(client, github_manifest_ListOptions)

	// Fetch and print GitLab manifests
	// fmt.Println("\nFetching GitLab manifests...")
	// fetchAndPrintManifests(client, gitlab_manifest_ListOptions)

	// Fetch and print manifest by ID
	fmt.Println("\nFetching manifest by ID...")
	fetchAndPrintManifests(client, get_by_id_ListOptions)
}

func fetchAndPrintManifests(client *enbuild.Client, options *manifests.ManifestListOptions) {
	manifests, err := client.Manifests.List(options)
	if err != nil {
		log.Printf("Error listing manifests: %v", err)
		return // Return early on error
	}

	fmt.Printf("Found %d manifests\n", len(manifests))
	for _, manifest := range manifests {
		fmt.Printf("ID: %v, Name: %v, Type: %v, Slug: %v, VCS: %v\n",
			manifest.ID, manifest.Name, manifest.Type, manifest.Slug, manifest.VCS)
	}
}
