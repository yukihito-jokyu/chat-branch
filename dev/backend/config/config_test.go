package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// 正常な設定ファイルの内容
	validConfigContent := `
server:
  address: ":8080"
database:
  user: "testuser"
  password: "testpassword"
  host: "localhost"
  port: "3306"
  dbname: "testdb"
jwt:
  secret: "testsecret"
  expiration: 24h
`

	// 不正なYAMLファイルの内容
	invalidConfigContent := `
server:
  address: ":8080"
database:
  user: "testuser"
  password: [
`

	type args struct {
		path    string
		content string // テスト実行時に作成するファイルの内容
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "正常系: 正しい設定ファイルを読み込めること",
			args: args{
				path:    "valid_config.yaml",
				content: validConfigContent,
			},
			want: &Config{
				Server: ServerConfig{
					Address: ":8080",
				},
				Database: DatabaseConfig{
					User:     "testuser",
					Password: "testpassword",
					Host:     "localhost",
					Port:     "3306",
					DBName:   "testdb",
				},
				JWT: JWTConfig{
					Secret:     "testsecret",
					Expiration: 24 * time.Hour,
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: 存在しないファイルを指定した場合エラーになること",
			args: args{
				path:    "non_existent_config.yaml",
				content: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: 不正なYAML形式の場合エラーになること",
			args: args{
				path:    "invalid_config.yaml",
				content: invalidConfigContent,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用ファイルの作成（contentが空でない場合のみ）
			if tt.args.content != "" {
				err := os.WriteFile(tt.args.path, []byte(tt.args.content), 0644)
				if err != nil {
					t.Fatalf("failed to create temp config file: %v", err)
				}
				defer os.Remove(tt.args.path)
			}

			got, err := LoadConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
