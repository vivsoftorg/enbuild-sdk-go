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
	defaultBaseURL = "https://enbuild.vivplatform.io/enbuild-bk"
	defaultTimeout = 30 * time.Second
	apiVersionPath = "/api/v1/"
)

// Client represents the ENBUILD API client
type Client struct {
	httpClient *request.Client
	authManager *AuthManager

	// Services
	Catalogs *Service
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client) error

// NewClient creates a new ENBUILD API client
func NewClient(options ...ClientOption) (*Client, error) {
	// Process the default base URL to ensure it has the API version path
	defaultURLWithAPI := defaultBaseURL
	if !strings.HasSuffix(defaultURLWithAPI, strings.TrimPrefix(apiVersionPath, "/")) && 
	   !strings.Contains(defaultURLWithAPI, apiVersionPath) {
		defaultURLWithAPI += apiVersionPath
	}
	
	baseURL, _ := url.Parse(defaultURLWithAPI)
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

	// If no token provider was set and no auth manager exists, use Keycloak with default credentials
	if c.httpClient.TokenProvider == nil && c.authManager == nil {
		username := os.Getenv("ENBUILD_USERNAME")
		password := os.Getenv("ENBUILD_PASSWORD")
		
		// If environment variables are not set, use default credentials
		if username == "" {
			username = "default-user"
		}
		if password == "" {
			password = "default-password"
		}
		
		// Create auth manager with default credentials
		authManager := NewAuthManager(username, password, c.httpClient.Debug)
		if err := authManager.Initialize(); err != nil {
			return nil, fmt.Errorf("failed to initialize default Keycloak authentication: %v", err)
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
		
		if c.httpClient.Debug {
			fmt.Printf("Using default Keycloak authentication with username: %s\n", username)
		}
	}

	// Initialize services
	c.Catalogs = NewService(c.httpClient)

	return c, nil
}
