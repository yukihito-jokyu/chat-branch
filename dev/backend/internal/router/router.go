package router

import (
	"backend/config"
	"backend/internal/handler"
	internalMiddleware "backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/usecase"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

// アプリケーションのルーティングを初期化する処理
func InitRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, genaiClient *genai.Client) {
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
	chatUsecase := usecase.NewChatUsecase(chatRepo, messageRepo, genaiClientWrapper)
	chatHandler := handler.NewChatHandler(chatUsecase)

	// Middleware の初期化
	authMiddleware := internalMiddleware.NewAuthMiddleware(cfg)

	e.GET("/health", h.HealthCheck)

	// auth関連
	{
		auth_router := e.Group("/api/auth")
		auth_router.POST("/signup", authHandler.Signup)
		auth_router.POST("/login", authHandler.Login)
	}

	// project関連
	{
		project_router := e.Group("/api/projects")
		project_router.Use(authMiddleware.Authenticate)
		project_router.GET("", projectHandler.GetProjects)
		project_router.POST("", projectHandler.CreateProject)
	}

	// chat関連
	{
		chat_router := e.Group("/api/chats")
		chat_router.Use(authMiddleware.Authenticate)
		chat_router.GET("/:chat_uuid", chatHandler.GetChat)
		chat_router.GET("/:chat_uuid/stream", chatHandler.FirstStreamChat)
	}
}
