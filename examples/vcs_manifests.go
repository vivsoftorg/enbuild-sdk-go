package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

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
		enbuild.WithDebug(false), // Disable debug mode
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
		for i, manifest := range githubManifests {
			idStr := ""
			switch id := manifest.ID.(type) {
			case string:
				idStr = id
			case float64:
				idStr = fmt.Sprintf("%.0f", id)
			case int:
				idStr = fmt.Sprintf("%d", id)
			default:
				idStr = fmt.Sprintf("%v (type: %s)", manifest.ID, reflect.TypeOf(manifest.ID))
			}
			
			fmt.Printf("%d. ID: %s, Name: %s, Type: %s, Slug: %s\n", 
				i+1, idStr, manifest.Name, manifest.Type, manifest.Slug)
		}
	}
}
