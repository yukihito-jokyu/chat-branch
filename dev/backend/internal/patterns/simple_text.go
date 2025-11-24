package patterns

import (
	"context"
	"fmt"

	"backend/internal/client"

	"github.com/google/generative-ai-go/genai"
)

// RunSimpleText は単純なテキスト生成を実行します
func RunSimpleText(ctx context.Context, c *genai.Client) error {
	fmt.Println("=== Pattern 1: Simple Text Generation ===")

	model := c.GenerativeModel(client.ModelGeminiPro)
	prompt := "Go言語の魅力を一言で教えてください。"

	fmt.Printf("Prompt: %s\n", prompt)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("no content generated")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			fmt.Printf("Response: %s\n", txt)
		}
	}
	fmt.Println()
	return nil
}
