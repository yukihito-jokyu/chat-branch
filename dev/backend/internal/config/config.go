package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーション設定を保持します
type Config struct {
	APIKey string
}

// LoadConfig は環境変数から設定を読み込みます
func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY environment variable is not set")
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}
