package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type chatORM struct {
	UUID                 string    `gorm:"primaryKey;column:uuid;size:255"`
	ProjectUUID          string    `gorm:"column:project_uuid;size:255"`
	ParentChatUUID       *string   `gorm:"column:parent_chat_uuid;size:255"`
	SourceMessageUUID    *string   `gorm:"column:source_message_uuid;size:255"`
	MessageSelectionUUID *string   `gorm:"column:message_selection_uuid;size:255"`
	Title                string    `gorm:"column:title;size:255"`
	Status               string    `gorm:"column:status;size:50"`
	ContextSummary       *string   `gorm:"column:context_summary;type:text"`
	PositionX            float64   `gorm:"column:position_x"`
	PositionY            float64   `gorm:"column:position_y"`
	CreatedID            string    `gorm:"column:created_id;size:255"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
	UpdatedID            *string   `gorm:"column:updated_id;size:255"`
}

func (chatORM) TableName() string {
	return "chats"
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepository{db: db}
}

// チャットを保存する
func (r *chatRepository) Create(ctx context.Context, chat *model.Chat) error {
	slog.DebugContext(ctx, "チャット作成処理を開始", "chat_uuid", chat.UUID)

	var contextSummary *string
	if chat.ContextSummary != "" {
		contextSummary = &chat.ContextSummary
	}

	orm := chatORM{
		UUID:                 chat.UUID,
		ProjectUUID:          chat.ProjectUUID,
		ParentChatUUID:       chat.ParentUUID,
		SourceMessageUUID:    chat.SourceMessageUUID,
		MessageSelectionUUID: chat.MessageSelectionUUID,
		Title:                chat.Title,
		Status:               chat.Status,
		ContextSummary:       contextSummary,
		CreatedID:            uuid.New().String(),
		CreatedAt:            chat.CreatedAt,
		UpdatedAt:            chat.UpdatedAt,
	}
	db := getDB(ctx, r.db)
	return db.WithContext(ctx).Create(&orm).Error
}

// 指定されたチャットIDのチャットを取得する
func (r *chatRepository) FindByID(ctx context.Context, uuid string) (*model.Chat, error) {
	slog.DebugContext(ctx, "チャット取得処理を開始", "chat_uuid", uuid)
	var orm chatORM
	db := getDB(ctx, r.db)
	if err := db.WithContext(ctx).Where("uuid = ?", uuid).First(&orm).Error; err != nil {
		return nil, err
	}

	var contextSummary string
	if orm.ContextSummary != nil {
		contextSummary = *orm.ContextSummary
	}

	return &model.Chat{
		UUID:                 orm.UUID,
		ProjectUUID:          orm.ProjectUUID,
		ParentUUID:           orm.ParentChatUUID,
		SourceMessageUUID:    orm.SourceMessageUUID,
		MessageSelectionUUID: orm.MessageSelectionUUID,
		Title:                orm.Title,
		Status:               orm.Status,
		ContextSummary:       contextSummary,
		CreatedAt:            orm.CreatedAt,
		UpdatedAt:            orm.UpdatedAt,
	}, nil
}
