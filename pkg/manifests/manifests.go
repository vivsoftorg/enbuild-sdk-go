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
	return &Service{
		client: client,
	}
}

// ManifestListOptions specifies the optional parameters to the ManifestsService.List method
type ManifestListOptions struct {
	VCS types.VCSType `url:"vcs,omitempty"`
}

// List returns a list of manifests
func (s *Service) List(opts *ManifestListOptions) ([]*types.Manifest, error) {
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
	
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// For GitHub/GitLab endpoints, the response structure is different
	if opts != nil && opts.VCS != "" {
		// Parse as a raw map to inspect the structure
		var rawResp map[string]interface{}
		_, err = s.client.Do(req, &rawResp)
		if err != nil {
			return nil, err
		}

		if s.client.Debug {
			fmt.Printf("Raw response from API: %+v\n", rawResp)
		}

		// Extract manifests from the catalogManifest field in the data
		var manifests []*types.Manifest
		
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
								var manifest types.Manifest
								if err := json.Unmarshal(manifestBytes, &manifest); err == nil {
									// If ID is still empty, try to use MongoID
									if manifest.ID == "" && manifest.MongoID != "" {
										manifest.ID = manifest.MongoID
									}
									if s.client.Debug {
										fmt.Printf("Parsed manifest: %+v\n", manifest)
									}
									manifests = append(manifests, &manifest)
								} else {
									if s.client.Debug {
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
			Data []*types.Manifest `json:"data"`
		}
		_, err = s.client.Do(req, &resp)
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
}
