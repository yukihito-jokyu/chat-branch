package handler

import (
	"backend/internal/usecase"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatHandler interface {
	FirstStreamChat(c echo.Context) error
}

type chatHandler struct {
	chatUsecase usecase.ChatUsecase
}

func NewChatHandler(chatUsecase usecase.ChatUsecase) ChatHandler {
	return &chatHandler{
		chatUsecase: chatUsecase,
	}
}

// チャットの最初のメッセージを元に、GenAI にストリームを送信する
func (h *chatHandler) FirstStreamChat(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	slog.InfoContext(ctx, "FirstStreamChat リクエスト受信", "chat_uuid", chatUUID)

	// SSE ヘッダーの設定
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	outputChan := make(chan string)
	errChan := make(chan error)

	go func() {
		// defer close(outputChan) // 成功時のみ閉じる
		// defer close(errChan) // 閉じない
		if err := h.chatUsecase.FirstStreamChat(ctx, chatUUID, outputChan); err != nil {
			errChan <- err
		} else {
			close(outputChan)
		}
	}()

	enc := json.NewEncoder(c.Response())

	for {
		select {
		case chunk, ok := <-outputChan:
			if !ok {
				// チャネルが閉じられたら終了
				data := map[string]string{
					"status": "done",
				}
				fmt.Fprintf(c.Response(), "data: ")
				enc.Encode(data)
				fmt.Fprintf(c.Response(), "\n\n")
				c.Response().Flush()
				return nil
			}
			// チャンクを送信
			data := map[string]string{
				"chunk":  chunk,
				"status": "processing",
			}
			fmt.Fprintf(c.Response(), "data: ")
			enc.Encode(data)
			fmt.Fprintf(c.Response(), "\n\n")
			c.Response().Flush()

		case err := <-errChan:
			if err != nil {
				slog.ErrorContext(ctx, "StreamChat エラー発生", "error", err)
				// エラーをクライアントに通知（必要であれば）
				// SSEの仕様上、接続を切るか、エラーイベントを送る
				return err
			}
			return nil

		case <-ctx.Done():
			// クライアント切断
			slog.InfoContext(ctx, "クライアント切断")
			return nil
		}
	}
}
