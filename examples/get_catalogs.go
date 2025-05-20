package main

import (
	"fmt"
	"log"
	"os"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

func printCatalogs(catalogs []*enbuild.Catalog) {
	fmt.Printf("Found %d catalogs\n", len(catalogs))
	for _, catalog := range catalogs {
		fmt.Printf("ID: %v, Name: %v, Type: %v, Slug: %v, VCS: %v\n",
			catalog.ID, catalog.Name, catalog.Type, catalog.Slug, catalog.VCS)
	}
}


func main() {

	// Get credentials from environment variables
	username := os.Getenv("ENBUILD_USERNAME")
	password := os.Getenv("ENBUILD_PASSWORD")
	baseURL := os.Getenv("ENBUILD_BASE_URL")

	// Print environment variables for debugging
	fmt.Printf("Using ENBUILD_BASE_URL: %s\n", baseURL)
	fmt.Printf("Using ENBUILD_USERNAME: %s\n", username)
	if password != "" {
		fmt.Printf("ENBUILD_PASSWORD is set\n")
	} else {
		fmt.Printf("ENBUILD_PASSWORD is not set\n")
	}

	// Create client options
	options := []enbuild.ClientOption{
		enbuild.WithDebug(true), // Enable debug mode
	}

	// Add base URL if provided
	if baseURL != "" {
		options = append(options, enbuild.WithBaseURL(baseURL))
	}

	// Add authentication if credentials are provided
	if username != "" && password != "" {
		options = append(options, enbuild.WithKeycloakAuth(username, password))
	}

	// Create a new client
	client, err := enbuild.NewClient(options...)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// -----------------------------------------------------------------------------------------------------------------
	// Example 1: List all Catalogs
	fmt.Println("Listing all Catalogs:")
	allCatalogs, err := client.Catalogs.List()
	if err != nil {
		log.Fatalf("Error listing Catalogs: %v", err)
	}
	printCatalogs(allCatalogs)
	// -----------------------------------------------------------------------------------------------------------------
	// Example 2: List GitHub Catalogs
	fmt.Println("\nListing GitHub Catalogs:")
	githubCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
		VCS: "github",
	})
	if err != nil {
		log.Fatalf("Error listing GitHub Catalogs: %v", err)
	}
	printCatalogs(githubCatalogs)
	// -----------------------------------------------------------------------------------------------------------------
	// Example 3: List GitLab Catalogs
	fmt.Println("\nListing GitLab Catalogs:")
	gitlabCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
		VCS: "gitlab",
	})
	if err != nil {
		log.Fatalf("Error listing GitLab Catalogs: %v", err)
	}
	printCatalogs(gitlabCatalogs)
	// -----------------------------------------------------------------------------------------------------------------
	// Example 4: Get catalog by ID
	id := "6638a128d6852d0012a27491"
	fmt.Printf("\nGetting catalog with ID %s:\n", id)
	catalog, err := client.Catalogs.Get(id, &enbuild.CatalogListOptions{})
	if err != nil {
		log.Fatalf("Error getting catalog: %v", err)
	}
	printCatalogs([]*enbuild.Catalog{catalog})
	// -----------------------------------------------------------------------------------------------------------------
	// Example 5: Filter Catalogs by type
	fmt.Println("\nFiltering Catalogs by type 'terraform':")
	terraformCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
		Type: "terraform",
	})
	if err != nil {
		log.Fatalf("Error filtering Catalogs: %v", err)
	}
	printCatalogs(terraformCatalogs)
	// -----------------------------------------------------------------------------------------------------------------
	// Example 6: Search Catalogs by name
	searchTerm := "Bang"
	fmt.Printf("\nSearching Catalogs with name containing '%s':\n", searchTerm)
	searchResults, err := client.Catalogs.List(&enbuild.CatalogListOptions{
		Name: searchTerm,
	})
	if err != nil {
		log.Fatalf("Error searching Catalogs: %v", err)
	}
	printCatalogs(searchResults)
	// -----------------------------------------------------------------------------------------------------------------
}
