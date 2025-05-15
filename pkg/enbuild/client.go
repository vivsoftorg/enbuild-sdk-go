package enbuild

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
)

const (
	defaultBaseURL = "https://enbuild.vivplatform.io/enbuild-bk"
	defaultTimeout = 30 * time.Second
	apiVersionPath = "/api/v1/"
	defaultToken   = "default-token" // Default token to use if none provided
)

// Client represents the ENBUILD API client
type Client struct {
	httpClient *request.Client

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
		AuthToken:  defaultToken, // Set default token
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

	// Initialize services
	c.Catalogs = NewService(c.httpClient)

	return c, nil
}
