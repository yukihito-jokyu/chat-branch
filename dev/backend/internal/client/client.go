package client

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	ModelGeminiPro       = "gemini-2.5-flash"
	ModelGeminiProVision = "gemini-2.5-flash" // 1.5 flash supports both text and images
)

// NewClient は新しいGeminiクライアントを作成します
func NewClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return client, nil
}
