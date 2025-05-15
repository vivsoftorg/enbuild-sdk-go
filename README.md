# ENBUILD SDK for Go

A Go client library for accessing the ENBUILD API.

## Version

Current version: 0.0.1 (Initial Release)

## Installation

```bash
go get github.com/vivsoftorg/enbuild-sdk-go
```

## Usage

```go
import "github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
```

Create a new ENBUILD client, then use the services on the client to access different parts of the ENBUILD API. For example:

```go
// Create a new client
client, err := enbuild.NewClient(
    enbuild.WithAuthToken("your-api-token"),
    enbuild.WithDebug(true), // Enable debug output
)
if err != nil {
    log.Fatalf("Error creating client: %v", err)
}

// List catalogs
catalogs, err := client.Catalogs.List()
if err != nil {
    log.Fatalf("Error listing catalogs: %v", err)
}
fmt.Printf("Found %d catalogs\n", len(catalogs))
```

## Authentication

The ENBUILD API uses bearer token authentication. You can set the authentication token when creating the client:

```go
client, err := enbuild.NewClient(
    enbuild.WithAuthToken("your-api-token"),
)
```

If no token is provided, the SDK will:
1. Look for the `ENBUILD_API_TOKEN` environment variable
2. Fall back to a default token if the environment variable is not set

## Configuration

The SDK can be configured using environment variables:

- `ENBUILD_API_TOKEN`: Authentication token for the API (optional, falls back to default token)
- `ENBUILD_BASE_URL`: Base URL for the API (optional, defaults to the production API endpoint)
  - Note: The SDK will automatically append `/api/v1/` to the base URL if it's not already included

Example:
```bash
export ENBUILD_API_TOKEN="your-api-token"
export ENBUILD_BASE_URL="https://enbuild-dev.vivplatform.io/enbuild-bk"
```

You can also configure the client programmatically:

```go
// Create client options
options := []enbuild.ClientOption{
    enbuild.WithDebug(true), // Enable debug mode
}

// Get API token from environment variable if provided
if token := os.Getenv("ENBUILD_API_TOKEN"); token != "" {
    options = append(options, enbuild.WithAuthToken(token))
}

// Get base URL from environment variable if provided
if baseURL := os.Getenv("ENBUILD_BASE_URL"); baseURL != "" {
    options = append(options, enbuild.WithBaseURL(baseURL))
}

// Create a new client
client, err := enbuild.NewClient(options...)
```

## Services

The ENBUILD SDK currently provides the following services:

- `Catalogs`: Access and manage catalogs/manifests

Additional services will be added in future releases:
- Users
- Roles
- Operations
- Repository
- MLDataset
- AdminSettings

## Examples

See the [examples](./examples) directory for more examples of using the SDK:

- `get_catalogs.go`: Demonstrates how to list, filter, and retrieve catalogs
- `get_manifests.go`: Alias for get_catalogs.go (for backward compatibility)

## Configuration Options

The client can be configured with the following options:

- `WithBaseURL`: Set a custom base URL for the API
- `WithTimeout`: Set a custom timeout for API requests
- `WithAuthToken`: Set the authentication token
- `WithDebug`: Enable or disable debug output

Example:

```go
client, err := enbuild.NewClient(
    enbuild.WithBaseURL("https://custom-api.enbuild.com"),
    enbuild.WithTimeout(60 * time.Second),
    enbuild.WithAuthToken("your-api-token"),
    enbuild.WithDebug(true),
)
```

## Debug Mode

When debug mode is enabled, the SDK will output detailed information about:
- Base URL being used
- Authentication token (masked for security)
- Request headers
- Request and response details

This is useful for troubleshooting API interactions.

## Error Handling

All API methods return an error as the last return value. If the API returns an error response, the error will contain the error message from the API.

```go
catalogs, err := client.Catalogs.List()
if err != nil {
    // Handle error
    log.Fatalf("Error listing catalogs: %v", err)
}
```

## License

This SDK is distributed under the MIT license.

## Notes on API Response Handling

- Timestamp fields like `createdOn` and `updatedOn` are defined as `interface{}` in the SDK models because the API may return them as either numeric timestamps or formatted strings.
- When working with these fields, you may need to check the type and convert accordingly:

```go
// Example of handling a timestamp field
if timestamp, ok := catalog.CreatedOn.(float64); ok {
    // Handle numeric timestamp
    createdTime := time.Unix(int64(timestamp/1000), 0)
    fmt.Printf("Created on: %s\n", createdTime.Format(time.RFC3339))
} else if timestampStr, ok := catalog.CreatedOn.(string); ok {
    // Handle string timestamp
    fmt.Printf("Created on: %s\n", timestampStr)
}
```

## VCS-Specific Catalogs

The SDK supports accessing VCS-specific catalogs (GitHub or GitLab):

```go
// List GitHub catalogs
githubCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
    VCS: "github",
})

// List GitLab catalogs
gitlabCatalogs, err := client.Catalogs.List(&enbuild.CatalogListOptions{
    VCS: "gitlab",
})

// Get a specific catalog by ID
catalog, err := client.Catalogs.Get(id, &enbuild.CatalogListOptions{})
```

## Recent Improvements

- Fixed URL construction to properly handle base URLs with or without trailing slashes
- Added default token support for easier development and testing
- Enhanced debug output with masked sensitive information
- Improved error handling and reporting
- Updated examples to demonstrate all key features
