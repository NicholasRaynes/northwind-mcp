package models

type EndpointSchema struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	Parameters  map[string]string `json:"parameters"`
	Returns     map[string]string `json:"returns"`
}
