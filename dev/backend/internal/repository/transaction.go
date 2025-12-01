package repository

import (
	"backend/internal/domain/repository"
	"context"

	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) repository.TransactionManager {
	return &transactionManager{db: db}
}

type txKey struct{}

// トランザクションを開始する処理
func (m *transactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, txKey{}, tx)
		return fn(ctxWithTx)
	})
}

// トランザクションを取得する処理
func getDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
