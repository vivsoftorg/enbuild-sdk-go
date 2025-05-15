package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client represents an HTTP client for making API requests
type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	HTTPClient *http.Client
	AuthToken  string
	Debug      bool
}

// NewRequest creates a new HTTP request
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	// Add auth token if available
	if c.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	}

	return req, nil
}

// Do sends an HTTP request and returns an HTTP response
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.Debug {
		fmt.Printf("Making request to: %s %s\n", req.Method, req.URL.String())
		
		// Print request headers in debug mode
		fmt.Println("Request headers:")
		for key, values := range req.Header {
			// Mask the Authorization header for security
			if strings.ToLower(key) == "authorization" {
				fmt.Printf("  %s: Bearer ****\n", key)
			} else {
				fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
			}
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.Debug {
		// Read the response body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Response from %s: %s\n", req.URL.String(), bodyString)

		// Create a new reader with the same data for the JSON decoder
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("API error: %s", resp.Status)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}
