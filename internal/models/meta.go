package models

type APISchema struct {
	OpenAPI string             `json:"openapi"`
	Info    APIInfo            `json:"info"`
	Tags    []APITag           `json:"tags,omitempty"`
	Paths   map[string]APIPath `json:"paths"`
}

type APIInfo struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type APITag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type APIPath struct {
	Get *APIEndpoint `json:"get,omitempty"`
}

type APIEndpoint struct {
	Tags        []string               `json:"tags,omitempty"`
	OperationID string                 `json:"operationId,omitempty"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Parameters  []APIParameter         `json:"parameters,omitempty"`
	Responses   map[string]APIResponse `json:"responses"`
}

type APIParameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	SchemaType  string `json:"schemaType"`
}

type APIResponse struct {
	Description string `json:"description"`
}
