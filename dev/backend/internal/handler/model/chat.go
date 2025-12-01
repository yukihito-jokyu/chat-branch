package model

type GetChatResponse struct {
	UUID           string  `json:"uuid"`
	ProjectUUID    string  `json:"project_uuid"`
	ParentUUID     *string `json:"parent_uuid"`
	Title          string  `json:"title"`
	Status         string  `json:"status"`
	ContextSummary string  `json:"context_summary"`
}
