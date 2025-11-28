-- +goose Up

-- 1. ユーザー作成
CREATE TABLE users (
    uuid VARCHAR(255) NOT NULL COMMENT 'UUID',
    google_id VARCHAR(255) COMMENT 'Google認証時のID',
    provider VARCHAR(50) NOT NULL DEFAULT 'guest' COMMENT 'google or guest',
    name VARCHAR(255) NOT NULL COMMENT '表示名',
    avatar_url TEXT COMMENT 'アイコン画像のURL',
    created_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_id VARCHAR(255) NULL,
    PRIMARY KEY (uuid),
    UNIQUE KEY uk_google_id (google_id)
) COMMENT='ユーザー管理テーブル';

-- 2. プロジェクト作成
CREATE TABLE projects (
    uuid VARCHAR(255) NOT NULL COMMENT 'UUID',
    user_uuid VARCHAR(255) NOT NULL COMMENT '所有ユーザーID',
    title VARCHAR(255) NOT NULL COMMENT 'プロジェクト名',
    created_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_id VARCHAR(255) NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_projects_user FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
) COMMENT='プロジェクト管理テーブル';

-- 3. チャット作成
CREATE TABLE chats (
    uuid VARCHAR(255) NOT NULL COMMENT 'UUID',
    project_uuid VARCHAR(255) NOT NULL COMMENT '所属プロジェクトID',
    parent_chat_uuid VARCHAR(255) NULL COMMENT '親チャットID',

    source_message_uuid VARCHAR(255) NULL COMMENT '分岐元の親メッセージID',
    selected_text TEXT NULL COMMENT '分岐時に選択されたテキスト',
    range_start INT NULL COMMENT '選択テキストの開始位置',
    range_end INT NULL COMMENT '選択テキストの終了位置',

    title VARCHAR(255) NOT NULL COMMENT 'チャットルーム名',
    status ENUM('open', 'merged', 'closed') NOT NULL DEFAULT 'open',
    context_summary TEXT NULL COMMENT '親から引き継いだ文脈',
    position_x FLOAT DEFAULT 0,
    position_y FLOAT DEFAULT 0,

    created_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_id VARCHAR(255) NULL,

    PRIMARY KEY (uuid),
    CONSTRAINT fk_chats_project FOREIGN KEY (project_uuid) REFERENCES projects(uuid) ON DELETE CASCADE,
    CONSTRAINT fk_chats_parent FOREIGN KEY (parent_chat_uuid) REFERENCES chats(uuid) ON DELETE CASCADE
) COMMENT='チャットルーム（ノード）管理テーブル';

-- 4. メッセージ作成
CREATE TABLE messages (
    uuid VARCHAR(255) NOT NULL COMMENT 'UUID',
    chat_uuid VARCHAR(255) NOT NULL COMMENT '所属するチャットID',

    role ENUM('user', 'assistant', 'system', 'merge_report') NOT NULL,
    content LONGTEXT NOT NULL,
    context_summary TEXT NULL,
    source_chat_uuid VARCHAR(255) NULL,

    created_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_id VARCHAR(255) NULL,

    PRIMARY KEY (uuid),
    CONSTRAINT fk_messages_chat FOREIGN KEY (chat_uuid) REFERENCES chats(uuid) ON DELETE CASCADE,
    CONSTRAINT fk_messages_source_chat FOREIGN KEY (source_chat_uuid) REFERENCES chats(uuid) ON DELETE SET NULL
) COMMENT='チャットメッセージ履歴テーブル';

ALTER TABLE chats
ADD CONSTRAINT fk_source_message
FOREIGN KEY (source_message_uuid) REFERENCES messages(uuid) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE chats DROP FOREIGN KEY fk_source_message;

DROP TABLE messages;
DROP TABLE chats;
DROP TABLE projects;
DROP TABLE users;
