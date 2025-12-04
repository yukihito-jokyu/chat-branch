package repository

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestEdgeRepository_FindEdgesByChatID(t *testing.T) {
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
			name: "正常系: チャットIDに紐づくエッジが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&edgeORM{
					UUID:              "edge-1",
					ChatUUID:          "chat-uuid",
					SourceMessageUUID: "msg-1",
					TargetMessageUUID: "msg-2",
				})
				db.Create(&edgeORM{
					UUID:              "edge-2",
					ChatUUID:          "chat-uuid",
					SourceMessageUUID: "msg-2",
					TargetMessageUUID: "msg-3",
				})
				db.Create(&edgeORM{
					UUID:              "edge-other",
					ChatUUID:          "other-chat-uuid",
					SourceMessageUUID: "msg-other-1",
					TargetMessageUUID: "msg-other-2",
				})
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "正常系: エッジが存在しない場合は空のリストが返ること",
			args: args{
				chatUUID: "empty-chat-uuid",
			},
			setupData: func(db *gorm.DB) {},
			wantLen:   0,
			wantErr:   false,
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
			if err := db.AutoMigrate(&edgeORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			// データ投入
			if tt.setupData != nil {
				tt.setupData(db)
			}

			r := NewEdgeRepository(db)
			got, err := r.FindEdgesByChatID(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("edgeRepository.FindEdgesByChatID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("edgeRepository.FindEdgesByChatID() length = %v, want %v", len(got), tt.wantLen)
			}

			// データの中身の確認
			if tt.name == "正常系: チャットIDに紐づくエッジが取得できること" {
				found := false
				for _, edge := range got {
					if edge.UUID == "edge-1" {
						found = true
						if edge.SourceMessageUUID != "msg-1" || edge.TargetMessageUUID != "msg-2" {
							t.Errorf("edge-1 content mismatch")
						}
					}
				}
				if !found {
					t.Errorf("edge-1 not found")
				}
			}
		})
	}
}
