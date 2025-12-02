package router

import (
	"backend/config"
	"backend/internal/handler"
	internalMiddleware "backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/usecase"

	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

// アプリケーションのルーティングを初期化する処理
func InitRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, genaiClient *genai.Client, publisher message.Publisher) {
	// ミドルウェア
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // エラーをグローバルエラーハンドラに転送し、適切なステータスコードを決定できるようにする
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(c.Request().Context(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	h := handler.New(db)

	// Auth の依存関係注入
	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authUsecase, cfg)

	// Project の依存関係注入
	projectRepo := repository.NewProjectRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	txManager := repository.NewTransactionManager(db)
	projectUsecase := usecase.NewProjectUsecase(projectRepo, chatRepo, messageRepo, txManager)
	projectHandler := handler.NewProjectHandler(projectUsecase)

	// Chat の依存関係注入
	genaiClientWrapper := usecase.NewGenAIClientWrapper(genaiClient)
	messageSelectionRepo := repository.NewMessageSelectionRepository(db)
	chatUsecase := usecase.NewChatUsecase(chatRepo, messageRepo, messageSelectionRepo, txManager, genaiClientWrapper, publisher)
	chatHandler := handler.NewChatHandler(chatUsecase)

	// Middleware の初期化
	authMiddleware := internalMiddleware.NewAuthMiddleware(cfg)

	e.GET("/health", h.HealthCheck)

	// auth関連
	{
		auth_router := e.Group("/api/auth")
		// ゲストでサインアップする機能
		auth_router.POST("/signup", authHandler.Signup)
		// ゲストでログインする機能
		auth_router.POST("/login", authHandler.Login)
	}

	// project関連
	{
		project_router := e.Group("/api/projects")
		project_router.Use(authMiddleware.Authenticate)
		// ユーザーが過去に作成したプロジェクト一覧を取得する
		project_router.GET("", projectHandler.GetProjects)
		// 新しいプロジェクトを作成する
		project_router.POST("", projectHandler.CreateProject)
	}

	// chat関連
	{
		chat_router := e.Group("/api/chats")
		chat_router.Use(authMiddleware.Authenticate)
		// 特定のチャットの基本情報を取得する機能
		chat_router.GET("/:chat_uuid", chatHandler.GetChat)
		// 特定のチャット内の会話履歴を取得する機能
		chat_router.GET("/:chat_uuid/messages", chatHandler.GetMessages)
		// 特定のチャットにメッセージを送信する機能
		chat_router.POST("/:chat_uuid/message", chatHandler.SendMessage)
		// 特定のチャットにLLMによる文章を生成する機能(POST /api/chats/:chat_uuid/message の後に必ず呼び出す)
		chat_router.GET("/:chat_uuid/messages/stream", chatHandler.StreamMessage)
		// 特定のチャットにLLMによる文章を生成する機能(初めてのチャット POST /api/projects の後に必ず呼び出す)
		chat_router.GET("/:chat_uuid/stream", chatHandler.FirstStreamChat)
		// 子チャット開始モーダルで、ユーザーが親チャットの要約を選択した場合、APIが実行され、ユーザーに確認させるためのプレビューを取得する機能
		chat_router.POST("/:chat_uuid/fork/preview", chatHandler.GenerateForkPreview)
		// 子チャットを生成する機能
		chat_router.POST("/:chat_uuid/fork", chatHandler.ForkChat)
		// 親にマージボタンを押した際、AIに子チャットの議論の流れと結論を要約を作らせる機能
		chat_router.POST("/:chat_uuid/merge/preview", chatHandler.GetMergePreview)
		// 子チャットを親チャットにマージする機能
		chat_router.POST("/:chat_uuid/merge", chatHandler.MergeChat)
		// チャットを閉じる機能
		chat_router.POST("/:chat_uuid/close", chatHandler.CloseChat)
		// チャットを開く機能
		chat_router.POST("/:chat_uuid/open", chatHandler.OpenChat)
	}
}
