# ENBUILD SDK for Go

A Go client library for accessing the ENBUILD API.

## Version

Current version: 0.0.2

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
    enbuild.WithKeycloakAuth("your-username", "your-password"),
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

The ENBUILD SDK uses Keycloak authentication by default. You can provide your credentials when creating the client:

```go
client, err := enbuild.NewClient(
    enbuild.WithKeycloakAuth("your-username", "your-password"),
)
```

If no credentials are provided, the SDK will:
1. Look for the `ENBUILD_USERNAME` and `ENBUILD_PASSWORD` environment variables
2. Fall back to default credentials if the environment variables are not set

You can set the username and password using environment variables:
```bash
export ENBUILD_USERNAME="your-username"
export ENBUILD_PASSWORD="your-password"
```

The SDK will automatically:
1. Fetch the Keycloak configuration from the ENBUILD admin settings
2. Authenticate with Keycloak to get an access token
3. Automatically refresh the token when it expires

## Configuration

The SDK can be configured using environment variables:

- `ENBUILD_USERNAME`: Username for Keycloak authentication
- `ENBUILD_PASSWORD`: Password for Keycloak authentication
- `ENBUILD_BASE_URL`: Base URL for the API (optional, defaults to the production API endpoint)
  - Note: The SDK will automatically append `/api/v1/` to the base URL if it's not already included

Example:
```bash
export ENBUILD_USERNAME="your-username"
export ENBUILD_PASSWORD="your-password"
export ENBUILD_BASE_URL="https://enbuild-dev.vivplatform.io/enbuild-bk"
```

You can also configure the client programmatically:

```go
// Create client options
options := []enbuild.ClientOption{
    enbuild.WithDebug(true), // Enable debug mode
    enbuild.WithKeycloakAuth("your-username", "your-password"),
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
- `keycloak_auth.go`: Shows how to authenticate with Keycloak using username and password

## Configuration Options

The client can be configured with the following options:

- `WithBaseURL`: Set a custom base URL for the API
- `WithTimeout`: Set a custom timeout for API requests
- `WithKeycloakAuth`: Set the Keycloak authentication credentials
- `WithDebug`: Enable or disable debug output

Example:

```go
client, err := enbuild.NewClient(
    enbuild.WithBaseURL("https://custom-api.enbuild.com"),
    enbuild.WithTimeout(60 * time.Second),
    enbuild.WithKeycloakAuth("your-username", "your-password"),
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

- Implemented Keycloak authentication with automatic token refresh
- Fixed URL construction to properly handle base URLs with or without trailing slashes
- Enhanced debug output with masked sensitive information
- Improved error handling and reporting
- Updated examples to demonstrate all key features
