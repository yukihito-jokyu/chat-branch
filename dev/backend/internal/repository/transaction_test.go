package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTransactionManager_Do(t *testing.T) {
	// テスト用のテーブル定義
	type TestModel struct {
		ID   uint   `gorm:"primaryKey"`
		Name string `gorm:"unique"`
	}

	tests := []struct {
		name      string
		fn        func(ctx context.Context, db *gorm.DB) error
		wantErr   bool
		checkData func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "正常系: トランザクションがコミットされること",
			fn: func(ctx context.Context, db *gorm.DB) error {
				// トランザクション内でデータを挿入
				// 注意: ここではリポジトリの内部実装(getDB)を模倣して、コンテキストからDBを取得する代わりに
				// テストヘルパー内で渡されたdbを使うか、あるいはTransactionManagerが正しくコンテキストにDBをセットしているかを確認する必要がある。
				// しかし、TransactionManager.Doはコールバックにコンテキストを渡すので、
				// そのコンテキストを使ってDB操作をする必要がある。
				// ここでは簡易的に、TransactionManagerの実装が正しいと仮定し、
				// コールバック内のコンテキストからDBを取り出すロジック(getDB)はprivateなのでテストできない。
				// 代わりに、TransactionManagerと同じパッケージにいるので getDB を呼べるはず。

				tx := getDB(ctx, db)
				return tx.Create(&TestModel{Name: "commit-test"}).Error
			},
			wantErr: false,
			checkData: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&TestModel{}).Where("name = ?", "commit-test").Count(&count)
				assert.Equal(t, int64(1), count)
			},
		},
		{
			name: "異常系: エラー時にロールバックされること",
			fn: func(ctx context.Context, db *gorm.DB) error {
				tx := getDB(ctx, db)
				if err := tx.Create(&TestModel{Name: "rollback-test"}).Error; err != nil {
					return err
				}
				return errors.New("error for rollback")
			},
			wantErr: true,
			checkData: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&TestModel{}).Where("name = ?", "rollback-test").Count(&count)
				assert.Equal(t, int64(0), count)
			},
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
			if err := db.AutoMigrate(&TestModel{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			tm := NewTransactionManager(db)

			err = tm.Do(context.Background(), func(ctx context.Context) error {
				return tt.fn(ctx, db)
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionManager.Do() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.checkData(t, db)
		})
	}
}
