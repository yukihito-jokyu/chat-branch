package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	domainUsecase "backend/internal/domain/usecase"
	"backend/internal/infrastructure/queue"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"google.golang.org/genai"
)

// GenAIClientWrapper は GenAIClient インターフェースの実装
type GenAIClientWrapper struct {
	client *genai.Client
}

func NewGenAIClientWrapper(client *genai.Client) domainUsecase.GenAIClient {
	return &GenAIClientWrapper{client: client}
}

func (w *GenAIClientWrapper) GenerateContentStream(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) func(func(*genai.GenerateContentResponse, error) bool) {
	return w.client.Models.GenerateContentStream(ctx, model, parts, config)
}

func (w *GenAIClientWrapper) GenerateContent(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	return w.client.Models.GenerateContent(ctx, model, parts, config)
}

type chatUsecase struct {
	chatRepo             repository.ChatRepository
	messageRepo          repository.MessageRepository
	messageSelectionRepo repository.MessageSelectionRepository
	transactionManager   repository.TransactionManager
	genaiClient          domainUsecase.GenAIClient
	publisher            message.Publisher
}

func NewChatUsecase(
	chatRepo repository.ChatRepository,
	messageRepo repository.MessageRepository,
	messageSelectionRepo repository.MessageSelectionRepository,
	transactionManager repository.TransactionManager,
	genaiClient domainUsecase.GenAIClient,
	publisher message.Publisher,
) domainUsecase.ChatUsecase {
	return &chatUsecase{
		chatRepo:             chatRepo,
		messageRepo:          messageRepo,
		messageSelectionRepo: messageSelectionRepo,
		transactionManager:   transactionManager,
		genaiClient:          genaiClient,
		publisher:            publisher,
	}
}

// ユーザーのメッセージを元に、GenAI にストリームを送信する
func (u *chatUsecase) StreamMessage(ctx context.Context, chatUUID string, outputChan chan<- string) error {
	slog.InfoContext(ctx, "メッセージストリーム処理開始", "chat_uuid", chatUUID)

	// 1. 最新のサマリを持つメッセージを取得
	latestSummaryMessage, err := u.messageRepo.FindLatestMessageWithSummary(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "最新サマリ取得失敗", "chat_uuid", chatUUID, "error", err)
		return err
	}

	// 2. メッセージ履歴の取得
	allMessages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "メッセージ履歴取得失敗", "chat_uuid", chatUUID, "error", err)
		return err
	}

	var contextMessages []*model.Message
	if latestSummaryMessage != nil {
		// サマリ以降のメッセージを抽出
		found := false
		for _, msg := range allMessages {
			if found {
				contextMessages = append(contextMessages, msg)
			}
			if msg.UUID == latestSummaryMessage.UUID {
				found = true
			}
		}
	} else {
		contextMessages = allMessages
	}

	// 3. プロンプト構築
	var parts []*genai.Content

	// サマリがあれば追加
	if latestSummaryMessage != nil && latestSummaryMessage.ContextSummary != nil {
		parts = append(parts, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "以下の会話の要約を踏まえて回答してください:\n" + *latestSummaryMessage.ContextSummary},
			},
		})
	}

	for _, msg := range contextMessages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		parts = append(parts, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: msg.Content},
			},
		})
	}

	// 4. GenAI 呼び出し
	client := u.genaiClient
	modelName := "gemini-2.5-flash"

	if len(contextMessages) == 0 && latestSummaryMessage == nil {
		return errors.New("no context to generate response")
	}

	if len(parts) == 0 {
		return errors.New("empty prompt")
	}

	iter := client.GenerateContentStream(ctx, modelName, parts, nil)

	var fullResponse string
	for resp, err := range iter {
		if err != nil {
			slog.ErrorContext(ctx, "GenAI APIからの受信エラー", "error", err)
			return err
		}
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					chunk := part.Text
					fullResponse += chunk
					outputChan <- chunk
				}
			}
		}
	}

	// 5. 生成された文章の保存
	assistantMessage := &model.Message{
		UUID:      uuid.New().String(),
		ChatUUID:  chatUUID,
		Role:      "assistant",
		Content:   fullResponse,
		CreatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, assistantMessage); err != nil {
		slog.ErrorContext(ctx, "アシスタントメッセージの保存に失敗しました", "error", err)
		return err
	}

	// 6. サマリ生成タスクのPublish
	topic := "chat_summary"
	payload, err := json.Marshal(chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "payloadのJSON変換に失敗しました", "error", err)
		return nil // 非同期タスクの失敗はメイン処理のエラーにはしない
	}

	if err := queue.PublishTask(u.publisher, topic, payload); err != nil {
		slog.ErrorContext(ctx, "サマリ生成タスクの登録に失敗しました", "error", err)
	} else {
		slog.InfoContext(ctx, "サマリ生成タスクを登録しました", "chat_uuid", chatUUID)
	}

	slog.InfoContext(ctx, "メッセージストリーム処理完了", "chat_uuid", chatUUID)
	return nil
}

