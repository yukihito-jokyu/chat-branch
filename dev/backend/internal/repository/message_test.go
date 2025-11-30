package repository

import (
	"backend/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestMessageRepository_Create(t *testing.T) {
	type args struct {
		message *model.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正常系: メッセージが作成できること",
			args: args{
				message: &model.Message{
					UUID:      "message-uuid",
					ChatUUID:  "chat-uuid",
					Role:      "user",
					Content:   "hello",
					UserUUID:  "user-uuid",
					CreatedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: UUIDが重複している場合エラーになること",
			args: args{
				message: &model.Message{
					UUID:      "duplicate-uuid",
					ChatUUID:  "chat-uuid",
					Role:      "user",
					Content:   "duplicate message",
					UserUUID:  "user-uuid",
					CreatedAt: time.Now(),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// インメモリDBのセットアップ
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			// マイグレーション
			if err := db.AutoMigrate(&messageORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			r := NewMessageRepository(db)

			// 重複エラーのテストケースのための事前データ投入
			if tt.name == "異常系: UUIDが重複している場合エラーになること" {
				if err := r.Create(context.Background(), tt.args.message); err != nil {
					t.Fatalf("failed to create pre-existing message: %v", err)
				}
			}

			if err := r.Create(context.Background(), tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
