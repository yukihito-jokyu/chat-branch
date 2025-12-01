package model

import "time"

type Chat struct {
	UUID           string
	ProjectUUID    string
	ParentUUID     *string
	Title          string
	Status         string
	ContextSummary string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// usecase から呼び出されるのでここに配置する
type ForkPreviewRequest struct {
	TargetMessageUUID string `json:"target_message_uuid"`
	SelectedText      string `json:"selected_text"`
	RangeStart        int    `json:"range_start"`
	RangeEnd          int    `json:"range_end"`
}

type ForkPreviewResponse struct {
	SuggestedTitle   string `json:"suggested_title"`
	GeneratedContext string `json:"generated_context"`
}
