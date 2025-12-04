package repository

import (
	"backend/internal/domain/model"
	"context"
)

type EdgeRepository interface {
	FindEdgesByChatID(ctx context.Context, chatUUID string) ([]*model.Edge, error)
	Create(ctx context.Context, edge *model.Edge) error
}
