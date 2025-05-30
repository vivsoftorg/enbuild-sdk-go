package enbuild

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// List returns a list of catalogs.
func (s *Enbuild) ListCatalog(ctx context.Context, opts ...*CatalogListOptions) ([]*Catalog, error) {
	var options *CatalogListOptions
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	} else {
		// If options are for query params, this should be nil for a GET request body.
		// If options ARE the body, then this is fine.
		// For now, preserving original logic of passing it as body.
		options = &CatalogListOptions{}
	}

	// Assuming 'options' is intended as the body. If it's for query params,
	// path should be constructed dynamically and body should be nil.
	req, err := s.client.NewRequest(ctx, http.MethodGet, "manifests", options)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Catalog `json:"data"`
	}
	if _, err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	for _, catalog := range resp.Data {
		if id, ok := catalog.ID.(float64); ok {
			catalog.ID = fmt.Sprintf("%v", int64(id))
		} else if id, ok := catalog.ID.(string); ok {
			catalog.ID = id
		}
	}

	return s.filterCatalogs(resp.Data, options), nil
}

// Get returns a single catalog by ID.
func (s *Enbuild) GetCatalog(ctx context.Context, id string, opts *CatalogListOptions) (*Catalog, error) {
	if id == "" {
		return nil, fmt.Errorf("catalog ID is required")
	}

	// Assuming 'opts' is intended as the body.
	path := fmt.Sprintf("manifests/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, opts)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Catalog `json:"data"`
	}
	if _, err := s.client.Do(ctx, req, &resp); err != nil {
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

func (s *Enbuild) filterCatalogs(catalogs []*Catalog, opts *CatalogListOptions) []*Catalog {
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
