# ENBUILD SDK for Go

A Go client library for accessing the ENBUILD API.

## Version

Current version: 0.0.4

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