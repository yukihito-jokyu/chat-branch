package patterns

import (
	"context"
	"fmt"

	"backend/internal/client"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// RunStreaming はストリーミング応答を実行します
func RunStreaming(ctx context.Context, c *genai.Client) error {
	fmt.Println("=== Pattern 4: Streaming Response ===")

	model := c.GenerativeModel(client.ModelGeminiPro)
	prompt := "Go言語の歴史について300文字程度で語ってください。"

	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Print("Response: ")

	iter := model.GenerateContentStream(ctx, genai.Text(prompt))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to generate content stream: %w", err)
		}

		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					fmt.Print(txt)
				}
			}
		}
	}
	fmt.Println("\n")
	return nil
}
