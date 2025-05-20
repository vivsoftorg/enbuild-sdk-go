package enbuild

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
)

const (
	defaultBaseURL = "https://enbuild.vivplatform.io"
	defaultTimeout = 30 * time.Second
	apiVersionPath = "/enbuild-bk/api/v1/"
	adminSettingsPath = "/enbuild-user/api/v1/adminSettings"
)

// Client represents the ENBUILD API client
type Client struct {
	httpClient  *request.Client
	authManager *AuthManager

	// Services
	Catalogs *Service
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client) error

// NewClient creates a new ENBUILD API client
func NewClient(options ...ClientOption) (*Client, error) {
	// Get base URL from environment variable if provided
	baseURLEnv := os.Getenv("ENBUILD_BASE_URL")
	baseURLToUse := defaultBaseURL
	if baseURLEnv != "" {
		baseURLToUse = baseURLEnv
	}

	// Process the base URL to ensure it has the API version path
	if !strings.HasSuffix(baseURLToUse, strings.TrimPrefix(apiVersionPath, "/")) &&
		!strings.Contains(baseURLToUse, apiVersionPath) {
		baseURLToUse += apiVersionPath
	}

	baseURL, _ := url.Parse(baseURLToUse)
	httpClient := &request.Client{
		BaseURL:    baseURL,
		UserAgent:  "enbuild-sdk-go",
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		Debug:      false,
	}

	c := &Client{
		httpClient: httpClient,
	}

	// Apply options
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	// If no token provider was set and no auth manager exists, use default credentials
	if c.httpClient.TokenProvider == nil && c.authManager == nil {
		username := os.Getenv("ENBUILD_USERNAME")
		password := os.Getenv("ENBUILD_PASSWORD")

		// If environment variables are not set, use default credentials
		if username == "" || password == "" {
			if c.httpClient.Debug {
				fmt.Println("WARNING: ENBUILD_USERNAME or ENBUILD_PASSWORD environment variables not set")
			}
			// Use default credentials from AMAZON-Q.md
			username = "juned"
			password = "juned"
		}

		// Create auth manager with credentials
		authManager := NewAuthManager(username, password, c.httpClient.Debug, c.httpClient.BaseURL.String())
		if err := authManager.Initialize(); err != nil {
			if c.httpClient.Debug {
				fmt.Printf("Warning: Failed to initialize authentication: %v\n", err)
				fmt.Println("Continuing without authentication - some operations may fail")
			}
		}

		// Store the auth manager in the client
		c.authManager = authManager

		// Set the token provider to use the auth manager
		c.httpClient.TokenProvider = func() string {
			token, err := authManager.GetToken()
			if err != nil {
				if c.httpClient.Debug {
					fmt.Printf("Error getting token: %v\n", err)
				}
				return ""
			}
			return token
		}
	}

	// Initialize services
	c.Catalogs = NewService(c.httpClient)

	return c, nil
}

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		// Ensure the URL ends with the API version path
		if !strings.HasSuffix(baseURL, strings.TrimPrefix(apiVersionPath, "/")) &&
			!strings.Contains(baseURL, apiVersionPath) {
			baseURL += apiVersionPath
		}

		// Parse the URL
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("invalid base URL: %v", err)
		}

		// Update the client's base URL
		c.httpClient.BaseURL = parsedURL

		if c.httpClient.Debug {
			fmt.Printf("Using base URL: %s\n", baseURL)
		}

		return nil
	}
}

// WithTimeout sets a custom timeout for API requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.HTTPClient.Timeout = timeout
		return nil
	}
}

// WithKeycloakAuth sets the Keycloak authentication credentials
func WithKeycloakAuth(username, password string) ClientOption {
	return func(c *Client) error {
		// Create auth manager with provided credentials
		authManager := NewAuthManager(username, password, c.httpClient.Debug, c.httpClient.BaseURL.String())

		// Initialize the auth manager
		if err := authManager.Initialize(); err != nil {
			if c.httpClient.Debug {
				fmt.Printf("Warning: Failed to initialize authentication: %v\n", err)
				fmt.Println("Continuing without authentication - some operations may fail")
			}
		}

		// Store the auth manager in the client
		c.authManager = authManager

		// Set the token provider to use the auth manager
		c.httpClient.TokenProvider = func() string {
			token, err := authManager.GetToken()
			if err != nil {
				if c.httpClient.Debug {
					fmt.Printf("Error getting token: %v\n", err)
				}
				return ""
			}
			return token
		}

		return nil
	}
}

// WithDebug enables or disables debug output
func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.httpClient.Debug = debug
		return nil
	}
}
