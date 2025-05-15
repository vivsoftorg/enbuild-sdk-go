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
	VCS string `url:"vcs,omitempty"`
}

// List returns a list of manifests
func (s *Service) List(opts *ManifestListOptions) ([]*types.Manifest, error) {
	path := "manifests"
	if vcsPath := getVCSPath(opts); vcsPath != "" {
		path = vcsPath
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if opts == nil || opts.VCS == "" {
		return s.getStandardManifests(req)
	}
	return s.getCustomEndpointManifests(req)
}

func getVCSPath(opts *ManifestListOptions) string {
	if opts == nil || opts.VCS == "" {
		return ""
	}
	switch strings.ToLower(string(opts.VCS)) {
	case "github":
		return "githubManifest"
	case "gitlab":
		return "gitlabManifest"
	default:
		return ""
	}
}

func (s *Service) getStandardManifests(req *http.Request) ([]*types.Manifest, error) {
	var resp struct {
		Data []*types.Manifest `json:"data"`
	}
	if _, err := s.client.Do(req, &resp); err != nil {
		return nil, err
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
