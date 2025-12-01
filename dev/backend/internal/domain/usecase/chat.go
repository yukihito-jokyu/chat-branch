package usecase

import (
	"backend/internal/domain/model"
	"context"
)

type ChatUsecase interface {
	// チャットの最初のメッセージを元に、GenAI にストリームを送信する
	FirstStreamChat(ctx context.Context, chatUUID string, outputChan chan<- string) error
	// チャットを取得する
	GetChat(ctx context.Context, chatUUID string) (*model.Chat, error)
	// チャットのメッセージ一覧を取得する
	GetMessages(ctx context.Context, chatUUID string) ([]*model.Message, error)
}
