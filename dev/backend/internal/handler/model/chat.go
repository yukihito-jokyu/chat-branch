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

type ForkPreviewResponse struct {
	SuggestedTitle   string `json:"suggested_title"`
	GeneratedContext string `json:"generated_context"`
}

type MergePreviewResponse struct {
	SuggestedSummary string `json:"suggested_summary"`
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

type ForkChatRequest struct {
	TargetMessageUUID string `json:"target_message_uuid"`
	ParentChatUUID    string `json:"parent_chat_uuid"`
	SelectedText      string `json:"selected_text"`
	RangeStart        int    `json:"range_start"`
	RangeEnd          int    `json:"range_end"`
	Title             string `json:"title"`
	ContextSummary    string `json:"context_summary"`
}

type ForkChatResponse struct {
	NewChatID string `json:"new_chat_id"`
	Message   string `json:"message"`
}
