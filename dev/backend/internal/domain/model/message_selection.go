package model

import "time"

type MessageSelection struct {
	UUID         string
	SelectedText string
	RangeStart   int
	RangeEnd     int
	CreatedAt    time.Time
}