// チャットの最初のメッセージを元に、GenAI にストリームを送信する
func (u *chatUsecase) FirstStreamChat(ctx context.Context, chatUUID string, outputChan chan<- string) error {
	slog.InfoContext(ctx, "チャットストリーム処理開始", "chat_uuid", chatUUID)

	// 1. チャットの存在確認
	_, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャットが見つかりません", "chat_uuid", chatUUID, "error", err)
		return err
	}

	// 2. メッセージ履歴の取得
	messages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "メッセージ履歴の取得に失敗しました", "chat_uuid", chatUUID, "error", err)
		return err
	}

	if len(messages) != 1 {
		return errors.New("invalid message state: expected exactly one initial message")
	}

	// 3. GenAI クライアント (注入されたものを使用)
	client := u.genaiClient

	modelName := "gemini-2.5-flash"

	// 4. プロンプトの構築とストリーム送信
	targetMessage := messages[0]

	// GenerateContentStream の呼び出し
	iter := client.GenerateContentStream(ctx, modelName, genai.Text(targetMessage.Content), nil)

	var fullResponse string

	// 5. ストリーム処理
	for resp, err := range iter {
		if err != nil {
			slog.ErrorContext(ctx, "GenAI APIからの受信エラー", "error", err)
			return err
		}

		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					chunk := part.Text
					fullResponse += chunk
					outputChan <- chunk
				}
			}
		}
	}

	// 6. 生成された文章の保存
	assistantMessage := &model.Message{
		UUID:      uuid.New().String(),
		ChatUUID:  chatUUID,
		Role:      "assistant",
		Content:   fullResponse,
		CreatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, assistantMessage); err != nil {
		slog.ErrorContext(ctx, "アシスタントメッセージの保存に失敗しました", "error", err)
		return err
	}

	slog.InfoContext(ctx, "チャットストリーム処理完了", "chat_uuid", chatUUID)
	return nil
}

// チャットを取得する
func (u *chatUsecase) GetChat(ctx context.Context, chatUUID string) (*model.Chat, error) {
	slog.InfoContext(ctx, "チャット取得処理開始", "chat_uuid", chatUUID)
	chat, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャット取得失敗", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}
	return chat, nil
}

// チャットのメッセージ一覧を取得する
func (u *chatUsecase) GetMessages(ctx context.Context, chatUUID string) ([]*model.Message, error) {
	slog.InfoContext(ctx, "メッセージ一覧取得処理開始", "chat_uuid", chatUUID)
	allMessages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "メッセージ一覧取得失敗", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}

	var mainMessages []*model.Message
	reportsByParent := make(map[string][]*model.Message)

	for _, msg := range allMessages {
		if msg.Role == "merge_report" {
			if msg.ParentMessageUUID != nil {
				reportsByParent[*msg.ParentMessageUUID] = append(reportsByParent[*msg.ParentMessageUUID], msg)
			}
		} else {
			mainMessages = append(mainMessages, msg)
		}
	}

	for _, msg := range mainMessages {
		if reports, ok := reportsByParent[msg.UUID]; ok {
			msg.MergeReports = reports
		}
	}

	return mainMessages, nil
}

