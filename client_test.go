package enbuild

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL.String(), defaultBaseURL)
	}

	if c.UserAgent != "enbuild-sdk-go" {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, "enbuild-sdk-go")
	}
}

func TestWithBaseURL(t *testing.T) {
	customURL := "https://custom-api.enbuild.com/api/v1/"
	c, err := NewClient(WithBaseURL(customURL))
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if c.BaseURL.String() != customURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL.String(), customURL)
	}
}

func TestWithTimeout(t *testing.T) {
	customTimeout := 60 * time.Second
	c, err := NewClient(WithTimeout(customTimeout))
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if c.httpClient.Timeout != customTimeout {
		t.Errorf("NewClient Timeout = %v, expected %v", c.httpClient.Timeout, customTimeout)
	}
}

func TestWithAuthToken(t *testing.T) {
	token := "test-token"
	c, err := NewClient(WithAuthToken(token))
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if c.authToken != token {
		t.Errorf("NewClient authToken = %v, expected %v", c.authToken, token)
	}
}

func TestNewRequest(t *testing.T) {
	c, _ := NewClient(WithAuthToken("test-token"))
	
	req, err := c.newRequest(http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("newRequest returned unexpected error: %v", err)
	}

	// Test that the auth token is set
	if got, want := req.Header.Get("Authorization"), "Bearer test-token"; got != want {
		t.Errorf("newRequest Authorization header = %v, expected %v", got, want)
	}

	// Test that the user agent is set
	if got, want := req.Header.Get("User-Agent"), "enbuild-sdk-go"; got != want {
		t.Errorf("newRequest User-Agent header = %v, expected %v", got, want)
	}
}

func TestDo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data": {"id": "123", "name": "test"}}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	baseURL, _ := url.Parse(server.URL + "/")
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  "enbuild-sdk-go",
		httpClient: server.Client(),
	}

	req, _ := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	
	var resp struct {
		Data struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	
	_, err := c.do(req, &resp)
	if err != nil {
		t.Fatalf("do returned unexpected error: %v", err)
	}

	if resp.Data.ID != "123" {
		t.Errorf("do response ID = %v, expected %v", resp.Data.ID, "123")
	}

	if resp.Data.Name != "test" {
		t.Errorf("do response Name = %v, expected %v", resp.Data.Name, "test")
	}
}
