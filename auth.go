package enbuild

import (
	"fmt"
	"net/http"
)

// AuthLocalService handles communication with the authentication related endpoints
type AuthLocalService struct {
	client *Client
}

// LocalAdmin represents a local admin user
type LocalAdmin struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	CreatedOn string `json:"createdOn,omitempty"`
	UpdatedOn string `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// Create creates a new local admin user
func (s *AuthLocalService) Create(admin *LocalAdmin) (*LocalAdmin, error) {
	req, err := s.client.newRequest(http.MethodPost, "authLocal", admin)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *LocalAdmin `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Update updates an existing local admin user
func (s *AuthLocalService) Update(id string, admin *LocalAdmin) (*LocalAdmin, error) {
	path := fmt.Sprintf("authLocal/%s", id)
	req, err := s.client.newRequest(http.MethodPut, path, admin)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *LocalAdmin `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
