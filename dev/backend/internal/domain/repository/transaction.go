package repository

import "context"

type TransactionManager interface {
	// トランザクションを実行する処理
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
