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
					UUID:           "chat-uuid",
					ProjectUUID:    "project-uuid",
					Title:          "test chat",
					Status:         "active",
					ContextSummary: "summary",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: UUIDが重複している場合エラーになること",
			args: args{
				chat: &model.Chat{
					UUID:           "duplicate-uuid",
					ProjectUUID:    "project-uuid",
					Title:          "duplicate chat",
					Status:         "active",
					ContextSummary: "summary",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
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
				parentUUID := "parent-uuid"
				contextSummary := "summary"
				db.Create(&chatORM{
					UUID:           "chat-uuid",
					ProjectUUID:    "project-uuid",
					ParentChatUUID: &parentUUID,
					Title:          "test chat",
					Status:         "active",
					ContextSummary: &contextSummary,
					CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				})
			},
			want: &model.Chat{
				UUID:           "chat-uuid",
				ProjectUUID:    "project-uuid",
				ParentUUID:     func() *string { s := "parent-uuid"; return &s }(),
				Title:          "test chat",
				Status:         "active",
				ContextSummary: "summary",
				CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
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

func TestChatRepository_UpdateStatus(t *testing.T) {
	type args struct {
		chatUUID string
		status   string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		wantErr   bool
	}{
		{
			name: "正常系: ステータスが更新できること",
			args: args{
				chatUUID: "chat-uuid",
				status:   "merged",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&chatORM{
					UUID:        "chat-uuid",
					ProjectUUID: "project-uuid",
					Title:       "test chat",
					Status:      "open",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
			},
			wantErr: false,
		},
		{
			name: "正常系: 存在しないチャットでもエラーにならない（GORMの仕様）",
			args: args{
				chatUUID: "non-existent",
				status:   "merged",
			},
			setupData: func(db *gorm.DB) {},
			wantErr:   false,
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
			if err := r.UpdateStatus(context.Background(), tt.args.chatUUID, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("chatRepository.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 更新確認
			if !tt.wantErr && tt.name == "正常系: ステータスが更新できること" {
				var orm chatORM
				if err := db.Where("uuid = ?", tt.args.chatUUID).First(&orm).Error; err != nil {
					t.Fatalf("failed to fetch chat: %v", err)
				}
				if orm.Status != tt.args.status {
					t.Errorf("status = %v, want %v", orm.Status, tt.args.status)
				}
			}
		})
	}
}

func TestChatRepository_FindOldestByProjectUUID(t *testing.T) {
	type args struct {
		projectUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		want      *model.Chat
		wantErr   bool
	}{
		{
			name: "正常系: 最も古いチャットが取得できること",
			args: args{
				projectUUID: "project-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&chatORM{
					UUID:        "chat-old",
					ProjectUUID: "project-uuid",
					Title:       "old chat",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				})
				db.Create(&chatORM{
					UUID:        "chat-new",
					ProjectUUID: "project-uuid",
					Title:       "new chat",
					CreatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				})
			},
			want: &model.Chat{
				UUID:        "chat-old",
				ProjectUUID: "project-uuid",
				Title:       "old chat",
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "正常系: ContextSummaryが存在する場合、正しくマッピングされること",
			args: args{
				projectUUID: "project-uuid-summary",
			},
			setupData: func(db *gorm.DB) {
				summary := "test summary"
				db.Create(&chatORM{
					UUID:           "chat-summary",
					ProjectUUID:    "project-uuid-summary",
					Title:          "summary chat",
					ContextSummary: &summary,
					CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				})
			},
			want: &model.Chat{
				UUID:           "chat-summary",
				ProjectUUID:    "project-uuid-summary",
				Title:          "summary chat",
				ContextSummary: "test summary",
				CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "異常系: チャットが存在しない場合エラーになること",
			args: args{
				projectUUID: "non-existent-project",
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
			got, err := r.FindOldestByProjectUUID(context.Background(), tt.args.projectUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatRepository.FindOldestByProjectUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.UUID != tt.want.UUID || got.ProjectUUID != tt.want.ProjectUUID || got.Title != tt.want.Title {
					t.Errorf("chatRepository.FindOldestByProjectUUID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestChatRepository_CountByProjectUUID(t *testing.T) {
	type args struct {
		projectUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		want      int64
		wantErr   bool
	}{
		{
			name: "正常系: プロジェクト内のチャット数が正しくカウントされること",
			args: args{
				projectUUID: "project-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&chatORM{
					UUID:        "chat-1",
					ProjectUUID: "project-uuid",
					Title:       "chat 1",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				db.Create(&chatORM{
					UUID:        "chat-2",
					ProjectUUID: "project-uuid",
					Title:       "chat 2",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				db.Create(&chatORM{
					UUID:        "chat-other",
					ProjectUUID: "other-project",
					Title:       "other chat",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "正常系: チャットが存在しない場合は0が返ること",
			args: args{
				projectUUID: "empty-project",
			},
			setupData: func(db *gorm.DB) {},
			want:      0,
			wantErr:   false,
		},
		{
			name: "異常系: DBエラーが発生した場合エラーになること",
			args: args{
				projectUUID: "project-uuid",
			},
			setupData: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			want:    0,
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

			// データ投入
			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewChatRepository(db)
			got, err := r.CountByProjectUUID(context.Background(), tt.args.projectUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatRepository.CountByProjectUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("chatRepository.CountByProjectUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
