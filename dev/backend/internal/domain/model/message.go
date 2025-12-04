package model

import "time"

type Fork struct {
	ChatUUID     string
	SelectedText string
	RangeStart   int
	RangeEnd     int
}

type Message struct {
	UUID              string
	ChatUUID          string
	ParentMessageUUID *string // 追加
	Role              string  // user or assistant
	Content           string
	ContextSummary    *string
	SourceChatUUID    *string
	PositionX         float64
	PositionY         float64
	Forks             []Fork
	MergeReports      []*Message
	CreatedAt         time.Time
}
