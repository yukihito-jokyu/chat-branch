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
	UUID        string    `gorm:"primaryKey;column:uuid;size:36"`
	ProjectUUID string    `gorm:"column:project_uuid;size:36"`
	Title       string    `gorm:"column:title;size:255"`
	CreatedID   string    `gorm:"column:created_id;size:255"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
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
	orm := chatORM{
		UUID:        chat.UUID,
		ProjectUUID: chat.ProjectUUID,
		Title:       chat.Title,
		CreatedID:   uuid.New().String(),
		CreatedAt:   chat.CreatedAt,
		UpdatedAt:   chat.UpdatedAt,
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

	return &model.Chat{
		UUID:        orm.UUID,
		ProjectUUID: orm.ProjectUUID,
		Title:       orm.Title,
		CreatedAt:   orm.CreatedAt,
		UpdatedAt:   orm.UpdatedAt,
	}, nil
}
