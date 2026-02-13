package models

type APIResponse struct {
	Products    map[string]string `json:"products"`
	Categories  map[string]string `json:"categories"`
	Environment string            `json:"environment"`
	Message     string            `json:"message"`
	Version     string            `json:"version"`
}
