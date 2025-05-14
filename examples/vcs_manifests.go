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
		enbuild.WithDebug(true), // Enable debug mode
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

	fmt.Println("Fetching GitHub manifests...")
	// List GitHub manifests
	githubManifests, err := client.Manifests.List(&enbuild.ManifestListOptions{
		VCS: enbuild.VCSTypeGitHub,
	})
	if err != nil {
		log.Printf("Error listing GitHub manifests: %v", err)
	} else {
		fmt.Printf("Found %d GitHub manifests\n", len(githubManifests))
		for _, manifest := range githubManifests {
			fmt.Printf("- %s (%s)\n", manifest.Name, manifest.ID)
		}
	}

	// fmt.Println("\nFetching GitLab manifests...")
	// // List GitLab manifests
	// gitlabManifests, err := client.Manifests.List(&enbuild.ManifestListOptions{
	// 	VCS: enbuild.VCSTypeGitLab,
	// })
	// if err != nil {
	// 	log.Printf("Error listing GitLab manifests: %v", err)
	// } else {
	// 	fmt.Printf("Found %d GitLab manifests\n", len(gitlabManifests))
	// 	for _, manifest := range gitlabManifests {
	// 		fmt.Printf("- %s (%s)\n", manifest.Name, manifest.ID)
	// 	}
	// }

	// fmt.Println("\nFetching regular manifests...")
	// // List regular manifests (no VCS specified)
	// manifests, err := client.Manifests.List(nil)
	// if err != nil {
	// 	log.Printf("Error listing regular manifests: %v", err)
	// } else {
	// 	fmt.Printf("Found %d regular manifests\n", len(manifests))
	// 	for _, manifest := range manifests {
	// 		fmt.Printf("- %s (%s)\n", manifest.Name, manifest.ID)
	// 	}
	// }

	// // Get a specific GitHub manifest by ID if any exist
	// if len(githubManifests) > 0 {
	// 	id := githubManifests[0].ID
	// 	fmt.Printf("\nFetching GitHub manifest with ID: %s\n", id)
	// 	manifest, err := client.Manifests.Get(id, &enbuild.ManifestListOptions{
	// 		VCS: enbuild.VCSTypeGitHub,
	// 	})
	// 	if err != nil {
	// 		log.Printf("Error getting GitHub manifest: %v", err)
	// 	} else {
	// 		fmt.Printf("GitHub Manifest Details:\n")
	// 		fmt.Printf("- ID: %s\n", manifest.ID)
	// 		fmt.Printf("- Name: %s\n", manifest.Name)
	// 		fmt.Printf("- Description: %s\n", manifest.Description)
	// 		fmt.Printf("- Version: %s\n", manifest.Version)
	// 	}
	// }
}
