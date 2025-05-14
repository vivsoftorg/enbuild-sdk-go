package enbuild

import (
	"fmt"
	"net/http"
)

// ManifestsService handles communication with the manifests related endpoints
type ManifestsService struct {
	client *Client
}

// Manifest represents an ENBUILD manifest
type Manifest struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	Version     string                 `json:"version,omitempty"`
	CreatedOn   string                 `json:"createdOn,omitempty"`
	UpdatedOn   string                 `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// List returns a list of manifests
func (s *ManifestsService) List() ([]*Manifest, error) {
	req, err := s.client.newRequest(http.MethodGet, "manifests", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Manifest `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Get returns a single manifest by ID
func (s *ManifestsService) Get(id string) (*Manifest, error) {
	path := fmt.Sprintf("manifests/%s", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Manifest `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
