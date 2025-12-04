package model

import "time"

type Project struct {
	UUID      string
	UserUUID  string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProjectTree struct {
	Nodes []ProjectNode
	Edges []ProjectEdge
}

type ProjectNode struct {
	ID       string
	ChatUUID string
	Data     ProjectNodeData
	Position ProjectNodePosition
}

type ProjectNodeData struct {
	UserMessage *string
	Assistant   string
}

type ProjectNodePosition struct {
	X float64
	Y float64
}

type ProjectEdge struct {
	ID     string
	Source string
	Target string
}
