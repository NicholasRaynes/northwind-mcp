package models

type APISchema struct {
	OpenAPI string             `json:"openapi"`
	Info    APIInfo            `json:"info"`
	Paths   map[string]APIPath `json:"paths"`
}

type APIInfo struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type APIPath struct {
	Get *APIEndpoint `json:"get,omitempty"`
}

type APIEndpoint struct {
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Parameters  []APIParameter         `json:"parameters,omitempty"`
	Responses   map[string]APIResponse `json:"responses"`
}

type APIParameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	SchemaType  string `json:"type"`
}

type APIResponse struct {
	Description string `json:"description"`
}
