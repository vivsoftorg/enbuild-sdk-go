package types

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
	VCS         string                 `json:"vcs,omitempty"`
	Slug        string                 `json:"slug,omitempty"`
	Type        string                 `json:"type,omitempty"`
}
