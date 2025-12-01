-- +goose Up
-- 1. 選択テキスト情報を管理する新規テーブル作成
CREATE TABLE message_selections (
    uuid VARCHAR(255) NOT NULL COMMENT 'UUID',
    selected_text TEXT NULL COMMENT '選択されたテキスト',
    range_start INT NULL COMMENT '選択テキストの開始位置',
    range_end INT NULL COMMENT '選択テキストの終了位置',

    created_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_id VARCHAR(255) NULL,

    PRIMARY KEY (uuid)
) COMMENT='メッセージ内のテキスト選択情報';

-- 2. chats テーブルの修正
-- 新しいFKカラムの追加
ALTER TABLE chats
ADD COLUMN message_selection_uuid VARCHAR(255) NULL COMMENT '選択テキスト情報ID' AFTER source_message_uuid;

-- FK制約の追加
ALTER TABLE chats
ADD CONSTRAINT fk_chats_selection
FOREIGN KEY (message_selection_uuid) REFERENCES message_selections(uuid) ON DELETE SET NULL;

-- 既存カラムの削除（※注意: 既存データがある場合、ここのデータは消失します）
ALTER TABLE chats DROP COLUMN selected_text;
ALTER TABLE chats DROP COLUMN range_start;
ALTER TABLE chats DROP COLUMN range_end;

-- 3. messages テーブルの修正
-- 新しいFKカラムの追加
ALTER TABLE messages
ADD COLUMN message_selection_uuid VARCHAR(255) NULL COMMENT '関連する選択テキスト情報ID' AFTER source_chat_uuid;

-- FK制約の追加
ALTER TABLE messages
ADD CONSTRAINT fk_messages_selection
FOREIGN KEY (message_selection_uuid) REFERENCES message_selections(uuid) ON DELETE SET NULL;

-- +goose Down
-- 1. messages テーブルの戻し
ALTER TABLE messages DROP FOREIGN KEY fk_messages_selection;
ALTER TABLE messages DROP COLUMN message_selection_uuid;

-- 2. chats テーブルの戻し
-- 削除したカラムの復元
ALTER TABLE chats ADD COLUMN selected_text TEXT NULL COMMENT '分岐時に選択されたテキスト' AFTER source_message_uuid;
ALTER TABLE chats ADD COLUMN range_start INT NULL COMMENT '選択テキストの開始位置' AFTER selected_text;
ALTER TABLE chats ADD COLUMN range_end INT NULL COMMENT '選択テキストの終了位置' AFTER range_start;

-- FKとカラムの削除
ALTER TABLE chats DROP FOREIGN KEY fk_chats_selection;
ALTER TABLE chats DROP COLUMN message_selection_uuid;

-- 3. 新規テーブルの削除
DROP TABLE message_selections;
