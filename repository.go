package enbuild

import (
	"fmt"
	"net/http"
)

// RepositoryService handles communication with the repository related endpoints
type RepositoryService struct {
	client *Client
}

// Repository represents an ENBUILD repository
type Repository struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedOn   string                 `json:"createdOn,omitempty"`
	UpdatedOn   string                 `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// List returns a list of repositories
func (s *RepositoryService) List() ([]*Repository, error) {
	req, err := s.client.newRequest(http.MethodGet, "repository", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Repository `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Get returns a single repository by ID
func (s *RepositoryService) Get(id string) (*Repository, error) {
	path := fmt.Sprintf("repository/%s", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Repository `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