// メッセージを送信する
func (u *chatUsecase) SendMessage(ctx context.Context, chatUUID string, content string) (*model.Message, error) {
	slog.InfoContext(ctx, "メッセージ送信処理開始", "chat_uuid", chatUUID)

	// チャットの存在確認
	_, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャットが見つかりません", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}

	message := &model.Message{
		UUID:      uuid.New().String(),
		ChatUUID:  chatUUID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, message); err != nil {
		slog.ErrorContext(ctx, "メッセージ保存失敗", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "メッセージ送信成功", "message_uuid", message.UUID)
	return message, nil
}

// フォークプレビューを生成する
func (u *chatUsecase) GenerateForkPreview(ctx context.Context, chatUUID string, req model.ForkPreviewRequest) (*model.ForkPreviewResponse, error) {
	slog.InfoContext(ctx, "フォークプレビュー生成開始", "chat_uuid", chatUUID, "target_message_uuid", req.TargetMessageUUID)

	// 1. メッセージ履歴の取得
	allMessages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// 2. 対象メッセージの特定とサマリの探索
	var targetMessage *model.Message
	var latestSummaryMessage *model.Message
	var targetIndex int = -1

	// 対象メッセージのインデックスを探す
	for i, msg := range allMessages {
		if msg.UUID == req.TargetMessageUUID {
			targetMessage = msg
			targetIndex = i
			break
		}
	}

	if targetMessage == nil {
		return nil, fmt.Errorf("target message not found")
	}

	// 対象メッセージから遡ってサマリを探す（対象メッセージ含む）
	for i := targetIndex; i >= 0; i-- {
		if allMessages[i].ContextSummary != nil && *allMessages[i].ContextSummary != "" {
			latestSummaryMessage = allMessages[i]
			break
		}
	}

	// 3. プロンプト構築用のメッセージ抽出
	var targetMessages []*model.Message

	if latestSummaryMessage != nil {

		startIndex := -1
		for i, msg := range allMessages {
			if msg.UUID == latestSummaryMessage.UUID {
				startIndex = i
				break
			}
		}

		if startIndex != -1 && startIndex < targetIndex {
			// サマリメッセージの次のメッセージから、対象メッセージの一つ前まで
			targetMessages = allMessages[startIndex+1 : targetIndex]
		} else if startIndex == targetIndex {
			// 対象メッセージ自体がサマリを持っている場合、履歴は空
			targetMessages = []*model.Message{}
		}
	} else {
		// サマリがない場合
		// 最初から対象メッセージの一つ前まで
		targetMessages = allMessages[0:targetIndex]
	}

	// 4. プロンプト構築
	var parts []*genai.Content

	if latestSummaryMessage != nil && latestSummaryMessage.ContextSummary != nil {
		parts = append(parts, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "以下の会話の要約を踏まえてください:\n" + *latestSummaryMessage.ContextSummary},
			},
		})
	}

	for _, msg := range targetMessages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		parts = append(parts, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: msg.Content},
			},
		})
	}

	// 対象メッセージの内容
	// 対象メッセージは必ず含める（roleに応じて）
	targetRole := "user"
	if targetMessage.Role == "assistant" {
		targetRole = "model"
	}
	parts = append(parts, &genai.Content{
		Role: targetRole,
		Parts: []*genai.Part{
			{Text: targetMessage.Content},
		},
	})

	// 指示プロンプト
	prompt := fmt.Sprintf(`
ユーザーは上記の会話の最後のメッセージの以下の部分を選択して、新しい話題（チャット）を開始しようとしています。
選択範囲: "%s"
(範囲: %d文字目から%d文字目)

以下のJSON形式で、新しいチャットのタイトル案と、これまでの文脈を考慮した「新しいチャットの冒頭に設定するコンテキスト（要約）」を生成してください。
コンテキストは、選択された話題について深掘りするための導入として機能するようにしてください。
生成はJSONのみで良いです。説明は不要です。

JSON形式:
{
  "suggested_title": "タイトル案",
  "generated_context": "生成されたコンテキスト"
}
`, req.SelectedText, req.RangeStart, req.RangeEnd)

	parts = append(parts, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: prompt},
		},
	})

	// 5. GenAI 呼び出し
	client := u.genaiClient
	modelName := "gemini-2.5-flash"

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	resp, err := client.GenerateContent(ctx, modelName, parts, config)
	if err != nil {
		return nil, fmt.Errorf("GenAI呼び出しに失敗: %w", err)
	}

	// 6. レスポンス解析
	var generatedText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				generatedText += part.Text
			}
		}
	}

	var result model.ForkPreviewResponse
	if err := json.Unmarshal([]byte(generatedText), &result); err != nil {
		return nil, fmt.Errorf("JSON出力に失敗: %w", err)
	}

	return &result, nil
}

