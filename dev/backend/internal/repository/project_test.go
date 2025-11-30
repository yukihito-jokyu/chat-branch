package repository

import (
	"backend/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProjectRepository_FindAllByUserUUID(t *testing.T) {
	// テスト用のDBセットアップ
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// マイグレーション
	if err := db.AutoMigrate(&projectORM{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	type args struct {
		userUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupDB   func(db *gorm.DB)
		wantCount int
		wantErr   bool
	}{
		{
			name: "正常系: 指定したユーザーのプロジェクトが取得できること",
			args: args{
				userUUID: "user-1",
			},
			setupDB: func(db *gorm.DB) {
				// テストデータの投入
				projects := []projectORM{
					{UUID: "p1", UserID: "user-1", Title: "Project 1", UpdatedAt: time.Now()},
					{UUID: "p2", UserID: "user-1", Title: "Project 2", UpdatedAt: time.Now()},
					{UUID: "p3", UserID: "user-2", Title: "Project 3", UpdatedAt: time.Now()},
				}
				db.Create(&projects)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "正常系: プロジェクトが存在しない場合は空のリストが返ること",
			args: args{
				userUUID: "user-3",
			},
			setupDB: func(db *gorm.DB) {
				// データなし
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "異常系: データベースエラーの場合エラーになること",
			args: args{
				userUUID: "user-1",
			},
			setupDB: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストごとにデータをクリア
			db.Exec("DELETE FROM projects")
			tt.setupDB(db)

			r := NewProjectRepository(db)
			got, err := r.FindAllByUserUUID(context.Background(), tt.args.userUUID)

			if (err != nil) != tt.wantErr {
				t.Errorf("projectRepository.FindAllByUserUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Len(t, got, tt.wantCount)
				for _, p := range got {
					assert.Equal(t, tt.args.userUUID, p.UserID)
					assert.IsType(t, &model.Project{}, p)
				}
			}
		})
	}
}
