package main

import (
    "fmt"
    "log"
    "os"

    "github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

const debug = true

func printStacks(Stacks []*enbuild.Catalog) {
    fmt.Printf("Found %d Stacks\n", len(Stacks))
    for _, catalog := range Stacks {
        fmt.Printf("ID: %v, Name: %v, Type: %v, Slug: %v, VCS: %v\n",
            catalog.ID, catalog.Name, catalog.Type, catalog.Slug, catalog.VCS)
    }
}

func createClient() (*enbuild.Client, error) {
    username := os.Getenv("ENBUILD_USERNAME")
    password := os.Getenv("ENBUILD_PASSWORD")
    baseURL := os.Getenv("ENBUILD_BASE_URL")

    fmt.Printf("Using ENBUILD_BASE_URL: %s\n", baseURL)
    fmt.Printf("Using ENBUILD_USERNAME: %s\n", username)
    if password != "" {
        fmt.Printf("ENBUILD_PASSWORD is set\n")
    } else {
        fmt.Printf("ENBUILD_PASSWORD is not set\n")
		os.Exit(1)
    }

    options := []enbuild.ClientOption{
        enbuild.WithDebug(debug),
    }

    if baseURL != "" {
        options = append(options, enbuild.WithBaseURL(baseURL))
    }

    if username != "" && password != "" {
        options = append(options, enbuild.WithKeycloakAuth(username, password))
    }

    return enbuild.NewClient(options...)
}

func listAllStacks(client *enbuild.Client) {
    fmt.Println("Listing all Stacks:")
    allStacks, err := client.Stacks.List()
    if err != nil {
        log.Fatalf("Error listing Stacks: %v", err)
    }
    printStacks(allStacks)
}

func listGitHubStacks(client *enbuild.Client) {
    fmt.Println("\nListing GitHub Stacks:")
    githubStacks, err := client.Stacks.List(&enbuild.CatalogListOptions{
        VCS: "github",
    })
    if err != nil {
        log.Fatalf("Error listing GitHub Stacks: %v", err)
    }
    printStacks(githubStacks)
}

func listGitLabStacks(client *enbuild.Client) {
    fmt.Println("\nListing GitLab Stacks:")
    gitlabStacks, err := client.Stacks.List(&enbuild.CatalogListOptions{
        VCS: "gitlab",
    })
    if err != nil {
        log.Fatalf("Error listing GitLab Stacks: %v", err)
    }
    printStacks(gitlabStacks)
}

func getCatalogByID(client *enbuild.Client, id string) {
    fmt.Printf("\nGetting catalog with ID %s:\n", id)
    catalog, err := client.Stacks.Get(id, &enbuild.CatalogListOptions{})
    if err != nil {
        log.Fatalf("Error getting catalog: %v", err)
    }
    printStacks([]*enbuild.Catalog{catalog})
}

func filterStacksByType(client *enbuild.Client) {
    fmt.Println("\nFiltering Stacks by type 'terraform':")
    terraformStacks, err := client.Stacks.List(&enbuild.CatalogListOptions{
        Type: "terraform",
    })
    if err != nil {
        log.Fatalf("Error filtering Stacks: %v", err)
    }
    printStacks(terraformStacks)
}

func searchStacksByName(client *enbuild.Client) {
    searchTerm := "Bang"
    fmt.Printf("\nSearching Stacks with name containing '%s':\n", searchTerm)
    searchResults, err := client.Stacks.List(&enbuild.CatalogListOptions{
        Name: searchTerm,
    })
    if err != nil {
        log.Fatalf("Error searching Stacks: %v", err)
    }
    printStacks(searchResults)
}

func main() {
    client, err := createClient()
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }

    getCatalogByID(client, "6638a128d6852d0012a27491")
	// listGitLabStacks(client)
}
