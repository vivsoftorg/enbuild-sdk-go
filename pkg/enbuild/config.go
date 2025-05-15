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

		// Ensure the base URL ends with a slash
		if !strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}

		// Ensure the base URL includes the API version path
		if !strings.HasSuffix(baseURL, apiVersionPath) && !strings.Contains(baseURL, apiVersionPath) {
			baseURL += apiVersionPath
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("invalid base URL: %v", err)
		}

		c.httpClient.BaseURL = parsedURL
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
				return fmt.Errorf("authentication token is required")
			}
		}
		c.httpClient.AuthToken = token
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
