package enbuild

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	// "strings"
)

// ListStacks returns a list of Stacks.
// It accepts context, page, limit, and searchTerm for pagination and searching.
func (s *Enbuild) ListStacks(ctx context.Context, page int, limit int, searchTerm string) ([]*Stack, error) {
	// Construct the request path with query parameters
	// Ensure searchTerm is URL-encoded to handle special characters
	encodedSearchTerm := url.QueryEscape(searchTerm)
	path := fmt.Sprintf("stacks?page=%d&limit=%d&search=%s", page, limit, encodedSearchTerm)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*Stack `json:"data"`
	}
	if _, err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	// Directly return the received data as it does not require any modification.
	return resp.Data, nil
}

// // Get returns a single Stack by ID.
// func (s *Enbuild) Get(id string, opts *Stack) (*Stack, error) {
// 	if id == "" {
// 		return nil, fmt.Errorf("Stack ID is required")
// 	}

// 	path := fmt.Sprintf("manifests/%s", id)
// 	req, err := s.client.NewRequest(http.MethodGet, path, opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var resp struct {
// 		Data []*Stack `json:"data"`
// 	}
// 	if _, err := s.client.Do(req, &resp); err != nil {
// 		return nil, err
// 	}

// 	for _, Stack := range resp.Data {
// 		if id, ok := Stack.ID.(float64); ok {
// 			Stack.ID = fmt.Sprintf("%v", int64(id))
// 		} else if id, ok := Stack.ID.(string); ok {
// 			Stack.ID = id
// 		}
// 	}

// 	return resp.Data[0], nil
// }

// func (s *Enbuild) filterStacks(Stacks []*Stack, opts *Stack) []*Stack {
// 	if opts == nil {
// 		return Stacks
// 	}

// 	var filtered []*Stack
// 	for _, m := range Stacks {
// 		if opts.ID != "" && !strings.EqualFold(m.ID.(string), opts.ID) {
// 			continue
// 		}
// 		if opts.VCS != "" && !strings.EqualFold(m.VCS, opts.VCS) {
// 			continue
// 		}
// 		if opts.Type != "" && !strings.EqualFold(m.Type, opts.Type) {
// 			continue
// 		}
// 		if opts.Slug != "" && !strings.EqualFold(m.Slug, opts.Slug) {
// 			continue
// 		}
// 		if opts.Name != "" && !strings.Contains(strings.ToLower(m.Name), strings.ToLower(opts.Name)) {
// 			continue
// 		}
// 		if opts.Description != "" && !strings.Contains(strings.ToLower(m.Description), strings.ToLower(opts.Description)) {
// 			continue
// 		}
// 		if opts.Version != "" && !strings.EqualFold(m.Version, opts.Version) {
// 			continue
// 		}

// 		filtered = append(filtered, m)
// 	}

// 	return filtered
// }
