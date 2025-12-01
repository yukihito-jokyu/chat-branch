package model

import "time"

type Chat struct {
	UUID        string
	ProjectUUID string
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
