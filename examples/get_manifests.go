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
	// Get API token from environment variable
	token := os.Getenv("ENBUILD_API_TOKEN")
	if token == "" {
		log.Fatal("ENBUILD_API_TOKEN environment variable is required")
	}

	// Create client options
	options := []enbuild.ClientOption{
		enbuild.WithAuthToken(token),
		enbuild.WithDebug(true), // Enable debug mode
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

	// // Example 1: List all manifests
	// fmt.Println("Listing all manifests:")
	// allManifests, err := client.Manifests.List(nil)
	// if err != nil {
	// 	log.Fatalf("Error listing manifests: %v", err)
	// }
	// printManifests(allManifests)

	// // Example 2: List GitHub manifests
	// fmt.Println("\nListing GitHub manifests:")
	// githubManifests, err := client.Manifests.List(&manifests.ManifestListOptions{
	// 	VCS: "github",
	// })
	// if err != nil {
	// 	log.Fatalf("Error listing GitHub manifests: %v", err)
	// }
	// printManifests(githubManifests)

	// // Example 3: List GitLab manifests
	// fmt.Println("\nListing GitLab manifests:")
	// gitlabManifests, err := client.Manifests.List(&manifests.ManifestListOptions{
	// 	VCS: "gitlab",
	// })
	// if err != nil {
	// 	log.Fatalf("Error listing GitLab manifests: %v", err)
	// }
	// printManifests(gitlabManifests)

	// Example 4: Get manifest by ID
	id := "6638a128d6852d0012a27491"
	fmt.Printf("\nGetting manifest with ID %s:\n", id)
	manifest, err := client.Manifests.Get(id, &manifests.ManifestListOptions{})
	if err != nil {
		log.Fatalf("Error getting manifest: %v", err)
	}
	printManifests([]*types.Manifest{manifest})

	// // Example 5: Filter manifests by type
	// fmt.Println("\nFiltering manifests by type 'terraform':")
	// terraformManifests, err := client.Manifests.List(&manifests.ManifestListOptions{
	// 	Type: "terraform",
	// })
	// if err != nil {
	// 	log.Fatalf("Error filtering manifests: %v", err)
	// }
	// printManifests(terraformManifests)

	// // Example 6: Search manifests by name
	// searchTerm := "Bang"
	// fmt.Printf("\nSearching manifests with name containing '%s':\n", searchTerm)
	// searchResults, err := client.Manifests.List(&manifests.ManifestListOptions{
	// 	Name: searchTerm,
	// })
	// if err != nil {
	// 	log.Fatalf("Error searching manifests: %v", err)
	// }
	// printManifests(searchResults)
}
