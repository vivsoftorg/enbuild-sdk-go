# ENBUILD SDK for Go

# THIS IS A WORK IN PROGRESS FOR NOW JUST A POC

A Go client library for accessing the ENBUILD API.

## Installation

```bash
go get github.com/vivsoftorg/enbuild-sdk-go
```

## Usage

```go
import "github.com/vivsoftorg/enbuild-sdk-go"
```

Create a new ENBUILD client, then use the services on the client to access different parts of the ENBUILD API. For example:

```go
// Create a new client with an authentication token
client, err := enbuild.NewClient(
    enbuild.WithAuthToken("your-api-token"),
)
if err != nil {
    log.Fatalf("Error creating client: %v", err)
}

// List users
users, err := client.Users.List(nil)
if err != nil {
    log.Fatalf("Error listing users: %v", err)
}
fmt.Printf("Found %d users\n", len(users))
```

## Authentication

The ENBUILD API uses bearer token authentication. You can set the authentication token when creating the client:

```go
client, err := enbuild.NewClient(
    enbuild.WithAuthToken("your-api-token"),
)
```

## Configuration

The SDK can be configured using environment variables:

- `ENBUILD_API_TOKEN`: Authentication token for the API (required)
- `ENBUILD_BASE_URL`: Base URL for the API (optional, defaults to the production API endpoint)
  - Note: The SDK will automatically append `/api/v1/` to the base URL if it's not already included

Example:
```bash
export ENBUILD_API_TOKEN="your-api-token"
export ENBUILD_BASE_URL="https://enbuild-dev.vivplatform.io/enbuild-bk"
```

You can also configure the client programmatically:

```go
// Get configuration from environment variables
token := os.Getenv("ENBUILD_API_TOKEN")
baseURL := os.Getenv("ENBUILD_BASE_URL")

// Create client options
options := []enbuild.ClientOption{
    enbuild.WithAuthToken(token),
}

// Add base URL if provided
if baseURL != "" {
    options = append(options, enbuild.WithBaseURL(baseURL))
}

// Create a new client
client, err := enbuild.NewClient(options...)
```

## Services

The ENBUILD SDK provides the following services:

- `Users`: Manage users
- `Roles`: Manage roles and permissions
- `Operations`: Manage operations
- `Repository`: Access repositories
- `Manifests`: Access manifests
- `MLDataset`: Access ML datasets
- `AdminSettings`: Access admin settings

## Examples

See the [examples](./examples) directory for more examples of using the SDK.

## Configuration Options

The client can be configured with the following options:

- `WithBaseURL`: Set a custom base URL for the API
- `WithTimeout`: Set a custom timeout for API requests
- `WithAuthToken`: Set the authentication token

Example:

```go
client, err := enbuild.NewClient(
    enbuild.WithBaseURL("https://custom-api.enbuild.com/api/v1/"),
    enbuild.WithTimeout(60 * time.Second),
    enbuild.WithAuthToken("your-api-token"),
)
```

## Error Handling

All API methods return an error as the last return value. If the API returns an error response, the error will contain the error message from the API.

```go
users, err := client.Users.List(nil)
if err != nil {
    // Handle error
    log.Fatalf("Error listing users: %v", err)
}
```

## License

This SDK is distributed under the MIT license.

## Notes on API Response Handling

- Timestamp fields like `createdOn` and `updatedOn` are defined as `interface{}` in the SDK models because the API may return them as either numeric timestamps or formatted strings.
- When working with these fields, you may need to check the type and convert accordingly:

```go
// Example of handling a timestamp field
if timestamp, ok := user.CreatedOn.(float64); ok {
    // Handle numeric timestamp
    createdTime := time.Unix(int64(timestamp/1000), 0)
    fmt.Printf("Created on: %s\n", createdTime.Format(time.RFC3339))
} else if timestampStr, ok := user.CreatedOn.(string); ok {
    // Handle string timestamp
    fmt.Printf("Created on: %s\n", timestampStr)
}
```
### VCS-Specific Manifests

The SDK supports accessing VCS-specific manifests (GitHub or GitLab):

```go
// List GitHub manifests
githubManifests, err := client.Manifests.List(&enbuild.ManifestListOptions{
    VCS: enbuild.VCSTypeGitHub,
})

// List GitLab manifests
gitlabManifests, err := client.Manifests.List(&enbuild.ManifestListOptions{
    VCS: enbuild.VCSTypeGitLab,
})

// Get a specific GitHub manifest by ID
manifest, err := client.Manifests.Get(id, &enbuild.ManifestListOptions{
    VCS: enbuild.VCSTypeGitHub,
})
```

Note: The GitHub and GitLab manifest endpoints return data in a different structure than the standard endpoints. The SDK handles this difference automatically by parsing the `catalogManifest` array in the response.
