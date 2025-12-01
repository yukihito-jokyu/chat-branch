package model

type GetChatResponse struct {
	UUID           string  `json:"uuid"`
	ProjectUUID    string  `json:"project_uuid"`
	ParentUUID     *string `json:"parent_uuid"`
	Title          string  `json:"title"`
	Status         string  `json:"status"`
	ContextSummary string  `json:"context_summary"`
}

type MessageResponse struct {
	UUID           string         `json:"uuid"`
	Role           string         `json:"role"`
	Content        string         `json:"content"`
	Forks          []ForkResponse `json:"forks"`
	SourceChatUUID *string        `json:"source_chat_uuid,omitempty"`
}

type ForkResponse struct {
	ChatUUID     string `json:"chat_uuid"`
	SelectedText string `json:"selected_text"`
	RangeStart   int    `json:"range_start"`
	RangeEnd     int    `json:"range_end"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}
