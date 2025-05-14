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

	// List users
	users, err := client.Users.List(nil)
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}
	fmt.Printf("Found %d users\n", len(users))
	for _, user := range users {
		fmt.Printf("- %s (%s)\n", user.Username, user.Email)
	}

	// List operations
	operations, err := client.Operations.List(&enbuild.OperationListOptions{
		Limit: 5,
		Sort:  "-createdOn",
	})
	if err != nil {
		log.Fatalf("Error listing operations: %v", err)
	}
	fmt.Printf("\nFound %d operations\n", len(operations))
	for _, op := range operations {
		fmt.Printf("- %s (%s): %s\n", op.Name, op.ID, op.Status)
	}

	// Get admin settings
	settings, err := client.AdminSettings.Get()
	if err != nil {
		log.Fatalf("Error getting admin settings: %v", err)
	}
	fmt.Printf("\nAdmin Settings: %+v\n", settings)
}
