# ENBUILD SDK for Go

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

Example:
```bash
export ENBUILD_API_TOKEN="your-api-token"
export ENBUILD_BASE_URL="https://api.staging.enbuild.com/api/v1/"
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
