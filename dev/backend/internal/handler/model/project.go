package model

import "time"

type ProjectResponse struct {
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	InitialMessage string `json:"initial_message"`
}

type CreateProjectResponse struct {
	ProjectUUID string      `json:"project_uuid"`
	ChatUUID    string      `json:"chat_uuid"`
	MessageInfo MessageInfo `json:"message_info"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type MessageInfo struct {
	MessageUUID string `json:"message_uuid"`
	Message     string `json:"message"`
}
