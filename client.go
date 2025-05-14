package enbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.enbuild.com/api/v1/"
	defaultTimeout = 30 * time.Second
	apiVersionPath = "/api/v1/"
)

// Client represents the ENBUILD API client
type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
	authToken  string
	debug      bool

	// Services
	Users         *UsersService
	Roles         *RolesService
	Operations    *OperationsService
	Repository    *RepositoryService
	Manifests     *ManifestsService
	MLDataset     *MLDatasetService
	AdminSettings *AdminSettingsService
	AuthLocal     *AuthLocalService
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client) error

// NewClient creates a new ENBUILD API client
func NewClient(options ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  "enbuild-sdk-go",
		httpClient: &http.Client{Timeout: defaultTimeout},
	}

	// Apply options
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	// Initialize services
	c.Users = &UsersService{client: c}
	c.Roles = &RolesService{client: c}
	c.Operations = &OperationsService{client: c}
	c.Repository = &RepositoryService{client: c}
	c.Manifests = &ManifestsService{client: c}
	c.MLDataset = &MLDatasetService{client: c}
	c.AdminSettings = &AdminSettingsService{client: c}
	c.AuthLocal = &AuthLocalService{client: c}

	return c, nil
}

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
		c.BaseURL = parsedURL
		return nil
	}
}

// WithTimeout sets the timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithAuthToken sets the authentication token for the API client
func WithAuthToken(token string) ClientOption {
	return func(c *Client) error {
		c.authToken = token
		return nil
	}
}

// WithDebug enables debug mode for the client
func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

// newRequest creates a new HTTP request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
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
	if c.authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	}

	return req, nil
}

// do sends an HTTP request and returns an HTTP response
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.debug {
		fmt.Printf("Making request to: %s %s\n", req.Method, req.URL.String())
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.debug {
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

// Response is a wrapper for API responses
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
