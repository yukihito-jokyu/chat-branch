package model

import "time"

type Message struct {
	UUID      string
	ChatUUID  string
	Role      string // user or assistant
	Content   string
	CreatedAt time.Time
}
