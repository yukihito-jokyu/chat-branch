-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    uuid CHAR(36) NOT NULL COMMENT 'UUID',
    google_id VARCHAR(255) COMMENT 'Google認証時のID (ゲストの場合はNULL)',
    provider VARCHAR(50) NOT NULL DEFAULT 'guest' COMMENT 'google または guest',
    name VARCHAR(255) NOT NULL COMMENT '表示名',
    avatar_url TEXT COMMENT 'アイコン画像のURL',
    created_id VARCHAR(36) NOT NULL COMMENT '作成者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    updated_id VARCHAR(36) NULL COMMENT '更新者ID',
    PRIMARY KEY (uuid),
    UNIQUE KEY uk_google_id (google_id)
) COMMENT='ユーザー管理テーブル';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
