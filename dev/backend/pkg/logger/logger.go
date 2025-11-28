package logger

import (
	"backend/config"
	"log/slog"
	"os"
)

// 環境に基づいてグローバルロガーを初期化する処理
// isProduction が true の場合は JSONHandler を使用し、それ以外の場合は TextHandler を使用する。
func InitLogger(cfg *config.LoggerConfig) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	if cfg.IsProduction {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
