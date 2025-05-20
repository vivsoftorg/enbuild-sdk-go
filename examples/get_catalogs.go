package main

import (
    "fmt"
    "log"
    "os"

    "github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

const debug = false

func printCatalogs(catalogs []*enbuild.Catalog) {
    fmt.Printf("Found %d catalogs\n", len(catalogs))
    for _, catalog := range catalogs {
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

func listAllCatalogs(client *enbuild.Client) {
    fmt.Println("Listing all Catalogs:")
    allCatalogs, err := client.Catalogs.List()
    if err != nil {
        log.Fatalf("Error listing Catalogs: %v", err)
    }
    printCatalogs(allCatalogs)
}

func listGitHubCatalogs(client *enbuild.Client) {
    fmt.Println("\nListing GitHub Catalogs:")
    githubCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
        VCS: "github",
    })
    if err != nil {
        log.Fatalf("Error listing GitHub Catalogs: %v", err)
    }
    printCatalogs(githubCatalogs)
}

func listGitLabCatalogs(client *enbuild.Client) {
    fmt.Println("\nListing GitLab Catalogs:")
    gitlabCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
        VCS: "gitlab",
    })
    if err != nil {
        log.Fatalf("Error listing GitLab Catalogs: %v", err)
    }
    printCatalogs(gitlabCatalogs)
}

func getCatalogByID(client *enbuild.Client, id string) {
    fmt.Printf("\nGetting catalog with ID %s:\n", id)
    catalog, err := client.Catalogs.Get(id, &enbuild.CatalogListOptions{})
    if err != nil {
        log.Fatalf("Error getting catalog: %v", err)
    }
    printCatalogs([]*enbuild.Catalog{catalog})
}

func filterCatalogsByType(client *enbuild.Client) {
    fmt.Println("\nFiltering Catalogs by type 'terraform':")
    terraformCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
        Type: "terraform",
    })
    if err != nil {
        log.Fatalf("Error filtering Catalogs: %v", err)
    }
    printCatalogs(terraformCatalogs)
}

func searchCatalogsByName(client *enbuild.Client) {
    searchTerm := "Bang"
    fmt.Printf("\nSearching Catalogs with name containing '%s':\n", searchTerm)
    searchResults, err := client.Catalogs.List(&enbuild.CatalogListOptions{
        Name: searchTerm,
    })
    if err != nil {
        log.Fatalf("Error searching Catalogs: %v", err)
    }
    printCatalogs(searchResults)
}

func main() {
    client, err := createClient()
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }

    getCatalogByID(client, "6638a128d6852d0012a27491")
	// listGitLabCatalogs(client)
}
