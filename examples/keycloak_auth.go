package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

func main() {
	// Get username and password from environment variables or command line
	username := os.Getenv("ENBUILD_USERNAME")
	password := os.Getenv("ENBUILD_PASSWORD")
	
	// If not provided in environment variables, check command line arguments
	if username == "" && len(os.Args) > 1 {
		username = os.Args[1]
	}
	
	if password == "" && len(os.Args) > 2 {
		password = os.Args[2]
	}
	
	// Validate credentials
	if username == "" || password == "" {
		log.Fatalf("Usage: %s <username> <password>\nOr set ENBUILD_USERNAME and ENBUILD_PASSWORD environment variables", os.Args[0])
	}
	
	fmt.Printf("Authenticating with username: %s\n", username)
	
	// Create a new client with Keycloak authentication
	client, err := enbuild.NewClient(
		enbuild.WithKeycloakAuth(username, password),
		enbuild.WithDebug(true),
	)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	
	fmt.Println("Authentication successful!")
	fmt.Println("Listing all Catalogs:")
	
	// List catalogs
	catalogs, err := client.Catalogs.List(nil)
	if err != nil {
		log.Fatalf("Error listing catalogs: %v", err)
	}
	
	fmt.Printf("Found %d catalogs\n", len(catalogs))
	
	// Print catalog details
	for _, catalog := range catalogs {
		fmt.Printf("ID: %s, Name: %s, Type: %s, Slug: %s, VCS: %s\n", 
			catalog.ID, 
			catalog.Name, 
			catalog.Type, 
			catalog.Slug,
			catalog.VCS,
		)
	}
}
