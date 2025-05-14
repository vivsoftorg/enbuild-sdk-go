package enbuild

import (
	"fmt"
	"net/http"
	"net/url"
)

// UsersService handles communication with the user related endpoints
type UsersService struct {
	client *Client
}

// User represents an ENBUILD user
type User struct {
	ID        string      `json:"id,omitempty"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName,omitempty"`
	LastName  string      `json:"lastName,omitempty"`
	CreatedBy string      `json:"createdBy,omitempty"`
	CreatedOn interface{} `json:"createdOn,omitempty"`
	UpdatedOn interface{} `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// UserListOptions specifies the optional parameters to the UsersService.List method
type UserListOptions struct {
	CreatedBy string `url:"createdBy,omitempty"`
	Limit     int    `url:"limit,omitempty"`
	Page      int    `url:"page,omitempty"`
	Sort      string `url:"sort,omitempty"`
}

// List returns a list of users
func (s *UsersService) List(opts *UserListOptions) ([]*User, error) {
	path := "users"
	if opts != nil {
		v := url.Values{}
		if opts.CreatedBy != "" {
			v.Add("createdBy", opts.CreatedBy)
		}
		if opts.Limit > 0 {
			v.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Page > 0 {
			v.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Sort != "" {
			v.Add("sort", opts.Sort)
		}
		if len(v) > 0 {
			path = fmt.Sprintf("%s?%s", path, v.Encode())
		}
	}

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*User `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Get returns a single user by ID
func (s *UsersService) Get(id string) (*User, error) {
	path := fmt.Sprintf("users/%s", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *User `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Create creates a new user
func (s *UsersService) Create(user *User) (*User, error) {
	req, err := s.client.newRequest(http.MethodPost, "users", user)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *User `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Update updates an existing user
func (s *UsersService) Update(id string, user *User) (*User, error) {
	path := fmt.Sprintf("users/%s", id)
	req, err := s.client.newRequest(http.MethodPut, path, user)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *User `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
