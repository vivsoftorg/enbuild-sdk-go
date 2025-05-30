# ENBUILD SDK for Go

A Go client library for accessing the ENBUILD API.

## Version

Current version: 0.0.3

## Installation

```bash
go get github.com/vivsoftorg/enbuild-sdk-go
```

## Usage

```go
import (
    "context"
    "log"
    "os"

    "github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

func main() {
    // Read credentials from environment variables
    username := os.Getenv("ENBUILD_USERNAME")
    password := os.Getenv("ENBUILD_PASSWORD")
    baseURL := os.Getenv("ENBUILD_BASE_URL")

    options := []enbuild.ClientOption{
        enbuild.WithDebug(true), // Enable debug output
    }
    if baseURL != "" {
        options = append(options, enbuild.WithBaseURL(baseURL))
    }
    if username != "" && password != "" {
        options = append(options, enbuild.WithKeycloakAuth(username, password))
    }

    // Create a new client with context
    client, err := enbuild.NewClient(context.Background(), options...)
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }

    // List all catalogs
    catalogs, err := client.Catalogs.ListCatalog(context.Background())
    if err != nil {
        log.Fatalf("Error listing catalogs: %v", err)
    }
    log.Printf("Found %d catalogs\n", len(catalogs))
}
```

## Examples

See the [examples](./examples) directory for more usage patterns:

- **get_catalogs.go**:  
  Demonstrates listing all catalogs, filtering by VCS (`github`, `gitlab`), filtering by type, searching by name, and getting a catalog by ID.
- **get_stacks.go**:  
  Shows how to list all stacks with pagination and search term.

**Example: List all stacks**

```go
client, err := enbuild.NewClient(context.Background(), options...)
if err != nil {
    log.Fatalf("Error creating client: %v", err)
}
page := 0
limit := 10
searchTerm := ""
stacks, err := client.Stacks.ListStacks(context.Background(), page, limit, searchTerm)
if err != nil {
    log.Fatalf("Error listing stacks: %v", err)
}
for _, stack := range stacks {
    log.Printf("ID: %v Name: %v Type: %v Status: %v\n", stack.ID, stack.Name, stack.Type, stack.Status)
}
```

**Example: List GitHub catalogs**

```go
githubCatalogs, err := client.Catalogs.ListCatalog(context.Background(), &enbuild.CatalogListOptions{
    VCS: "github",
})
if err != nil {
    log.Fatalf("Error listing GitHub catalogs: %v", err)
}
```

**Example: Get a catalog by ID**

```go
catalog, err := client.Catalogs.GetCatalog(context.Background(), "catalog-id", &enbuild.CatalogListOptions{})
if err != nil {
    log.Fatalf("Error getting catalog: %v", err)
}
```