package enbuild

import (
	"net/http"
)

// MLDatasetService handles communication with the ML dataset related endpoints
type MLDatasetService struct {
	client *Client
}

// MLDataset represents an ENBUILD ML dataset
type MLDataset struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Size        int64                  `json:"size,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedOn   string                 `json:"createdOn,omitempty"`
	UpdatedOn   string                 `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// List returns a list of ML datasets
func (s *MLDatasetService) List() ([]*MLDataset, error) {
	req, err := s.client.newRequest(http.MethodGet, "mlDataset", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*MLDataset `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
