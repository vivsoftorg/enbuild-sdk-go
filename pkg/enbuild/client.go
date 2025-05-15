package enbuild

import (
	"net/http"
	"net/url"
	"time"

	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/manifests"
)

const (
	defaultBaseURL = "https://enbuild-dev.vivplatform.io/enbuild-bk"
	defaultTimeout = 30 * time.Second
	apiVersionPath = "/api/v1/"
)

// Client represents the ENBUILD API client
type Client struct {
	httpClient *request.Client

	// Services
	Manifests *manifests.Service
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client) error

// NewClient creates a new ENBUILD API client
func NewClient(options ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
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

	// Initialize services
	c.Manifests = manifests.NewService(c.httpClient)

	return c, nil
}
