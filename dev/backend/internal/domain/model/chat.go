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
