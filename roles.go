package enbuild

import (
	"fmt"
	"net/http"
)

// RolesService handles communication with the role related endpoints
type RolesService struct {
	client *Client
}

// Role represents an ENBUILD role
type Role struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Permissions []string    `json:"permissions,omitempty"`
	CreatedOn   interface{} `json:"createdOn,omitempty"`
	UpdatedOn   interface{} `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// List returns a list of roles
func (s *RolesService) List() ([]*Role, error) {
	req, err := s.client.newRequest(http.MethodGet, "roles", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Role `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Create creates a new role
func (s *RolesService) Create(role *Role) (*Role, error) {
	req, err := s.client.newRequest(http.MethodPost, "roles", role)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Role `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Update updates an existing role
func (s *RolesService) Update(id string, role *Role) (*Role, error) {
	path := fmt.Sprintf("roles/%s", id)
	req, err := s.client.newRequest(http.MethodPut, path, role)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Role `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// CheckAuth checks if the user's token is valid
func (s *RolesService) CheckAuth() (bool, error) {
	req, err := s.client.newRequest(http.MethodGet, "roles/auth", nil)
	if err != nil {
		return false, err
	}

	resp, err := s.client.do(req, nil)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}
