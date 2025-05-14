package enbuild

import (
	"net/http"
)

// AdminSettingsService handles communication with the admin settings related endpoints
type AdminSettingsService struct {
	client *Client
}

// AdminSettings represents ENBUILD admin settings
type AdminSettings struct {
	ID        string                 `json:"id,omitempty"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
	CreatedOn interface{}            `json:"createdOn,omitempty"`
	UpdatedOn interface{}            `json:"updatedOn,omitempty"`
	// Add other fields as needed
}

// Get returns all admin settings
func (s *AdminSettingsService) Get() (*AdminSettings, error) {
	req, err := s.client.newRequest(http.MethodGet, "adminSettings", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data *AdminSettings `json:"data"`
	}
	_, err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
