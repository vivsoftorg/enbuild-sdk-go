package main

import (
    "fmt"
    "log"
    "os"

    "github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
    "github.com/vivsoftorg/enbuild-sdk-go/pkg/manifests"
    "github.com/vivsoftorg/enbuild-sdk-go/pkg/types"
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

    fmt.Println("Fetching GitHub manifests...")
    // List GitHub manifests
    githubManifests, err := client.Manifests.List(&manifests.ManifestListOptions{
        VCS: types.VCSTypeGitHub,
    })
    if err != nil {
        log.Printf("Error listing GitHub manifests: %v", err)
        return  // Return early on error
    }

    fmt.Printf("Found %d GitHub manifests\n", len(githubManifests))
    for _, manifest := range githubManifests {
        fmt.Printf("ID: %v, Name: %v, Type: %v, Slug: %v\n",
           manifest.ID, manifest.Name, manifest.Type, manifest.Slug)
    }
}
