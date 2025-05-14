package enbuild

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUsersList(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test request parameters
		if r.URL.Path != "/users" {
			t.Errorf("Expected URL path %v, got %v", "/users", r.URL.Path)
		}
		
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": [
				{
					"id": "1",
					"username": "user1",
					"email": "user1@example.com"
				},
				{
					"id": "2",
					"username": "user2",
					"email": "user2@example.com"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	baseURL, _ := url.Parse(server.URL + "/")
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  "enbuild-sdk-go",
		httpClient: server.Client(),
	}
	c.Users = &UsersService{client: c}

	// Test List method
	users, err := c.Users.List(nil)
	if err != nil {
		t.Fatalf("Users.List returned unexpected error: %v", err)
	}

	// Check response
	if len(users) != 2 {
		t.Errorf("Users.List returned %d users, expected %d", len(users), 2)
	}

	// Check first user
	if users[0].ID != "1" {
		t.Errorf("Users.List first user ID = %v, expected %v", users[0].ID, "1")
	}
	if users[0].Username != "user1" {
		t.Errorf("Users.List first user Username = %v, expected %v", users[0].Username, "user1")
	}
	if users[0].Email != "user1@example.com" {
		t.Errorf("Users.List first user Email = %v, expected %v", users[0].Email, "user1@example.com")
	}
}

func TestUsersGet(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test request parameters
		if r.URL.Path != "/users/1" {
			t.Errorf("Expected URL path %v, got %v", "/users/1", r.URL.Path)
		}
		
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": {
				"id": "1",
				"username": "user1",
				"email": "user1@example.com",
				"firstName": "John",
				"lastName": "Doe"
			}
		}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	baseURL, _ := url.Parse(server.URL + "/")
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  "enbuild-sdk-go",
		httpClient: server.Client(),
	}
	c.Users = &UsersService{client: c}

	// Test Get method
	user, err := c.Users.Get("1")
	if err != nil {
		t.Fatalf("Users.Get returned unexpected error: %v", err)
	}

	// Check response
	if user.ID != "1" {
		t.Errorf("Users.Get ID = %v, expected %v", user.ID, "1")
	}
	if user.Username != "user1" {
		t.Errorf("Users.Get Username = %v, expected %v", user.Username, "user1")
	}
	if user.Email != "user1@example.com" {
		t.Errorf("Users.Get Email = %v, expected %v", user.Email, "user1@example.com")
	}
	if user.FirstName != "John" {
		t.Errorf("Users.Get FirstName = %v, expected %v", user.FirstName, "John")
	}
	if user.LastName != "Doe" {
		t.Errorf("Users.Get LastName = %v, expected %v", user.LastName, "Doe")
	}
}

func TestUsersCreate(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test request parameters
		if r.URL.Path != "/users" {
			t.Errorf("Expected URL path %v, got %v", "/users", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method %v, got %v", http.MethodPost, r.Method)
		}
		
		// Test request body
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		if user.Username != "newuser" {
			t.Errorf("Expected username %v, got %v", "newuser", user.Username)
		}
		if user.Email != "newuser@example.com" {
			t.Errorf("Expected email %v, got %v", "newuser@example.com", user.Email)
		}
		
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": {
				"id": "3",
				"username": "newuser",
				"email": "newuser@example.com"
			}
		}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	baseURL, _ := url.Parse(server.URL + "/")
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  "enbuild-sdk-go",
		httpClient: server.Client(),
	}
	c.Users = &UsersService{client: c}

	// Test Create method
	newUser := &User{
		Username: "newuser",
		Email:    "newuser@example.com",
	}
	user, err := c.Users.Create(newUser)
	if err != nil {
		t.Fatalf("Users.Create returned unexpected error: %v", err)
	}

	// Check response
	if user.ID != "3" {
		t.Errorf("Users.Create ID = %v, expected %v", user.ID, "3")
	}
	if user.Username != "newuser" {
		t.Errorf("Users.Create Username = %v, expected %v", user.Username, "newuser")
	}
	if user.Email != "newuser@example.com" {
		t.Errorf("Users.Create Email = %v, expected %v", user.Email, "newuser@example.com")
	}
}
