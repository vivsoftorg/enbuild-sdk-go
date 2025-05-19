package enbuild

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// WithBaseURL sets the base URL for the API client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		if baseURL == "" {
			return nil
		}

		// Remove trailing slash if present
		baseURL = strings.TrimSuffix(baseURL, "/")

		// Ensure the base URL includes the API version path
		if !strings.HasSuffix(baseURL, strings.TrimPrefix(apiVersionPath, "/")) && !strings.Contains(baseURL, apiVersionPath) {
			baseURL += apiVersionPath
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("invalid base URL: %v", err)
		}

		c.httpClient.BaseURL = parsedURL
		
		// Log base URL in debug mode
		if c.httpClient.Debug {
			fmt.Printf("Using base URL: %s\n", parsedURL.String())
		}
		
		return nil
	}
}

// WithTimeout sets the timeout for API requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.HTTPClient.Timeout = timeout
		return nil
	}
}

// WithDebug enables or disables debug mode
func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.httpClient.Debug = debug
		return nil
	}
}

// WithKeycloakAuth creates a client option that configures authentication using Keycloak
func WithKeycloakAuth(username, password string) ClientOption {
	return func(c *Client) error {
		authManager := NewAuthManager(username, password, c.httpClient.Debug)
		if err := authManager.Initialize(); err != nil {
			return err
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

// maskToken masks a token for secure logging
func maskToken(token string) string {
	if len(token) <= 4 {
		return "****"
	}
	return token[:2] + "****" + token[len(token)-2:]
}
