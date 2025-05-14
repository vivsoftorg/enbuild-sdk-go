package enbuild

import (
	"fmt"
	"net/http"
	"net/url"
)

// OperationsService handles communication with the operations related endpoints
type OperationsService struct {
	client *Client
}

// Operation represents an ENBUILD operation
type Operation struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy   string                 `json:"createdBy,omitempty"`
	CreatedOn   string                 `json:"createdOn,omitempty"`
	UpdatedOn   string                 `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// OperationListOptions specifies the optional parameters to the OperationsService.List method
type OperationListOptions struct {
	Limit int    `url:"limit,omitempty"`
	Page  int    `url:"page,omitempty"`
	Sort  string `url:"sort,omitempty"`
	// Add other filter options as needed
}

// List returns a list of operations
func (s *OperationsService) List(opts *OperationListOptions) ([]*Operation, error) {
	path := "operations"
	if opts != nil {
		v := url.Values{}
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
		Data []*Operation `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Get returns a single operation by ID
func (s *OperationsService) Get(id string) (*Operation, error) {
	path := fmt.Sprintf("operations/%s", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Operation `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Create creates a new operation
func (s *OperationsService) Create(operation *Operation) (*Operation, error) {
	req, err := s.client.newRequest(http.MethodPost, "operations", operation)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Operation `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Update updates an existing operation
func (s *OperationsService) Update(id string, operation *Operation) (*Operation, error) {
	path := fmt.Sprintf("operations/%s", id)
	req, err := s.client.newRequest(http.MethodPut, path, operation)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *Operation `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
