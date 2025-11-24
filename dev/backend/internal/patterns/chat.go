package patterns

import (
	"context"
	"fmt"

	"backend/internal/client"

	"github.com/google/generative-ai-go/genai"
)

// RunChatSession はチャットセッション（履歴管理）を実行します
func RunChatSession(ctx context.Context, c *genai.Client) error {
	fmt.Println("=== Pattern 2: Chat Session (History Management) ===")

	model := c.GenerativeModel(client.ModelGeminiPro)
	cs := model.StartChat()

	// 1回目のメッセージ
	msg1 := "こんにちは、私はGoエンジニアです。"
	fmt.Printf("User: %s\n", msg1)
	res1, err := cs.SendMessage(ctx, genai.Text(msg1))
	if err != nil {
		return fmt.Errorf("failed to send message 1: %w", err)
	}
	printResponse(res1)

	// 2回目のメッセージ（文脈依存）
	msg2 := "私が普段使っているプログラミング言語は何かわかりますか？"
	fmt.Printf("User: %s\n", msg2)
	res2, err := cs.SendMessage(ctx, genai.Text(msg2))
	if err != nil {
		return fmt.Errorf("failed to send message 2: %w", err)
	}
	printResponse(res2)

	fmt.Println()
	return nil
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				fmt.Printf("Gemini: %s\n", txt)
			}
		}
	}
}
