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

func Test_messageSelectionRepository_Create(t *testing.T) {
	type args struct {
		ctx       context.Context
		selection *model.MessageSelection
	}
	tests := []struct {
		name      string
		args      args
		setupDB   func(db *gorm.DB)
		wantErr   bool
		assertion func(t *testing.T, db *gorm.DB, args args)
	}{
		{
			name: "正常系: メッセージ選択を保存できる",
			args: args{
				ctx: context.Background(),
				selection: &model.MessageSelection{
					UUID:         "test-uuid",
					SelectedText: "selected text",
					RangeStart:   10,
					RangeEnd:     20,
					CreatedAt:    time.Now(),
				},
			},
			setupDB: func(db *gorm.DB) {
				// 初期データなし
			},
			wantErr: false,
			assertion: func(t *testing.T, db *gorm.DB, args args) {
				var stored messageSelectionORM
				err := db.Where("uuid = ?", args.selection.UUID).First(&stored).Error
				assert.NoError(t, err)
				assert.Equal(t, args.selection.UUID, stored.UUID)
				assert.Equal(t, args.selection.SelectedText, stored.SelectedText)
				assert.Equal(t, args.selection.RangeStart, stored.RangeStart)
				assert.Equal(t, args.selection.RangeEnd, stored.RangeEnd)
				assert.NotEmpty(t, stored.CreatedID)
				assert.NotZero(t, stored.CreatedAt)
				assert.NotZero(t, stored.UpdatedAt)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In-memory SQLite DB setup
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}

			// Migration
			if err := db.AutoMigrate(&messageSelectionORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			r := NewMessageSelectionRepository(db)
			if err := r.Create(tt.args.ctx, tt.args.selection); (err != nil) != tt.wantErr {
				t.Errorf("messageSelectionRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.assertion != nil {
				tt.assertion(t, db, tt.args)
			}
		})
	}
}
