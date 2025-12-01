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
	UUID                 string    `gorm:"primaryKey;column:uuid;size:255"`
	ChatUUID             string    `gorm:"column:chat_uuid;size:255"`
	Role                 string    `gorm:"size:50"`
	Content              string    `gorm:"type:longtext"`
	ContextSummary       *string   `gorm:"column:context_summary;type:text"`
	SourceChatUUID       *string   `gorm:"column:source_chat_uuid;size:255"`
	MessageSelectionUUID *string   `gorm:"column:message_selection_uuid;size:255"`
	CreatedID            string    `gorm:"column:created_id;size:255"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
	UpdatedID            *string   `gorm:"column:updated_id;size:255"`
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

	if len(orms) == 0 {
		return []*model.Message{}, nil
	}

	messageUUIDs := make([]string, len(orms))
	for i, orm := range orms {
		messageUUIDs[i] = orm.UUID
	}

	type forkResult struct {
		ChatUUID          string
		SourceMessageUUID string
		SelectedText      string
		RangeStart        int
		RangeEnd          int
	}

	var forkResults []forkResult
	// chats テーブルと message_selections テーブルを結合して、指定されたメッセージUUIDのフォックスを取得する
	err := db.WithContext(ctx).Table("chats").
		Select("chats.uuid as chat_uuid, chats.source_message_uuid, ms.selected_text, ms.range_start, ms.range_end").
		Joins("JOIN message_selections ms ON chats.message_selection_uuid = ms.uuid").
		Where("chats.source_message_uuid IN ?", messageUUIDs).
		Scan(&forkResults).Error
	if err != nil {
		return nil, err
	}

	forksMap := make(map[string][]model.Fork)
	for _, res := range forkResults {
		forksMap[res.SourceMessageUUID] = append(forksMap[res.SourceMessageUUID], model.Fork{
			ChatUUID:     res.ChatUUID,
			SelectedText: res.SelectedText,
			RangeStart:   res.RangeStart,
			RangeEnd:     res.RangeEnd,
		})
	}

	var messages []*model.Message
	for _, orm := range orms {
		messages = append(messages, &model.Message{
			UUID:           orm.UUID,
			ChatUUID:       orm.ChatUUID,
			Role:           orm.Role,
			Content:        orm.Content,
			SourceChatUUID: orm.SourceChatUUID,
			Forks:          forksMap[orm.UUID],
			CreatedAt:      orm.CreatedAt,
		})
	}
	return messages, nil
}
