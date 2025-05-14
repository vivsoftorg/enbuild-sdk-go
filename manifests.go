package enbuild

import (
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
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	Version     string                 `json:"version,omitempty"`
	CreatedOn   interface{}            `json:"createdOn,omitempty"`
	UpdatedOn   interface{}            `json:"updatedOn,omitempty"`
	VCS         VCSType                `json:"vcs,omitempty"`
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

	// Try different response formats for VCS-specific endpoints
	if opts != nil && opts.VCS != "" {
		// First try the nested manifests structure
		var vcsResp struct {
			Data struct {
				Manifests []*Manifest `json:"manifests"`
			} `json:"data"`
		}
		
		_, err = s.client.do(req, &vcsResp)
		if err == nil && len(vcsResp.Data.Manifests) > 0 {
			return vcsResp.Data.Manifests, nil
		}
		
		// If that fails, try direct array in data field
		req, err = s.client.newRequest(http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		
		var directResp struct {
			Data []*Manifest `json:"data"`
		}
		_, err = s.client.do(req, &directResp)
		if err == nil && len(directResp.Data) > 0 {
			return directResp.Data, nil
		}
		
		// Try with a different structure where manifests are directly in the response
		req, err = s.client.newRequest(http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		
		var flatResp []*Manifest
		_, err = s.client.do(req, &flatResp)
		if err == nil && len(flatResp) > 0 {
			return flatResp, nil
		}
		
		// If all attempts fail, return empty array rather than error
		return []*Manifest{}, nil
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

	// Try different response formats for VCS-specific endpoints
	if opts != nil && opts.VCS != "" {
		// First try the nested manifest structure
		var vcsResp struct {
			Data struct {
				Manifest *Manifest `json:"manifest"`
			} `json:"data"`
		}
		
		_, err = s.client.do(req, &vcsResp)
		if err == nil && vcsResp.Data.Manifest != nil {
			return vcsResp.Data.Manifest, nil
		}
		
		// If that fails, try direct object in data field
		req, err = s.client.newRequest(http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		
		var directResp struct {
			Data *Manifest `json:"data"`
		}
		_, err = s.client.do(req, &directResp)
		if err == nil && directResp.Data != nil {
			return directResp.Data, nil
		}
		
		// Try with a direct response
		req, err = s.client.newRequest(http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		
		var flatResp *Manifest
		_, err = s.client.do(req, &flatResp)
		if err == nil && flatResp != nil {
			return flatResp, nil
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
