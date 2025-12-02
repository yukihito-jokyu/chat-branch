package repository

import (
	"backend/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

func TestMessageRepository_UpdateContextSummary(t *testing.T) {
	type args struct {
		messageUUID string
		summary     string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		wantErr   bool
		check     func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "正常系: コンテキストサマリが更新できること",
			args: args{
				messageUUID: "msg-1",
				summary:     "updated summary",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&messageORM{
					UUID:           "msg-1",
					ChatUUID:       "chat-uuid",
					Role:           "user",
					Content:        "hello",
					ContextSummary: nil,
					CreatedAt:      time.Now(),
				})
			},
			wantErr: false,
			check: func(t *testing.T, db *gorm.DB) {
				var m messageORM
				db.First(&m, "uuid = ?", "msg-1")
				if m.ContextSummary == nil || *m.ContextSummary != "updated summary" {
					t.Errorf("ContextSummary not updated")
				}
			},
		},
		{
			name: "異常系: DBエラーが発生した場合エラーになること",
			args: args{
				messageUUID: "msg-1",
				summary:     "updated summary",
			},
			setupData: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			wantErr: true,
			check:   nil,
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

			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewMessageRepository(db)
			if err := r.UpdateContextSummary(context.Background(), tt.args.messageUUID, tt.args.summary); (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.UpdateContextSummary() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.check != nil {
				tt.check(t, db)
			}
		})
	}
}

func TestMessageRepository_FindLatestMessageWithSummary(t *testing.T) {
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		want      *model.Message
		wantErr   bool
	}{
		{
			name: "正常系: 最新のサマリ付きメッセージが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				summary1 := "summary 1"
				summary2 := "summary 2"
				db.Create(&messageORM{
					UUID:           "msg-1",
					ChatUUID:       "chat-uuid",
					Role:           "assistant",
					Content:        "content 1",
					ContextSummary: &summary1,
					CreatedAt:      time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				})
				db.Create(&messageORM{
					UUID:           "msg-2",
					ChatUUID:       "chat-uuid",
					Role:           "assistant",
					Content:        "content 2",
					ContextSummary: &summary2,
					CreatedAt:      time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
				})
				db.Create(&messageORM{
					UUID:           "msg-3",
					ChatUUID:       "chat-uuid",
					Role:           "user",
					Content:        "content 3",
					ContextSummary: nil,
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				})
			},
			want: &model.Message{
				UUID:     "msg-2",
				ChatUUID: "chat-uuid",
				Role:     "assistant",
				Content:  "content 2",
			},
			wantErr: false,
		},
		{
			name: "正常系: サマリ付きメッセージが存在しない場合はnilが返ること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&messageORM{
					UUID:           "msg-1",
					ChatUUID:       "chat-uuid",
					Role:           "user",
					Content:        "content 1",
					ContextSummary: nil,
					CreatedAt:      time.Now(),
				})
			},
			want:    nil,
			wantErr: false,
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
			want:    nil,
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

			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewMessageRepository(db)
			got, err := r.FindLatestMessageWithSummary(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.FindLatestMessageWithSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("messageRepository.FindLatestMessageWithSummary() got = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("messageRepository.FindLatestMessageWithSummary() got nil, want %v", tt.want)
				} else {
					if got.UUID != tt.want.UUID {
						t.Errorf("messageRepository.FindLatestMessageWithSummary() UUID = %v, want %v", got.UUID, tt.want.UUID)
					}
				}
			}
		})
	}
}

func TestMessageRepository_FindLatestMessageByRole(t *testing.T) {
	type args struct {
		chatUUID string
		role     string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		want      *model.Message
		wantErr   bool
	}{
		{
			name: "正常系: 指定されたロールの最新メッセージが取得できること",
			args: args{
				chatUUID: "chat-uuid",
				role:     "assistant",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&messageORM{
					UUID:      "msg-1",
					ChatUUID:  "chat-uuid",
					Role:      "assistant",
					Content:   "content 1",
					CreatedAt: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				})
				db.Create(&messageORM{
					UUID:      "msg-2",
					ChatUUID:  "chat-uuid",
					Role:      "user",
					Content:   "content 2",
					CreatedAt: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
				})
				db.Create(&messageORM{
					UUID:      "msg-3",
					ChatUUID:  "chat-uuid",
					Role:      "assistant",
					Content:   "content 3",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				})
			},
			want: &model.Message{
				UUID:     "msg-3",
				ChatUUID: "chat-uuid",
				Role:     "assistant",
				Content:  "content 3",
			},
			wantErr: false,
		},
		{
			name: "正常系: 指定されたロールのメッセージが存在しない場合はnilが返ること",
			args: args{
				chatUUID: "chat-uuid",
				role:     "assistant",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&messageORM{
					UUID:      "msg-1",
					ChatUUID:  "chat-uuid",
					Role:      "user",
					Content:   "content 1",
					CreatedAt: time.Now(),
				})
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "異常系: DBエラーが発生した場合エラーになること",
			args: args{
				chatUUID: "chat-uuid",
				role:     "assistant",
			},
			setupData: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			want:    nil,
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

			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewMessageRepository(db)
			got, err := r.FindLatestMessageByRole(context.Background(), tt.args.chatUUID, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageRepository.FindLatestMessageByRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("messageRepository.FindLatestMessageByRole() got = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("messageRepository.FindLatestMessageByRole() got nil, want %v", tt.want)
				} else {
					if got.UUID != tt.want.UUID {
						t.Errorf("messageRepository.FindLatestMessageByRole() UUID = %v, want %v", got.UUID, tt.want.UUID)
					}
				}
			}
		})
	}
}
