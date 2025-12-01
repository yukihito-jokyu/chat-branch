package model

import "time"

type Fork struct {
	ChatUUID     string
	SelectedText string
	RangeStart   int
	RangeEnd     int
}

type Message struct {
	UUID           string
	ChatUUID       string
	Role           string // user or assistant
	Content        string
	ContextSummary *string
	SourceChatUUID *string
	Forks          []Fork
	CreatedAt      time.Time
}
