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

// メッセージを保存する
func (r *messageRepository) Create(ctx context.Context, message *model.Message) error {
	slog.DebugContext(ctx, "メッセージ作成処理を開始", "message_uuid", message.UUID)
	orm := messageORM{
		UUID:      message.UUID,
		ChatUUID:  message.ChatUUID,
		Role:      message.Role,
		Content:   message.Content,
		CreatedID: uuid.New().String(),
		CreatedAt: message.CreatedAt,
	}
	db := getDB(ctx, r.db)
	return db.WithContext(ctx).Create(&orm).Error
}

// 指定されたチャットIDのメッセージを取得する
func (r *messageRepository) FindMessagesByChatID(ctx context.Context, chatUUID string) ([]*model.Message, error) {
	slog.DebugContext(ctx, "メッセージ取得処理を開始", "chat_uuid", chatUUID)
	var orms []messageORM
	db := getDB(ctx, r.db)
	if err := db.WithContext(ctx).Where("chat_uuid = ?", chatUUID).Order("created_at asc").Find(&orms).Error; err != nil {
		return nil, err
	}

	var messages []*model.Message
	for _, orm := range orms {
		messages = append(messages, &model.Message{
			UUID:      orm.UUID,
			ChatUUID:  orm.ChatUUID,
			Role:      orm.Role,
			Content:   orm.Content,
			CreatedAt: orm.CreatedAt,
		})
	}
	return messages, nil
}
