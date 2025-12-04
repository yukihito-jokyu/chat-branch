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

type GetParentChatResponse struct {
	ChatUUID string `json:"chat_uuid"`
}

type GetProjectTreeResponse struct {
	Nodes []ProjectNode `json:"nodes"`
	Edges []ProjectEdge `json:"edges"`
}

type ProjectNode struct {
	ID       string              `json:"id"`
	ChatUUID string              `json:"chat_uuid"`
	Data     ProjectNodeData     `json:"data"`
	Position ProjectNodePosition `json:"position"`
}

type ProjectNodeData struct {
	UserMessage *string `json:"user_message"`
	Assistant   string  `json:"assistant"`
}

type ProjectNodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ProjectEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}
