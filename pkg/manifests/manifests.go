package manifests

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/types"
)

// Service handles communication with the manifests-related endpoints.
type Service struct {
	client *request.Client
}

// NewService creates a new manifests service.
func NewService(client *request.Client) *Service {
	return &Service{client: client}
}

// ManifestListOptions specifies the optional parameters to the ManifestsService.List method.
type ManifestListOptions struct {
	ID          string `url:"id,omitempty"`
	VCS         string `url:"vcs,omitempty"`
	Type        string `url:"type,omitempty"`
	Slug        string `url:"slug,omitempty"`
	Name        string `url:"name,omitempty"`
	Description string `url:"description,omitempty"`
	Version     string `url:"version,omitempty"`
}

// List returns a list of manifests.
func (s *Service) List(opts *ManifestListOptions) ([]*types.Manifest, error) {
	req, err := s.client.NewRequest(http.MethodGet, "manifests", opts)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*types.Manifest `json:"data"`
	}
	if _, err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	for _, manifest := range resp.Data {
		if id, ok := manifest.ID.(float64); ok {
			manifest.ID = fmt.Sprintf("%v", int64(id))
		} else if id, ok := manifest.ID.(string); ok {
			manifest.ID = id
		}
	}

	return s.filterManifests(resp.Data, opts), nil
}

// Get returns a single manifest by ID.
func (s *Service) Get(id string, opts *ManifestListOptions) (*types.Manifest, error) {
	if id == "" {
		return nil, fmt.Errorf("manifest ID is required")
	}

	path := fmt.Sprintf("manifests/%s", id)
	req, err := s.client.NewRequest(http.MethodGet, path, opts)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []*types.Manifest `json:"data"`
	}
	if _, err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	for _, manifest := range resp.Data {
		if id, ok := manifest.ID.(float64); ok {
			manifest.ID = fmt.Sprintf("%v", int64(id))
		} else if id, ok := manifest.ID.(string); ok {
			manifest.ID = id
		}
	}

	return resp.Data[0], nil
}

func (s *Service) filterManifests(manifests []*types.Manifest, opts *ManifestListOptions) []*types.Manifest {
	if opts == nil {
		return manifests
	}

	var filtered []*types.Manifest
	for _, m := range manifests {
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
