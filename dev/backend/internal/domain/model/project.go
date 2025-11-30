package model

import "time"

type Project struct {
	UUID      string
	UserUUID  string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
