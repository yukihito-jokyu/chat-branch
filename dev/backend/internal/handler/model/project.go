package model

import "time"

type ProjectResponse struct {
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}
