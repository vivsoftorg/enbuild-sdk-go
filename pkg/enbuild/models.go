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
// fot Stacks
type Stack struct {
    ID         string         `json:"id"`
    Name       string         `json:"name"`
    Stack      StackName      `json:"stack"`
    Catalog    CatalogInfo    `json:"catalog"`
    Components []Component    `json:"components"`
    Status     string         `json:"status"`
    Data       StackData      `json:"data"`
    Infrastructure Infra      `json:"infrastructure"`
    Type       string         `json:"type"`
    CreatedBy  string         `json:"createdBy"`
    UpdatedBy  string         `json:"updatedBy"`
    CreatedOn  interface{}          `json:"createdOn,omitempty"`
    Created_on interface{}          `json:"created_on,omitempty"`
    UpdatedOn  interface{}          `json:"updatedOn,omitempty"`
    V          int            `json:"__v"`
}

type StackName struct {
    Name string `json:"name"`
}

type CatalogInfo struct {
    ID     int         `json:"id"`
    Slug   string      `json:"slug"`
    Type   string      `json:"type"`
    Name   string      `json:"name"`
    Data   CatalogData `json:"data"`
}

type CatalogData struct {
    ID                string         `json:"_id"`
    Type              string         `json:"type"`
    Slug              string         `json:"slug"`
    Name              string         `json:"name"`
    Repository        string         `json:"repository"`
    ProjectID         string         `json:"project_id"`
    ReadmeFilePath    string         `json:"readme_file_path"`
    ValuesFolderPath  string         `json:"values_folder_path"`
    Ref               string         `json:"ref"`
    Sops              bool           `json:"sops,omitempty"`
    ImagePath         string         `json:"image_path"`
    MultiSelect       bool           `json:"multi_select"`
    VCS               string         `json:"vcs"`
    Status            string         `json:"status"`
    CreatedBy         string         `json:"created_by"`
    UpdatedBy         string         `json:"updatedBy"`
    Order             int            `json:"order"`
    Description       string         `json:"description"`
    RepositoryId      string         `json:"repositoryId"`
    Components        []ComponentCfg `json:"components"`
    Infrastructure    InfraData      `json:"infrastructure"`
    Configuration     interface{}    `json:"configuration"`
    CreatedBy2        string         `json:"createdBy"` // Note: duplicate json tag "createdBy"
    CreatedOn         interface{}          `json:"createdOn,omitempty"`
    Token             string         `json:"token"`
    ID2               int            `json:"id"` // Note: duplicate json tag "id"
    UpdatedOn         interface{}          `json:"updatedOn,omitempty"`
}

type ComponentCfg struct {
    Type            string `json:"type"`
    ProjectID       string `json:"project_id"`
    Name            string `json:"name"`
    Slug            string `json:"slug"`
    ToolType        string `json:"tool_type"`
    VariableFilePath string `json:"variable_file_path"`
    ImagePath       string `json:"image_path"`
    Mandatory       bool   `json:"mandatory"`
    RepositoryId    string `json:"repositoryId"`
    Ref             string `json:"ref"`
    Repository      string `json:"repository"`
    ID              int    `json:"id"`
    CatalogID       int    `json:"catalog_id"`
}

type InfraData struct {
    Slug         string         `json:"slug"`
    ShowKubeConfig bool         `json:"showKubeConfig"`
    Selections   []InfraSelect  `json:"selections"`
    CatalogID    int            `json:"catalog_id"`
}

type InfraSelect struct {
    Slug    string         `json:"slug"`
    Name    string         `json:"name"`
    Selected bool          `json:"selected"`
    Fields  []InfraField   `json:"fields"`
}

type InfraField struct {
    Name      string `json:"name"`
    Type      string `json:"type"`
    Key       string `json:"key,omitempty"`
    Variable  string `json:"variable,omitempty"`
    Plaintext bool   `json:"plaintext"`
}

type Component struct {
    Name string         `json:"name"`
    Data []ComponentData `json:"data"`
}

type ComponentData struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type StackData struct {
    StackName        string            `json:"stackName"`
    ProjectTemplateID string           `json:"projectTemplateId"`
    // EnvVariables     []EnvVar          `json:"envVariables"`
    ComponentConfigs []ComponentCfg2   `json:"componentConfigs"`
    RegistryConfig   RegistryConfig    `json:"registryConfig"`
    RepoConfig       RepoConfig        `json:"repoConfig"`
    User             UserInfo          `json:"user"`
    Assignee         string            `json:"assignee"`
}

type EnvVar struct {
    Key   string `json:"key"`
    Value string `json:"value"`
    Type  string `json:"type"`
}

type ComponentCfg2 struct {
    Type            string `json:"type"`
    ProjectID       string `json:"project_id"`
    Name            string `json:"name"`
    Slug            string `json:"slug"`
    ToolType        string `json:"tool_type"`
    VariableFilePath string `json:"variable_file_path"`
    ImagePath       string `json:"image_path"`
    Mandatory       bool   `json:"mandatory"`
    RepositoryId    string `json:"repositoryId"`
    Ref             string `json:"ref"`
    Repository      string `json:"repository"`
    ID              int    `json:"id"`
    CatalogID       int    `json:"catalog_id"`
    Code            string `json:"code"`
    OriginalCode    string `json:"originalCode"`
    Content         string `json:"content"`
    Path            string `json:"path"`
    Secrets         string `json:"secrets"`
    NoChange        bool   `json:"noChange"`
}

type RegistryConfig struct {
    Path   string `json:"path"`
    Config string `json:"config"`
}

type RepoConfig struct {
    Path   string `json:"path"`
    Config string `json:"config"`
}

type UserInfo struct {
    Email       string      `json:"email"`
    RealmAccess RealmAccess `json:"realmAccess"`
}

type RealmAccess struct {
    Roles []string `json:"roles"`
}

type Infra struct {
    Name   string                 `json:"name"`
    Config map[string]interface{} `json:"config"`
}