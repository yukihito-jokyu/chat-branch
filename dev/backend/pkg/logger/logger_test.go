package logger

import (
	"backend/config"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	type args struct {
		cfg *config.LoggerConfig
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(t *testing.T)
	}{
		{
			name: "本番環境: JSONHandlerが設定されること",
			args: args{
				cfg: &config.LoggerConfig{
					IsProduction: true,
				},
			},
			checkFunc: func(t *testing.T) {
				handler := slog.Default().Handler()
				_, ok := handler.(*slog.JSONHandler)
				assert.True(t, ok, "Expected *slog.JSONHandler")
			},
		},
		{
			name: "開発環境: TextHandlerが設定されること",
			args: args{
				cfg: &config.LoggerConfig{
					IsProduction: false,
				},
			},
			checkFunc: func(t *testing.T) {
				handler := slog.Default().Handler()
				_, ok := handler.(*slog.TextHandler)
				assert.True(t, ok, "Expected *slog.TextHandler")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLogger(tt.args.cfg)
			tt.checkFunc(t)
		})
	}
}
