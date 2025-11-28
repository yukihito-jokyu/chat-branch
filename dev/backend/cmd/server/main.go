package main

import (
	"backend/config"
	"backend/internal/router"
	"backend/pkg/logger"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

// アプリケーションのエントリーポイント
func main() {
	// 設定ファイルのパス
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yml"
	}

	// 設定の読み込み
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗: %v", err)
	}

	// ロガーの初期化
	logger.InitLogger(&cfg.Logger)

	// DB接続
	db, err := config.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("データベース接続に失敗: %v", err)
	}

	// Echo インスタンス
	e := echo.New()

	// ルーティング
	router.InitRoutes(e, db, cfg)

	// サーバーの起動
	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
