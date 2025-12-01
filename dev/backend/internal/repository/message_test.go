package repository

import (
	"backend/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type messageSelectionORM struct {
	UUID         string `gorm:"primaryKey;column:uuid;size:255"`
	SelectedText string `gorm:"column:selected_text;type:text"`
	RangeStart   int    `gorm:"column:range_start"`
	RangeEnd     int    `gorm:"column:range_end"`
	CreatedID    string `gorm:"column:created_id;size:255"`
}

func (messageSelectionORM) TableName() string {
	return "message_selections"
}

func strPtr(s string) *string {
	return &s
}

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
					CreatedAt: time.Now(),
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

func TestMessageRepository_FindMessagesByChatID(t *testing.T) {
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "正常系: チャットIDに紐づくメッセージが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&messageORM{
					UUID:      "msg-1",
					ChatUUID:  "chat-uuid",
					Role:      "user",
					Content:   "hello",
					CreatedAt: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				})
				db.Create(&messageORM{
					UUID:      "msg-2",
					ChatUUID:  "chat-uuid",
					Role:      "assistant",
					Content:   "hi",
					CreatedAt: time.Date(2023, 1, 1, 10, 0, 1, 0, time.UTC),
				})
				db.Create(&messageORM{
					UUID:      "msg-other",
					ChatUUID:  "other-chat-uuid",
					Role:      "user",
					Content:   "other",
					CreatedAt: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				})
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "正常系: メッセージが存在しない場合は空のリストが返ること",
			args: args{
				chatUUID: "empty-chat-uuid",
			},
			setupData: func(db *gorm.DB) {},
			wantLen:   0,
			wantErr:   false,
		},
		{
			name: "正常系: フォーク情報を含むメッセージが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				// メッセージの作成
				db.Create(&messageORM{
					UUID:      "msg-1",
					ChatUUID:  "chat-uuid",
					Role:      "assistant",
					Content:   "parent message",
					CreatedAt: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				})

				// message_selections テーブルの作成とデータ投入

				if err := db.AutoMigrate(&messageSelectionORM{}); err != nil {
					panic(err)
				}
				db.Create(&messageSelectionORM{
					UUID:         "selection-1",
					SelectedText: "selected",
					RangeStart:   0,
					RangeEnd:     5,
					CreatedID:    "test",
				})

				// chats テーブルの作成とデータ投入 (フォーク用)
				if err := db.AutoMigrate(&chatORM{}); err != nil {
					panic(err)
				}
				db.Create(&chatORM{
					UUID:                 "child-chat-1",
					ProjectUUID:          "project-1",
					SourceMessageUUID:    strPtr("msg-1"),
					MessageSelectionUUID: strPtr("selection-1"),
					Title:                "child chat",
					Status:               "open",
					CreatedID:            "test",
				})
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "異常系: フォーク取得時にDBエラーが発生した場合エラーになること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				// メッセージを作成して、最初のクエリが成功するようにする
				db.Create(&messageORM{
					UUID:      "msg-1",
					ChatUUID:  "chat-uuid",
					Role:      "user",
					Content:   "hello",
					CreatedAt: time.Now(),
				})

				// 2番目のクエリが失敗するようにテーブルを削除する
				if err := db.Migrator().DropTable(&messageSelectionORM{}); err != nil {
					panic(err)
				}
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "異常系: DBエラーが発生した場合エラーになること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			wantLen: 0,
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
			if err := db.AutoMigrate(&messageORM{}, &chatORM{}, &messageSelectionORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			// データ投入
			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewMessageRepository(db)
			got, err := r.FindMessagesByChatID(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.FindMessagesByChatID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("messageRepository.FindMessagesByChatID() length = %v, want %v", len(got), tt.wantLen)
			}
			// 順序の確認 (CreatedID順)
			if len(got) > 1 {
				if got[0].CreatedAt.After(got[1].CreatedAt) {
					t.Errorf("messageRepository.FindMessagesByChatID() order is wrong")
				}
			}
		})
	}
}
