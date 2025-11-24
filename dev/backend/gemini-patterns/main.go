package main

import (
	"context"
	"fmt"
	"log"

	"backend/internal/client"
	"backend/internal/config"
	"backend/internal/patterns"
)

func main() {
	// 1. 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	ctx := context.Background()

	// 2. クライアントの初期化
	c, err := client.NewClient(ctx, cfg.APIKey)
	if err != nil {
		log.Fatalf("Client initialization error: %v", err)
	}
	defer c.Close()

	// 3. 各パターンの実行
	if err := patterns.RunSimpleText(ctx, c); err != nil {
		log.Printf("Pattern 1 failed: %v", err)
	}

	if err := patterns.RunChatSession(ctx, c); err != nil {
		log.Printf("Pattern 2 failed: %v", err)
	}

	if err := patterns.RunStreaming(ctx, c); err != nil {
		log.Printf("Pattern 4 failed: %v", err)
	}

	fmt.Println("All patterns completed.")
}
