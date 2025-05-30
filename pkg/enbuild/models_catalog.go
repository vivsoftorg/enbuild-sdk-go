package enbuild

// Catalog represents an ENBUILD catalog
type Catalog struct {
	ID          interface{}            `json:"_id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	Version     string                 `json:"version,omitempty"`
	CreatedOn   interface{}            `json:"createdOn,omitempty"`
	UpdatedOn   interface{}            `json:"updatedOn,omitempty"`
	VCS         string                 `json:"vcs,omitempty"`
	Slug        string                 `json:"slug,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

// CatalogListOptions specifies the optional parameters to the catalogsService.List method.
type CatalogListOptions struct {
	ID          string `url:"id,omitempty"`
	VCS         string `url:"vcs,omitempty"`
	Type        string `url:"type,omitempty"`
	Slug        string `url:"slug,omitempty"`
	Name        string `url:"name,omitempty"`
	Description string `url:"description,omitempty"`
	Version     string `url:"version,omitempty"`
}
