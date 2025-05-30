package enbuild

type Stack struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name"`
	// Stack      StackName      `json:"stack"`
	Catalog    CatalogInfo `json:"catalog"`
	Components []Component `json:"components"`
	Status     string      `json:"status"`
	Type       string      `json:"type"`
	CreatedBy  string      `json:"createdBy"`
	UpdatedBy  string      `json:"updatedBy"`
	CreatedOn  string      `json:"createdOn,omitempty"`
	Created_on string      `json:"created_on,omitempty"`
	UpdatedOn  string      `json:"updatedOn,omitempty"`
	// V          int            `json:"__v,omitempty"`
	Logs []StackLog `json:"logs,omitempty"`
	// Project    *ProjectInfo   `json:"project,omitempty"`
	// Pipeline   []PipelineItem `json:"pipeline,omitempty"`
	// Pending    int            `json:"pending,omitempty"`
	// PermissionsSlug map[string]interface{} `json:"permissionsSlug,omitempty"`
	// Permissions     map[string]interface{} `json:"permissions,omitempty"`
}

type StackName struct {
	Name string `json:"name"`
}

type CatalogInfo struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Type string `json:"type"`
	Name string `json:"name"`
	// Data   CatalogData `json:"data"`
}

type CatalogData struct {
	ID         string `json:"_id"`
	Type       string `json:"type"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Repository string `json:"repository"`
	// ProjectID         string              `json:"project_id"`
	ReadmeFilePath    string `json:"readme_file_path"`
	ValuesFolderPath  string `json:"values_folder_path"`
	SecretsFolderPath string `json:"secrets_folder_path"`
	Ref               string `json:"ref"`
	// Sops              bool                `json:"sops,omitempty"`
	ImagePath      string              `json:"image_path"`
	MultiSelect    bool                `json:"multi_select"`
	VCS            string              `json:"vcs"`
	Status         string              `json:"status"`
	CreatedBy      string              `json:"created_by"`
	UpdatedBy      string              `json:"updatedBy"`
	Order          int                 `json:"order"`
	Description    string              `json:"description"`
	RepositoryId   string              `json:"repositoryId"`
	Components     []ComponentCfg      `json:"components"`
	Infrastructure InfraData           `json:"infrastructure"`
	Configuration  []ConfigurationItem `json:"configuration"`
	// CreatedOn         string        `json:"createdOn,omitempty"`
	UpdatedOn    string `json:"updatedOn,omitempty"`
	DownloadYAML bool   `json:"download_yaml,omitempty"`
}

type ConfigurationItem struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Fields []ConfigField `json:"fields"`
}

type ConfigField struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Variable     string `json:"variable,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Type         string `json:"type"`
}

type ComponentCfg struct {
	Type string `json:"type"`
	// ProjectID        string `json:"project_id"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	ToolType         string `json:"tool_type"`
	VariableFilePath string `json:"variable_file_path"`
	ImagePath        string `json:"image_path"`
	Mandatory        bool   `json:"mandatory"`
	RepositoryId     string `json:"repositoryId"`
	Ref              string `json:"ref"`
	Repository       string `json:"repository"`
	ID               int    `json:"id"`
	CatalogID        int    `json:"catalog_id"`
}

type InfraData struct {
	Slug           string        `json:"slug"`
	ShowKubeConfig bool          `json:"showKubeConfig"`
	Selections     []InfraSelect `json:"selections"`
	CatalogID      int           `json:"catalog_id"`
}

type InfraSelect struct {
	Slug     string       `json:"slug"`
	Name     string       `json:"name"`
	Selected bool         `json:"selected"`
	Fields   []InfraField `json:"fields"`
}

type InfraField struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Key       string `json:"key,omitempty"`
	Variable  string `json:"variable,omitempty"`
	Required  bool   `json:"required,omitempty"`
	Plaintext bool   `json:"plaintext"`
}

type Component struct {
	Name string          `json:"name"`
	Data []ComponentData `json:"data"`
}

type ComponentData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StackLog struct {
	Slug   string      `json:"slug"`
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Logs   []LogDetail `json:"logs"`
}

type LogDetail struct {
	Timestamp int64  `json:"timestamp"`
	Msg       string `json:"msg"`
	Level     string `json:"level"`
}

type ProjectInfo struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type PipelineItem struct {
	ID        int    `json:"id"`
	IID       int    `json:"iid"`
	ProjectID int    `json:"project_id"`
	SHA       string `json:"sha"`
	Ref       string `json:"ref"`
	Status    string `json:"status"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	WebURL    string `json:"web_url"`
	Name      string `json:"name"`
}
