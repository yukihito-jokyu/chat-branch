package main

import (
	"backend/config"
	"backend/internal/infrastructure/queue"
	"backend/internal/repository"
	"backend/internal/router"
	"backend/internal/usecase"
	"backend/internal/worker"
	"backend/pkg/logger"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"google.golang.org/genai"
	"gorm.io/gorm"
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

	// GenAI クライアントの初期化
	genaiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: cfg.Gemini.APIKey,
	})
	if err != nil {
		log.Fatalf("GenAIクライアントの作成に失敗: %v", err)
	}
	// defer genaiClient.Close() // main関数終了時に閉じる必要はないが、明示的に書くならここ

	// Watermill Publisher の初期化
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("DB接続の取得に失敗: %v", err)
	}
	publisher, err := queue.NewPublisher(sqlDB, slog.Default())
	if err != nil {
		log.Fatalf("Publisherの作成に失敗: %v", err)
	}
	defer publisher.Close()

	// Watermill Subscriber の初期化
	subscriber, err := queue.NewSubscriber(sqlDB, slog.Default())
	if err != nil {
		log.Fatalf("Subscriberの作成に失敗: %v", err)
	}
	defer subscriber.Close()

	// Worker の初期化と起動
	summaryWorker := setupWorker(db, genaiClient, subscriber)
	go func() {
		if err := summaryWorker.Run(context.Background()); err != nil {
			slog.Error("SummaryWorker failed", "error", err)
		}
	}()

	// サーバーの初期化
	e := setupServer(cfg, db, genaiClient, publisher)

	// サーバーの起動
	e.Logger.Fatal(e.Start(cfg.Server.Address))
}

// Workerの依存関係を初期化する
func setupWorker(db *gorm.DB, genaiClient *genai.Client, subscriber message.Subscriber) *worker.SummaryWorker {
	messageRepo := repository.NewMessageRepository(db)
	genaiClientWrapper := usecase.NewGenAIClientWrapper(genaiClient)
	return worker.NewSummaryWorker(subscriber, messageRepo, genaiClientWrapper)
}

// サーバーの依存関係を初期化する
func setupServer(cfg *config.Config, db *gorm.DB, genaiClient *genai.Client, publisher message.Publisher) *echo.Echo {
	e := echo.New()
	router.InitRoutes(e, db, cfg, genaiClient, publisher)
	return e
}