// チャットをフォークする
func (u *chatUsecase) ForkChat(ctx context.Context, params model.ForkChatParams) (string, error) {
	slog.InfoContext(ctx, "チャットフォーク処理開始", "parent_chat_uuid", params.ParentChatUUID, "target_message_uuid", params.TargetMessageUUID)

	// 1. 親チャットの存在確認
	parentChat, err := u.chatRepo.FindByID(ctx, params.ParentChatUUID)
	if err != nil {
		return "", fmt.Errorf("親チャットの存在確認に失敗: %w", err)
	}

	// 2. トランザクション処理
	// MessageSelection作成 -> Chat作成 -> Message作成
	newChatUUID := uuid.New().String()
	selectionUUID := uuid.New().String()
	messageUUID := uuid.New().String()

	err = u.transactionManager.Do(ctx, func(ctx context.Context) error {
		// 2-1. MessageSelection作成
		selection := &model.MessageSelection{
			UUID:         selectionUUID,
			SelectedText: params.SelectedText,
			RangeStart:   params.RangeStart,
			RangeEnd:     params.RangeEnd,
			CreatedAt:    time.Now(),
		}
		if err := u.messageSelectionRepo.Create(ctx, selection); err != nil {
			return fmt.Errorf("メッセージ選択の作成に失敗: %w", err)
		}

		// 2-2. Chat作成
		newChat := &model.Chat{
			UUID:                 newChatUUID,
			ProjectUUID:          parentChat.ProjectUUID, // 親チャットと同じプロジェクト
			ParentUUID:           &params.ParentChatUUID,
			SourceMessageUUID:    &params.TargetMessageUUID,
			MessageSelectionUUID: &selectionUUID,
			Title:                params.Title,
			Status:               "open",
			ContextSummary:       params.ContextSummary,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		if err := u.chatRepo.Create(ctx, newChat); err != nil {
			return fmt.Errorf("新しいチャットの作成に失敗: %w", err)
		}

		// 2-3. Message作成 (最初のメッセージ)
		// タイトルとコンテキストサマリを結合した文章をユーザーメッセージとして保存
		initialContent := fmt.Sprintf("%s\n\n%s", params.Title, params.ContextSummary)
		message := &model.Message{
			UUID:      messageUUID,
			ChatUUID:  newChatUUID,
			Role:      "assistant",
			Content:   initialContent,
			CreatedAt: time.Now(),
		}
		if err := u.messageRepo.Create(ctx, message); err != nil {
			return fmt.Errorf("初期メッセージの作成に失敗: %w", err)
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "チャットフォーク処理失敗", "error", err)
		return "", err
	}

	slog.InfoContext(ctx, "チャットフォーク処理完了", "new_chat_uuid", newChatUUID)

	return newChatUUID, nil
}

// マージプレビューを生成する
func (u *chatUsecase) GetMergePreview(ctx context.Context, chatUUID string) (*model.MergePreview, error) {
	slog.InfoContext(ctx, "マージプレビュー生成開始", "chat_uuid", chatUUID)

	// 1. チャット情報の取得 (親からforkした理由 = ContextSummary を取得)
	chat, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		return nil, fmt.Errorf("チャット取得失敗: %w", err)
	}

	// 2. 子チャットの最新のサマリを取得
	latestSummaryMessage, err := u.messageRepo.FindLatestMessageWithSummary(ctx, chatUUID)
	if err != nil {
		return nil, fmt.Errorf("最新サマリ取得失敗: %w", err)
	}

	// 3. role='assistant' の最新メッセージを取得
	latestAssistantMessage, err := u.messageRepo.FindLatestMessageByRole(ctx, chatUUID, "assistant")
	if err != nil {
		return nil, fmt.Errorf("最新アシスタントメッセージ取得失敗: %w", err)
	}

	// 4. プロンプト構築
	var parts []*genai.Content

	prompt := "以下の情報を元に、子チャットでの議論の流れと結論を要約してください。\n\n"

	prompt += "## 親チャットからForkした理由 (文脈)\n"
	if chat.ContextSummary != "" {
		prompt += chat.ContextSummary + "\n\n"
	} else {
		prompt += "なし\n\n"
	}

	prompt += "## 子チャットの最新のサマリ (途中経過)\n"
	if latestSummaryMessage != nil && latestSummaryMessage.ContextSummary != nil {
		prompt += *latestSummaryMessage.ContextSummary + "\n\n"
	} else {
		prompt += "なし\n\n"
	}

	prompt += "## 最新のAI回答 (直近の結論)\n"
	if latestAssistantMessage != nil {
		prompt += latestAssistantMessage.Content + "\n\n"
	} else {
		prompt += "なし\n\n"
	}

	prompt += `
出力フォーマット:
## 議論の流れ
(ここに議論の流れを記述)

## 結論
(ここに結論を記述)
`

	parts = append(parts, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: prompt},
		},
	})

	// 5. GenAI 呼び出し
	client := u.genaiClient
	modelName := "gemini-2.5-flash"

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "text/plain",
	}

	resp, err := client.GenerateContent(ctx, modelName, parts, config)
	if err != nil {
		return nil, fmt.Errorf("GenAI呼び出しに失敗: %w", err)
	}

	// 6. レスポンス解析
	var generatedText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				generatedText += part.Text
			}
		}
	}

	return &model.MergePreview{
		SuggestedSummary: generatedText,
	}, nil
}

