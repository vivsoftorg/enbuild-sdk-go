package enbuild

import (
	"fmt"
	"net/url"
	"os"
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

// WithAuthToken sets the authentication token for API requests
func WithAuthToken(token string) ClientOption {
	return func(c *Client) error {
		if token == "" {
			// Try to get token from environment variable
			token = os.Getenv("ENBUILD_API_TOKEN")
			if token == "" {
				// Use default token if environment variable is not set
				token = defaultToken
				if c.httpClient.Debug {
					fmt.Printf("Using default token: %s\n", token)
				}
			} else if c.httpClient.Debug {
				fmt.Printf("Using token from environment variable\n")
			}
		} else if c.httpClient.Debug {
			fmt.Printf("Using provided token\n")
		}
		
		c.httpClient.AuthToken = token
		
		// Log token in debug mode (masked for security)
		if c.httpClient.Debug {
			maskedToken := maskToken(token)
			fmt.Printf("Auth token: %s\n", maskedToken)
		}
		
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

// maskToken masks a token for secure logging
func maskToken(token string) string {
	if len(token) <= 4 {
		return "****"
	}
	return token[:2] + "****" + token[len(token)-2:]
}
