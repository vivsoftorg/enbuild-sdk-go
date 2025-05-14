package enbuild

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// VCSType represents the type of version control system
type VCSType string

const (
	// VCSTypeGitHub represents GitHub VCS
	VCSTypeGitHub VCSType = "GITHUB"
	// VCSTypeGitLab represents GitLab VCS
	VCSTypeGitLab VCSType = "GITLAB"
)

// ManifestsService handles communication with the manifests related endpoints
type ManifestsService struct {
	client *Client
}

// Manifest represents an ENBUILD manifest
type Manifest struct {
	ID          interface{}            `json:"id,omitempty"`
	MongoID     interface{}            `json:"_id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	Version     string                 `json:"version,omitempty"`
	CreatedOn   interface{}            `json:"createdOn,omitempty"`
	UpdatedOn   interface{}            `json:"updatedOn,omitempty"`
	VCS         VCSType                `json:"vcs,omitempty"`
	Slug        string                 `json:"slug,omitempty"`
	Type        string                 `json:"type,omitempty"`
	// Add other fields as needed
}

// ManifestListOptions specifies the optional parameters to the ManifestsService.List method
type ManifestListOptions struct {
	VCS VCSType `url:"vcs,omitempty"`
}

// List returns a list of manifests
func (s *ManifestsService) List(opts *ManifestListOptions) ([]*Manifest, error) {
	path := "manifests"
	
	// Use specific VCS endpoint if provided
	if opts != nil && opts.VCS != "" {
		vcsType := strings.ToLower(string(opts.VCS))
		if vcsType == "github" {
			path = "githubManifest"
		} else if vcsType == "gitlab" {
			path = "gitlabManifest"
		}
	}
	
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// For GitHub/GitLab endpoints, the response structure is different
	if opts != nil && opts.VCS != "" {
		// Parse as a raw map to inspect the structure
		var rawResp map[string]interface{}
		_, err = s.client.do(req, &rawResp)
		if err != nil {
			return nil, err
		}

		if s.client.debug {
			fmt.Printf("Raw response from API: %+v\n", rawResp)
		}

		// Extract manifests from the catalogManifest field in the data
		var manifests []*Manifest
		
		if dataField, ok := rawResp["data"]; ok {
			if dataObj, ok := dataField.(map[string]interface{}); ok {
				if catalogManifestField, ok := dataObj["catalogManifest"]; ok {
					if catalogArray, ok := catalogManifestField.([]interface{}); ok {
						for _, item := range catalogArray {
							if manifestMap, ok := item.(map[string]interface{}); ok {
								// Set ID field from either "id" or "_id" field
								if id, ok := manifestMap["id"]; ok {
									manifestMap["ID"] = id
								} else if id, ok := manifestMap["_id"]; ok {
									manifestMap["ID"] = id
								}
								
								manifestBytes, _ := json.Marshal(manifestMap)
								var manifest Manifest
								if err := json.Unmarshal(manifestBytes, &manifest); err == nil {
									// If ID is still empty, try to use MongoID
									if manifest.ID == "" && manifest.MongoID != "" {
										manifest.ID = manifest.MongoID
									}
									if s.client.debug {
										fmt.Printf("Parsed manifest: %+v\n", manifest)
									}
									manifests = append(manifests, &manifest)
								} else {
									if s.client.debug {
										fmt.Printf("Error unmarshaling manifest: %v\n", err)
									}
								}
							}
						}
					}
				}
			}
		}
		
		return manifests, nil
	} else {
		// Standard endpoint returns an array of manifests
		var resp struct {
			Data []*Manifest `json:"data"`
		}
		_, err = s.client.do(req, &resp)
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
}

// Get returns a single manifest by ID
func (s *ManifestsService) Get(id string, opts *ManifestListOptions) (*Manifest, error) {
	path := "manifests"
	
	// Use specific VCS endpoint if provided
	if opts != nil && opts.VCS != "" {
		vcsType := strings.ToLower(string(opts.VCS))
		if vcsType == "github" {
			path = "githubManifest"
		} else if vcsType == "gitlab" {
			path = "gitlabManifest"
		}
	}
	
	path = fmt.Sprintf("%s/%s", path, id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// For GitHub/GitLab endpoints, the response structure is different
	if opts != nil && opts.VCS != "" {
		// Parse as a raw map to inspect the structure
		var rawResp map[string]interface{}
		_, err = s.client.do(req, &rawResp)
		if err != nil {
			return nil, err
		}

		// Extract manifest from the data field
		if dataField, ok := rawResp["data"]; ok {
			if dataObj, ok := dataField.(map[string]interface{}); ok {
				// Check for catalogManifest field first
				if catalogManifestField, ok := dataObj["catalogManifest"]; ok {
					if catalogMap, ok := catalogManifestField.(map[string]interface{}); ok {
						manifestBytes, _ := json.Marshal(catalogMap)
						var manifest Manifest
						if err := json.Unmarshal(manifestBytes, &manifest); err == nil {
							return &manifest, nil
						}
					}
				}
				
				// If no catalogManifest field, try other fields
				for _, field := range []string{"manifest", "item"} {
					if manifestField, ok := dataObj[field]; ok {
						if manifestMap, ok := manifestField.(map[string]interface{}); ok {
							manifestBytes, _ := json.Marshal(manifestMap)
							var manifest Manifest
							if err := json.Unmarshal(manifestBytes, &manifest); err == nil {
								return &manifest, nil
							}
						}
					}
				}
				
				// If no specific field found, try using the data object itself
				manifestBytes, _ := json.Marshal(dataObj)
				var manifest Manifest
				if err := json.Unmarshal(manifestBytes, &manifest); err == nil && manifest.ID != "" {
					return &manifest, nil
				}
			}
		}
		
		// If all attempts fail, return error
		return nil, fmt.Errorf("could not parse manifest response")
	} else {
		// Standard endpoint returns a manifest directly
		var resp struct {
			Data *Manifest `json:"data"`
		}
		_, err = s.client.do(req, &resp)
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
}
