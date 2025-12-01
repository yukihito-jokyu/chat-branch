package handler

import (
	domainModel "backend/internal/domain/model"
	"backend/internal/domain/usecase"
	"backend/internal/handler/model"
	handlerModel "backend/internal/handler/model"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type chatHandler struct {
	chatUsecase usecase.ChatUsecase
}

func NewChatHandler(chatUsecase usecase.ChatUsecase) *chatHandler {
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

// チャットのメッセージをストリーミング送信する
func (h *chatHandler) StreamMessage(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	slog.InfoContext(ctx, "StreamMessage リクエスト受信", "chat_uuid", chatUUID)

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
		if err := h.chatUsecase.StreamMessage(ctx, chatUUID, outputChan); err != nil {
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
				slog.ErrorContext(ctx, "StreamMessage エラー発生", "error", err)
				// エラーをクライアントに通知（必要であれば）
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

// チャットを取得する
func (h *chatHandler) GetChat(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	slog.InfoContext(ctx, "GetChat リクエスト受信", "chat_uuid", chatUUID)

	chat, err := h.chatUsecase.GetChat(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "GetChat エラー", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := model.GetChatResponse{
		UUID:           chat.UUID,
		ProjectUUID:    chat.ProjectUUID,
		ParentUUID:     chat.ParentUUID,
		Title:          chat.Title,
		Status:         chat.Status,
		ContextSummary: chat.ContextSummary,
	}

	slog.InfoContext(ctx, "チャットの取得に成功", "chat_uuid", chat.UUID)
	return c.JSON(http.StatusOK, res)
}

// チャットのメッセージ一覧を取得する
func (h *chatHandler) GetMessages(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	messages, err := h.chatUsecase.GetMessages(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "GetMessages エラー", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := make([]model.MessageResponse, len(messages))
	for i, m := range messages {
		res[i] = mapMessageToResponse(m)
	}

	slog.InfoContext(ctx, "メッセージ一覧の取得に成功", "chat_uuid", chatUUID, "count", len(res))
	return c.JSON(http.StatusOK, res)
}

// メッセージを送信する
func (h *chatHandler) SendMessage(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	var req model.SendMessageRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "リクエストボディのバインドエラー", "error", err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "リクエストボディのバインドに失敗しました",
		})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "content が空です",
		})
	}

	message, err := h.chatUsecase.SendMessage(ctx, chatUUID, req.Content)
	if err != nil {
		slog.ErrorContext(ctx, "SendMessage エラー", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := model.MessageResponse{
		UUID:           message.UUID,
		Role:           message.Role,
		Content:        message.Content,
		SourceChatUUID: message.SourceChatUUID,
		// Forks は新規作成時は空
		Forks:        []model.ForkResponse{},
		MergeReports: []model.MessageResponse{},
	}

	slog.InfoContext(ctx, "メッセージ送信成功", "chat_uuid", chatUUID, "message_uuid", message.UUID)
	return c.JSON(http.StatusOK, res)
}

// フォークプレビューを生成する
func (h *chatHandler) GenerateForkPreview(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	var req domainModel.ForkPreviewRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "リクエストボディのバインドエラー", "error", err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "リクエストボディのバインドに失敗しました",
		})
	}

	res, err := h.chatUsecase.GenerateForkPreview(ctx, chatUUID, req)
	if err != nil {
		slog.ErrorContext(ctx, "GenerateForkPreview エラー", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// フォークチャットを生成する
func (h *chatHandler) ForkChat(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	var req handlerModel.ForkChatRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "リクエストボディのバインドエラー", "error", err)
		return c.JSON(http.StatusBadRequest, model.Response{Status: "error", Message: err.Error()})
	}

	// パスパラメータとリクエストボディの整合性チェック
	if req.ParentChatUUID != "" && req.ParentChatUUID != chatUUID {
		slog.WarnContext(ctx, "パスパラメータのチャットIDとリクエストボディの親チャットIDが一致しない", "path", chatUUID, "body", req.ParentChatUUID)
		return c.JSON(http.StatusBadRequest, model.Response{Status: "error", Message: "親チャットIDと一致しない"})
	}
	// リクエストボディに親チャットIDがない場合はパスパラメータを使用
	if req.ParentChatUUID == "" {
		req.ParentChatUUID = chatUUID
	}

	params := domainModel.ForkChatParams{
		TargetMessageUUID: req.TargetMessageUUID,
		ParentChatUUID:    req.ParentChatUUID,
		SelectedText:      req.SelectedText,
		RangeStart:        req.RangeStart,
		RangeEnd:          req.RangeEnd,
		Title:             req.Title,
		ContextSummary:    req.ContextSummary,
	}

	newChatID, err := h.chatUsecase.ForkChat(ctx, params)
	if err != nil {
		slog.ErrorContext(ctx, "フォークチャットの生成に失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{Status: "error", Message: err.Error()})
	}

	res := handlerModel.ForkChatResponse{
		NewChatID: newChatID,
		Message:   "子チャットを作成しました",
	}

	return c.JSON(http.StatusOK, res)
}

// マージプレビューを生成する
func (h *chatHandler) GetMergePreview(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	slog.InfoContext(ctx, "GetMergePreview リクエスト受信", "chat_uuid", chatUUID)

	preview, err := h.chatUsecase.GetMergePreview(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "GetMergePreview エラー", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := handlerModel.MergePreviewResponse{
		SuggestedSummary: preview.SuggestedSummary,
	}

	return c.JSON(http.StatusOK, res)
}

// チャットをマージする
func (h *chatHandler) MergeChat(c echo.Context) error {
	chatUUID := c.Param("chat_uuid")
	ctx := c.Request().Context()

	slog.InfoContext(ctx, "MergeChat リクエスト受信", "chat_uuid", chatUUID)

	var req handlerModel.MergeChatRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "リクエストボディのバインドエラー", "error", err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "リクエストボディのバインドに失敗しました",
		})
	}

	params := domainModel.MergeChatParams{
		ParentChatUUID: req.ParentChatUUID,
		SummaryContent: req.SummaryContent,
	}

	result, err := h.chatUsecase.MergeChat(ctx, chatUUID, params)
	if err != nil {
		slog.ErrorContext(ctx, "MergeChat エラー", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := handlerModel.MergeChatResponse{
		ReportMessageID: result.ReportMessageID,
		SummaryContent:  result.SummaryContent,
	}

	return c.JSON(http.StatusOK, res)
}

func mapMessageToResponse(m *domainModel.Message) model.MessageResponse {
	forks := make([]model.ForkResponse, len(m.Forks))
	for j, f := range m.Forks {
		forks[j] = model.ForkResponse{
			ChatUUID:     f.ChatUUID,
			SelectedText: f.SelectedText,
			RangeStart:   f.RangeStart,
			RangeEnd:     f.RangeEnd,
		}
	}

	mergeReports := make([]model.MessageResponse, len(m.MergeReports))
	for j, r := range m.MergeReports {
		mergeReports[j] = mapMessageToResponse(r)
	}

	return model.MessageResponse{
		UUID:           m.UUID,
		Role:           m.Role,
		Content:        m.Content,
		Forks:          forks,
		SourceChatUUID: m.SourceChatUUID,
		MergeReports:   mergeReports,
	}
}