// チャットをマージする
func (u *chatUsecase) MergeChat(ctx context.Context, chatUUID string, params model.MergeChatParams) (*model.MergeChatResult, error) {
	slog.InfoContext(ctx, "チャットマージ処理開始", "chat_uuid", chatUUID, "parent_chat_uuid", params.ParentChatUUID)

	// 1. 子チャットの取得
	childChat, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		return nil, fmt.Errorf("子チャットの取得に失敗: %w", err)
	}

	if childChat.SourceMessageUUID == nil {
		return nil, fmt.Errorf("子チャットにソースメッセージが設定されていません")
	}

	reportMessageID := uuid.New().String()

	// 2. トランザクション処理
	err = u.transactionManager.Do(ctx, func(ctx context.Context) error {
		// 2-1. マージレポートメッセージの作成
		// 親チャットに追加するので、ChatUUIDはParentChatUUIDになる
		reportMessage := &model.Message{
			UUID:              reportMessageID,
			ChatUUID:          params.ParentChatUUID,
			Role:              "merge_report",
			Content:           params.SummaryContent,
			ParentMessageUUID: childChat.SourceMessageUUID, // どのメッセージから派生したチャットがマージされたかを示す
			SourceChatUUID:    &chatUUID,                   // どのチャットがマージされたか
			CreatedAt:         time.Now(),
		}

		if err := u.messageRepo.Create(ctx, reportMessage); err != nil {
			return fmt.Errorf("マージレポートメッセージの作成に失敗: %w", err)
		}

		// 2-2. 子チャットのステータス更新
		if err := u.chatRepo.UpdateStatus(ctx, chatUUID, "merged"); err != nil {
			return fmt.Errorf("子チャットのステータス更新に失敗: %w", err)
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "チャットマージ処理失敗", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "チャットマージ処理完了", "chat_uuid", chatUUID)

	return &model.MergeChatResult{
		ReportMessageID: reportMessageID,
		SummaryContent:  params.SummaryContent,
	}, nil
}

// チャットをクローズする
func (u *chatUsecase) CloseChat(ctx context.Context, chatUUID string) (string, error) {
	slog.InfoContext(ctx, "チャットクローズ処理開始", "chat_uuid", chatUUID)

	// 1. チャットの存在確認
	_, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャットが見つかりません", "chat_uuid", chatUUID, "error", err)
		return "", err
	}

	// 2. ステータスを closed に更新
	if err := u.chatRepo.UpdateStatus(ctx, chatUUID, "closed"); err != nil {
		slog.ErrorContext(ctx, "チャットステータス更新失敗", "chat_uuid", chatUUID, "error", err)
		return "", err
	}

	slog.InfoContext(ctx, "チャットクローズ処理完了", "chat_uuid", chatUUID)
	return chatUUID, nil
}
