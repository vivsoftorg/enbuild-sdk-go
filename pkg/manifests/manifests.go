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
	// Determine the path based on VCS option
	path := "manifests"
	if opts != nil && opts.VCS != "" {
		vcsType := strings.ToLower(opts.VCS)
		if vcsType == "github" {
			path = "githubManifest"
		} else if vcsType == "gitlab" {
			path = "gitlabManifest"
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Get manifests
	var manifests []*types.Manifest

	// Handle different response formats based on the endpoint
	if path == "manifests" {
		// Standard endpoint
		manifests, err = s.getManifests(req)
	} else {
		// VCS-specific endpoints (GitHub/GitLab)
		manifests, err = s.getVCSManifests(req)
	}

	if err != nil {
		return nil, err
	}

	// If no options provided or only VCS filter was used, return all manifests
	if opts == nil || (opts.ID == "" && opts.Type == "" && opts.Slug == "" &&
		opts.Name == "" && opts.Description == "" && opts.Version == "") {
		return manifests, nil
	}

	// Filter manifests based on other options
	return s.filterManifests(manifests, opts), nil
}

// // Get returns a single manifest by ID
// func (s *Service) Get(id string, opts *ManifestListOptions) (*types.Manifest, error) {
// 	if id == "" {
// 		return nil, fmt.Errorf("manifest ID is required")
// 	}

// 	// Append ID to path
// 	path = fmt.Sprintf("%s/%s", "manifests", id)

// 	req, err := s.client.NewRequest(http.MethodGet, path, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Handle different response formats based on the endpoint
// 	if path == fmt.Sprintf("manifests/%s", id) {
// 		// Standard endpoint
// 		var resp struct {
// 			Data *types.Manifest `json:"data"`
// 		}
// 		if _, err := s.client.Do(req, &resp); err != nil {
// 			return nil, err
// 		}

// 		// Convert ID to string if needed
// 		if id, ok := resp.Data.ID.(float64); ok {
// 			resp.Data.ID = fmt.Sprintf("%v", int64(id))
// 		}

// 		return resp.Data, nil
// 	} else {
// 		// VCS-specific endpoint
// 		var rawResp map[string]interface{}
// 		if _, err := s.client.Do(req, &rawResp); err != nil {
// 			return nil, err
// 		}

// 		if s.client.Debug {
// 			fmt.Printf("Raw response from API: %+v\n", rawResp)
// 		}

// 		// Extract manifest from the data field
// 		if dataField, ok := rawResp["data"].(map[string]interface{}); ok {
// 			// Check for catalogManifest field first
// 			if catalogManifestField, ok := dataField["catalogManifest"].(map[string]interface{}); ok {
// 				// Convert ID to string format if it exists
// 				if id, ok := catalogManifestField["_id"]; ok {
// 					catalogManifestField["_id"] = fmt.Sprintf("%v", id) // Explicitly convert to string
// 				}

// 				manifestBytes, _ := json.Marshal(catalogManifestField)
// 				var manifest types.Manifest
// 				if err := json.Unmarshal(manifestBytes, &manifest); err == nil {
// 					return &manifest, nil
// 				}
// 			}

// 			// If no catalogManifest field, try using the data object itself
// 			// Convert ID to string format if it exists
// 			if id, ok := dataField["_id"]; ok {
// 				dataField["_id"] = fmt.Sprintf("%v", id) // Explicitly convert to string
// 			}

// 			manifestBytes, _ := json.Marshal(dataField)
// 			var manifest types.Manifest
// 			if err := json.Unmarshal(manifestBytes, &manifest); err == nil && manifest.ID != nil {
// 				return &manifest, nil
// 			}
// 		}

// 		return nil, fmt.Errorf("could not parse manifest response")
// 	}
// }

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

		// Filter by VCS if provided (already filtered by endpoint selection, but double-check)
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

func (s *Service) getManifests(req *http.Request) ([]*types.Manifest, error) {
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

func (s *Service) getVCSManifests(req *http.Request) ([]*types.Manifest, error) {
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
