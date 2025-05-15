package manifests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vivsoftorg/enbuild-sdk-go/internal/request"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/types"
)

// Service handles communication with the manifests related endpoints
type Service struct {
	client *request.Client
}

// NewService creates a new manifests service
func NewService(client *request.Client) *Service {
	return &Service{client: client}
}

// ManifestListOptions specifies the optional parameters to the ManifestsService.List method
type ManifestListOptions struct {
	ID          string `url:"id,omitempty"`          // Filter by manifest ID
	VCS         string `url:"vcs,omitempty"`         // Filter by VCS type (github, gitlab)
	Type        string `url:"type,omitempty"`        // Filter by manifest type
	Slug        string `url:"slug,omitempty"`        // Filter by manifest slug
	Name        string `url:"name,omitempty"`        // Filter by manifest name (search)
	Description string `url:"description,omitempty"` // Filter by manifest description (search)
	Version     string `url:"version,omitempty"`     // Filter by manifest version
}

// List returns a list of manifests
func (s *Service) List(opts *ManifestListOptions) ([]*types.Manifest, error) {
	// Always use the standard manifests endpoint
	path := "manifests"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Get all manifests first
	manifests, err := s.getAllManifests(req)
	if err != nil {
		return nil, err
	}

	// If no options provided, return all manifests
	if opts == nil {
		return manifests, nil
	}

	// Filter manifests based on options
	return s.filterManifests(manifests, opts), nil
}

// Get returns a single manifest by ID
func (s *Service) Get(id string, opts *ManifestListOptions) (*types.Manifest, error) {
	if id == "" {
		return nil, fmt.Errorf("manifest ID is required")
	}

	// Create options with ID if not provided
	if opts == nil {
		opts = &ManifestListOptions{ID: id}
	} else {
		opts.ID = id
	}

	// Use List with ID filter to get a single manifest
	manifests, err := s.List(opts)
	if err != nil {
		return nil, err
	}

	if len(manifests) == 0 {
		return nil, fmt.Errorf("manifest with ID %s not found", id)
	}

	return manifests[0], nil
}

func (s *Service) getAllManifests(req *http.Request) ([]*types.Manifest, error) {
	// Try standard endpoint first
	standardManifests, err := s.getStandardManifests(req)
	if err != nil {
		if s.client.Debug {
			fmt.Printf("Error getting standard manifests: %v\n", err)
		}
		// If standard endpoint fails, return empty list
		standardManifests = []*types.Manifest{}
	}

	// Try GitHub endpoint
	githubReq, _ := s.client.NewRequest(http.MethodGet, "githubManifest", nil)
	githubManifests, err := s.getCustomEndpointManifests(githubReq)
	if err != nil && s.client.Debug {
		fmt.Printf("Error getting GitHub manifests: %v\n", err)
	}

	// Try GitLab endpoint
	gitlabReq, _ := s.client.NewRequest(http.MethodGet, "gitlabManifest", nil)
	gitlabManifests, err := s.getCustomEndpointManifests(gitlabReq)
	if err != nil && s.client.Debug {
		fmt.Printf("Error getting GitLab manifests: %v\n", err)
	}

	// Combine all manifests
	allManifests := append(standardManifests, githubManifests...)
	allManifests = append(allManifests, gitlabManifests...)

	return allManifests, nil
}

func (s *Service) filterManifests(manifests []*types.Manifest, opts *ManifestListOptions) []*types.Manifest {
	if opts == nil {
		return manifests
	}

	var filtered []*types.Manifest
	for _, m := range manifests {
		// Filter by ID if provided
		if opts.ID != "" {
			ID, ok := m.ID.(string)
			if !ok || ID != opts.ID {
				continue
			}
		}

		// Filter by VCS if provided
		if opts.VCS != "" && !strings.EqualFold(m.VCS, opts.VCS) {
			continue
		}

		// Filter by Type if provided
		if opts.Type != "" && !strings.EqualFold(m.Type, opts.Type) {
			continue
		}

		// Filter by Slug if provided
		if opts.Slug != "" && !strings.EqualFold(m.Slug, opts.Slug) {
			continue
		}

		// Filter by Name if provided (case-insensitive contains)
		if opts.Name != "" && !strings.Contains(strings.ToLower(m.Name), strings.ToLower(opts.Name)) {
			continue
		}

		// Filter by Description if provided (case-insensitive contains)
		if opts.Description != "" && !strings.Contains(
			strings.ToLower(m.Description),
			strings.ToLower(opts.Description)) {
			continue
		}

		// Filter by Version if provided
		if opts.Version != "" && !strings.EqualFold(m.Version, opts.Version) {
			continue
		}

		// If all filters pass, add to filtered list
		filtered = append(filtered, m)
	}

	return filtered
}

func (s *Service) getStandardManifests(req *http.Request) ([]*types.Manifest, error) {
	var resp struct {
		Data []*types.Manifest `json:"data"`
	}
	if _, err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	// Convert ID to string for all manifests
	for _, manifest := range resp.Data {
		if id, ok := manifest.ID.(float64); ok {
			manifest.ID = fmt.Sprintf("%v", int64(id))
		} else if id, ok := manifest.ID.(string); ok {
			manifest.ID = id
		}
	}

	return resp.Data, nil
}

func (s *Service) getCustomEndpointManifests(req *http.Request) ([]*types.Manifest, error) {
	var rawResp map[string]interface{}
	if _, err := s.client.Do(req, &rawResp); err != nil {
		return nil, err
	}

	if s.client.Debug {
		fmt.Printf("Raw response from API: %+v\n", rawResp)
	}

	dataField, dataExists := rawResp["data"].(map[string]interface{})
	if !dataExists {
		return nil, fmt.Errorf("data field missing in response")
	}

	catalogArray, catalogExists := dataField["catalogManifest"].([]interface{})
	if !catalogExists {
		return nil, fmt.Errorf("catalogManifest field missing in response data")
	}

	return s.parseManifests(catalogArray), nil
}

func (s *Service) parseManifests(catalogArray []interface{}) []*types.Manifest {
	var manifests []*types.Manifest

	for _, item := range catalogArray {
		manifestMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Convert ID to string format if it exists
		if id, ok := manifestMap["_id"]; ok {
			manifestMap["_id"] = fmt.Sprintf("%v", id) // Explicitly convert to string
		}

		manifestBytes, _ := json.Marshal(manifestMap)
		var manifest types.Manifest
		if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
			if s.client.Debug {
				fmt.Printf("Error unmarshaling manifest: %v\n", err)
			}
			continue
		}

		if s.client.Debug {
			fmt.Printf("Parsed manifest: %+v\n", manifest)
		}

		manifests = append(manifests, &manifest)
	}

	return manifests
}
