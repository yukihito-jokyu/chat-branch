package model

import "time"

type Project struct {
	ID        string    `json:"uuid"`
	UserID    string    `json:"-"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}
