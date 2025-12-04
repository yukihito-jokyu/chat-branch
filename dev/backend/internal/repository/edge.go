package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type edgeORM struct {
	UUID              string `gorm:"primaryKey;column:uuid;size:255"`
	ChatUUID          string `gorm:"column:chat_uuid;size:255"`
	SourceMessageUUID string `gorm:"column:source_message_uuid;size:255"`
	TargetMessageUUID string `gorm:"column:target_message_uuid;size:255"`
}

func (edgeORM) TableName() string {
	return "edges"
}

type edgeRepository struct {
	db *gorm.DB
}

func NewEdgeRepository(db *gorm.DB) repository.EdgeRepository {
	return &edgeRepository{db: db}
}

// 指定されたチャットIDのエッジを取得する
func (r *edgeRepository) FindEdgesByChatID(ctx context.Context, chatUUID string) ([]*model.Edge, error) {
	slog.DebugContext(ctx, "エッジ取得処理を開始", "chat_uuid", chatUUID)
	var orms []edgeORM
	db := getDB(ctx, r.db)
	if err := db.WithContext(ctx).Where("chat_uuid = ?", chatUUID).Find(&orms).Error; err != nil {
		return nil, err
	}

	var edges []*model.Edge
	for _, orm := range orms {
		edges = append(edges, &model.Edge{
			UUID:              orm.UUID,
			ChatUUID:          orm.ChatUUID,
			SourceMessageUUID: orm.SourceMessageUUID,
			TargetMessageUUID: orm.TargetMessageUUID,
		})
	}
	return edges, nil
}

// エッジを作成する
func (r *edgeRepository) Create(ctx context.Context, edge *model.Edge) error {
	slog.DebugContext(ctx, "エッジ作成処理を開始", "edge_uuid", edge.UUID)
	orm := edgeORM{
		UUID:              edge.UUID,
		ChatUUID:          edge.ChatUUID,
		SourceMessageUUID: edge.SourceMessageUUID,
		TargetMessageUUID: edge.TargetMessageUUID,
	}
	db := getDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(&orm).Error; err != nil {
		return err
	}
	return nil
}
