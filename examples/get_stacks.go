package main

import (
	"context"
    "fmt"
    "log"
    "os"

    "github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

const debug = false

func printStacks(Stacks []*enbuild.Stack) {
    for _, Stack := range Stacks {
        fmt.Printf("ID: %v Name: %v\n", Stack.ID, Stack.Name)
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

    return enbuild.NewClient(context.Background(), options...)
}

func listAllStacks(client *enbuild.Client) {
    fmt.Println("Listing all Stacks (page 0, limit 10):")
    allStacks, err := client.Stacks.ListStacks(context.Background(), 0, 10, "")
    if err != nil {
        log.Fatalf("Error listing Stacks: %v", err)
    }

    fmt.Printf("Found: %d stacks\n", len(allStacks))
    printStacks(allStacks)

    // Demonstrate parameter usage
    fmt.Println("\nListing stacks with searchTerm 'test' (page 0, limit 5):")
    searchedStacks, err := client.Stacks.ListStacks(context.Background(), 0, 5, "test")
    if err != nil {
        log.Fatalf("Error listing Stacks with search: %v", err)
    }
    fmt.Printf("Found: %d stacks with search term 'test'\n", len(searchedStacks))
    printStacks(searchedStacks)
}

func main() {
    client, err := createClient()
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }
    listAllStacks(client)
}
