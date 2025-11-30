package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type messageORM struct {
	UUID      string    `gorm:"primaryKey;column:uuid;size:36"`
	ChatUUID  string    `gorm:"column:chat_uuid;size:36"`
	Role      string    `gorm:"size:50"`
	Content   string    `gorm:"type:text"`
	CreatedID string    `gorm:"column:created_id;size:255"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (messageORM) TableName() string {
	return "messages"
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *model.Message) error {
	slog.DebugContext(ctx, "メッセージ作成処理を開始", "message_uuid", message.UUID)
	orm := messageORM{
		UUID:      message.UUID,
		ChatUUID:  message.ChatUUID,
		Role:      message.Role,
		Content:   message.Content,
		CreatedID: message.UserUUID,
		CreatedAt: message.CreatedAt,
	}
	db := getDB(ctx, r.db)
	return db.WithContext(ctx).Create(&orm).Error
}
