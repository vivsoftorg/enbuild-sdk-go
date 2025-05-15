package enbuild

import (
	"net/url"
	"strings"
	"time"
)

// WithBaseURL sets the base URL for the API client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		// Ensure the base URL ends with the API version path
		if !strings.HasSuffix(baseURL, apiVersionPath) {
			if !strings.HasSuffix(baseURL, "/") {
				baseURL += "/"
			}
			if !strings.HasSuffix(baseURL, "api/v1/") {
				baseURL += "api/v1/"
			}
		}
		
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.httpClient.BaseURL = parsedURL
		return nil
	}
}

// WithTimeout sets the timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.HTTPClient.Timeout = timeout
		return nil
	}
}

// WithAuthToken sets the authentication token for the API client
func WithAuthToken(token string) ClientOption {
	return func(c *Client) error {
		c.httpClient.AuthToken = token
		return nil
	}
}

// WithDebug enables debug mode for the client
func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.httpClient.Debug = debug
		return nil
	}
}
