package enbuild

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ListStacks returns a list of Stacks.
// It accepts context, page, limit, and searchTerm for pagination and searching.
func (s *Enbuild) ListStacks(ctx context.Context, page int, limit int, searchTerm string) ([]*Stack, error) {
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
	return resp.Data, nil
}

