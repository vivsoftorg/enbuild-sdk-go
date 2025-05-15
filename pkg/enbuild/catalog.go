package enbuild

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
)

// Service handles communication with the catalogs-related endpoints.
type Service struct {
	client *request.Client
}

// NewService creates a new catalogs service.
func NewService(client *request.Client) *Service {
	return &Service{client: client}
}



// List returns a list of catalogs.
func (s *Service) List(opts *CatalogListOptions) ([]*Catalog, error) {
	req, err := s.client.NewRequest(http.MethodGet, "manifests", opts)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Catalog `json:"data"`
	}
	if _, err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	for _, catalog := range resp.Data {
		if id, ok := catalog.ID.(float64); ok {
			catalog.ID = fmt.Sprintf("%v", int64(id))
		} else if id, ok := catalog.ID.(string); ok {
			catalog.ID = id
		}
	}

	return s.filterCatalogs(resp.Data, opts), nil
}

// Get returns a single catalog by ID.
func (s *Service) Get(id string, opts *CatalogListOptions) (*Catalog, error) {
	if id == "" {
		return nil, fmt.Errorf("catalog ID is required")
	}

	path := fmt.Sprintf("manifests/%s", id)
	req, err := s.client.NewRequest(http.MethodGet, path, opts)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Catalog `json:"data"`
	}
	if _, err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	for _, catalog := range resp.Data {
		if id, ok := catalog.ID.(float64); ok {
			catalog.ID = fmt.Sprintf("%v", int64(id))
		} else if id, ok := catalog.ID.(string); ok {
			catalog.ID = id
		}
	}

	return resp.Data[0], nil
}

func (s *Service) filterCatalogs(catalogs []*Catalog, opts *CatalogListOptions) []*Catalog {
	if opts == nil {
		return catalogs
	}

	var filtered []*Catalog
	for _, m := range catalogs {
		if opts.ID != "" && !strings.EqualFold(m.ID.(string), opts.ID) {
			continue
		}
		if opts.VCS != "" && !strings.EqualFold(m.VCS, opts.VCS) {
			continue
		}
		if opts.Type != "" && !strings.EqualFold(m.Type, opts.Type) {
			continue
		}
		if opts.Slug != "" && !strings.EqualFold(m.Slug, opts.Slug) {
			continue
		}
		if opts.Name != "" && !strings.Contains(strings.ToLower(m.Name), strings.ToLower(opts.Name)) {
			continue
		}
		if opts.Description != "" && !strings.Contains(strings.ToLower(m.Description), strings.ToLower(opts.Description)) {
			continue
		}
		if opts.Version != "" && !strings.EqualFold(m.Version, opts.Version) {
			continue
		}

		filtered = append(filtered, m)
	}

	return filtered
}
