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
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

func TestChatRepository_FindByID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		want      *model.Chat
		wantErr   bool
	}{
		{
			name: "正常系: 存在するチャットが取得できること",
			args: args{
				uuid: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&chatORM{
					UUID:        "chat-uuid",
					ProjectUUID: "project-uuid",
					Title:       "test chat",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				})
			},
			want: &model.Chat{
				UUID:        "chat-uuid",
				ProjectUUID: "project-uuid",
				Title:       "test chat",
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "異常系: 存在しないチャットの場合エラーになること",
			args: args{
				uuid: "non-existent-uuid",
			},
			setupData: func(db *gorm.DB) {},
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// インメモリDBのセットアップ
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			// マイグレーション
			if err := db.AutoMigrate(&chatORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			// データ投入
			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewChatRepository(db)
			got, err := r.FindByID(context.Background(), tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 時間の比較は厳密に行うのが難しいため、EqualValuesなどで比較するか、各フィールドを比較する
				// ここでは assert を使っていないので手動比較
				if got.UUID != tt.want.UUID || got.ProjectUUID != tt.want.ProjectUUID || got.Title != tt.want.Title {
					t.Errorf("chatRepository.FindByID() = %v, want %v", got, tt.want)
				}
				if !got.CreatedAt.Equal(tt.want.CreatedAt) {
					t.Errorf("chatRepository.FindByID() CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
				}
			}
		})
	}
}
