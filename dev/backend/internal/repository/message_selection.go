package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageSelectionRepository struct {
	db *gorm.DB
}

func NewMessageSelectionRepository(db *gorm.DB) repository.MessageSelectionRepository {
	return &messageSelectionRepository{db: db}
}

type messageSelectionORM struct {
	UUID         string    `gorm:"primaryKey;column:uuid;size:36"`
	SelectedText string    `gorm:"column:selected_text;type:text"`
	RangeStart   int       `gorm:"column:range_start"`
	RangeEnd     int       `gorm:"column:range_end"`
	CreatedID    string    `gorm:"column:created_id;size:255"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	UpdatedID    *string   `gorm:"column:updated_id;size:255"`
}

func (messageSelectionORM) TableName() string {
	return "message_selections"
}

// メッセージ選択を保存する
func (r *messageSelectionRepository) Create(ctx context.Context, selection *model.MessageSelection) error {
	orm := messageSelectionORM{
		UUID:         selection.UUID,
		SelectedText: selection.SelectedText,
		RangeStart:   selection.RangeStart,
		RangeEnd:     selection.RangeEnd,
		CreatedID:    uuid.New().String(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	db := getDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(&orm).Error; err != nil {
		return fmt.Errorf("メッセージ選択の保存に失敗: %w", err)
	}

	return nil
}
