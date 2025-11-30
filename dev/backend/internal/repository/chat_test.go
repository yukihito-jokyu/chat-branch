package repository

import (
	"backend/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestChatRepository_Create(t *testing.T) {
	type args struct {
		chat *model.Chat
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正常系: チャットが作成できること",
			args: args{
				chat: &model.Chat{
					UUID:        "chat-uuid",
					ProjectUUID: "project-uuid",
					Title:       "test chat",
					UserUUID:    "user-uuid",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: UUIDが重複している場合エラーになること",
			args: args{
				chat: &model.Chat{
					UUID:        "duplicate-uuid",
					ProjectUUID: "project-uuid",
					Title:       "duplicate chat",
					UserUUID:    "user-uuid",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
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
			if err := db.AutoMigrate(&chatORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			r := NewChatRepository(db)

			// 重複エラーのテストケースのための事前データ投入
			if tt.name == "異常系: UUIDが重複している場合エラーになること" {
				if err := r.Create(context.Background(), tt.args.chat); err != nil {
					t.Fatalf("failed to create pre-existing chat: %v", err)
				}
			}

			if err := r.Create(context.Background(), tt.args.chat); (err != nil) != tt.wantErr {
				t.Errorf("chatRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
